package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	encryptedHeaderSize               = 6
	encryptedBlockSize                = 1024
	encryptedLogicalPrefixSize        = encryptedBlockSize - 1
	encryptedRawOffsetDelta           = encryptedHeaderSize + 1
	maxMetadataEntrySize       uint32 = 2 * 1024 * 1024
	maxCodeMetadataSize        uint32 = 1 * 1024 * 1024
	maxAppNameLength                  = 128
)

const (
	appNameSourceLocalMetadata   = "local-metadata"
	appNameSourcePackageConfig   = "package-config"
	appNameSourceNavigationTitle = "navigation-title"
	appNameSourceCodeCandidate   = "code-candidate"
)

var (
	appConfigNameFields = []string{
		"appName",
		"miniProgramName",
		"nickname",
		"brandName",
		"storeName",
		"mallName",
		"name",
		"title",
	}
	localNameFields = []string{
		"appName",
		"miniProgramName",
		"nickname",
		"brandName",
		"storeName",
		"mallName",
		"name",
		"title",
	}
	codeNamePatterns = []*regexp.Regexp{
		regexp.MustCompile("(?is)[\"']?appId[\"']?\\s*[:=]\\s*[\"'][^\"']{4,64}[\"']\\s*,\\s*[\"']?appName[\"']?\\s*[:=]\\s*[\"']([^\"'\\r\\n]{1,128})[\"']\\s*,\\s*[\"']?appVersion[\"']?"),
		regexp.MustCompile("(?is)[\"']?appName[\"']?\\s*[:=]\\s*[\"']([^\"'\\r\\n]{1,128})[\"']\\s*,\\s*[\"']?appVersion[\"']?\\s*[:=]\\s*[\"'][^\"']{1,64}[\"']\\s*,\\s*[\"']?appVersionCode[\"']?\\s*[:=]"),
		regexp.MustCompile("(?is)authFloatInfo\\s*[:=]\\s*\\{[^}]{0,1000}\\bname\\s*[:=]\\s*[\"']([^\"'\\r\\n]{1,128})[\"']"),
		regexp.MustCompile("(?m)\\bAPPNAME\\s*[:=]\\s*[\"']([^\"'\\r\\n]{1,128})[\"']"),
	}
	codeNameFilePriority = map[string]int{
		"common/vendor.js": 0,
		"app-service.js":   1,
		"app.js":           2,
		"common/app.js":    3,
		"common/main.js":   4,
		"main.js":          5,
		"vendor.js":        6,
		"manifest.js":      7,
	}
	knownLocalMetadataFiles = []string{
		"appinfo.json",
		"appInfo.json",
		"app-info.json",
		"info.json",
		"adapter-config.json",
		"initial-rendering-cache-config.json",
	}
)

type appNameResult struct {
	name   string
	source string
}

type wxapkgIndexEntry struct {
	name   string
	offset uint32
	size   uint32
}

// wxapkgDataReader exposes the logical, decrypted wxapkg bytes through
// ReaderAt without loading the package body into memory.
type wxapkgDataReader struct {
	file            *os.File
	logicalSize     int64
	decryptedPrefix []byte
	xorKey          byte
}

func (r *wxapkgDataReader) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("negative wxapkg read offset")
	}
	if len(p) == 0 {
		return 0, nil
	}
	if off >= r.logicalSize {
		return 0, io.EOF
	}

	readLen := len(p)
	if remaining := r.logicalSize - off; int64(readLen) > remaining {
		readLen = int(remaining)
	}

	if len(r.decryptedPrefix) == 0 {
		n, err := r.file.ReadAt(p[:readLen], off)
		if n != readLen && err == nil {
			err = io.ErrUnexpectedEOF
		}
		if readLen != len(p) && err == nil {
			err = io.EOF
		}
		return n, err
	}

	read := 0
	for read < readLen {
		logicalOffset := off + int64(read)
		if logicalOffset < int64(len(r.decryptedPrefix)) {
			n := readLen - read
			if available := len(r.decryptedPrefix) - int(logicalOffset); n > available {
				n = available
			}
			copy(p[read:read+n], r.decryptedPrefix[int(logicalOffset):int(logicalOffset)+n])
			read += n
			continue
		}

		rawOffset := logicalOffset + encryptedRawOffsetDelta
		n := readLen - read
		nRead, err := r.file.ReadAt(p[read:read+n], rawOffset)
		for i := read; i < read+nRead; i++ {
			p[i] ^= r.xorKey
		}
		read += nRead
		if err != nil {
			return read, err
		}
		if nRead == 0 {
			return read, io.ErrUnexpectedEOF
		}
	}

	if readLen != len(p) {
		return read, io.EOF
	}
	return read, nil
}

func (r *wxapkgDataReader) Close() error {
	return r.file.Close()
}

func openWxapkgDataReader(path, wxid string) (*wxapkgDataReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	closeWithError := func(err error) (*wxapkgDataReader, error) {
		_ = file.Close()
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return closeWithError(err)
	}
	if info.Size() < 14 {
		return closeWithError(errors.Errorf("wxapkg file is too short: %d bytes", info.Size()))
	}

	marker := make([]byte, 14)
	if _, err := io.ReadFull(file, marker); err != nil {
		return closeWithError(err)
	}
	if marker[0] == 0xBE && marker[13] == 0xED {
		return &wxapkgDataReader{file: file, logicalSize: info.Size()}, nil
	}

	if wxid == "" {
		return closeWithError(errors.New("encrypted wxapkg requires a wxid"))
	}
	minimumSize := int64(encryptedHeaderSize + encryptedBlockSize)
	if info.Size() < minimumSize {
		return closeWithError(errors.Errorf("encrypted wxapkg file is too short: %d bytes", info.Size()))
	}

	encryptedBlock := make([]byte, encryptedBlockSize)
	if _, err := file.ReadAt(encryptedBlock, encryptedHeaderSize); err != nil {
		return closeWithError(err)
	}

	key := pbkdf2.Key([]byte(wxid), []byte("saltiest"), 1000, 32, sha1.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return closeWithError(err)
	}
	decryptedPrefix := make([]byte, encryptedBlockSize)
	cipher.NewCBCDecrypter(block, []byte("the iv: 16 bytes")).CryptBlocks(decryptedPrefix, encryptedBlock)

	xorKey := byte(0x66)
	if len(wxid) >= 2 {
		xorKey = wxid[len(wxid)-2]
	}
	return &wxapkgDataReader{
		file:            file,
		logicalSize:     info.Size() - encryptedRawOffsetDelta,
		decryptedPrefix: decryptedPrefix[:encryptedLogicalPrefixSize],
		xorKey:          xorKey,
	}, nil
}

func readAtFull(reader io.ReaderAt, data []byte, offset int64) error {
	if offset < 0 {
		return errors.New("negative wxapkg read offset")
	}
	n, err := reader.ReadAt(data, offset)
	if n != len(data) {
		if err == nil {
			err = io.ErrUnexpectedEOF
		}
		return err
	}
	return nil
}

func parseWxapkgIndex(reader io.ReaderAt, logicalSize int64) ([]wxapkgIndexEntry, error) {
	if logicalSize < 18 {
		return nil, errors.Errorf("wxapkg file is too short: %d bytes", logicalSize)
	}

	offset := int64(0)
	readByte := func() (byte, error) {
		var data [1]byte
		if err := readAtFull(reader, data[:], offset); err != nil {
			return 0, err
		}
		offset++
		return data[0], nil
	}
	readUint32 := func() (uint32, error) {
		var data [4]byte
		if err := readAtFull(reader, data[:], offset); err != nil {
			return 0, err
		}
		offset += 4
		return binary.BigEndian.Uint32(data[:]), nil
	}

	firstMark, err := readByte()
	if err != nil {
		return nil, errors.Errorf("read wxapkg first marker: %v", err)
	}
	if _, err := readUint32(); err != nil {
		return nil, errors.Errorf("read wxapkg header: %v", err)
	}
	if _, err := readUint32(); err != nil {
		return nil, errors.Errorf("read wxapkg index length: %v", err)
	}
	if _, err := readUint32(); err != nil {
		return nil, errors.Errorf("read wxapkg body length: %v", err)
	}
	lastMark, err := readByte()
	if err != nil {
		return nil, errors.Errorf("read wxapkg last marker: %v", err)
	}
	if firstMark != 0xBE || lastMark != 0xED {
		return nil, errors.New("invalid wxapkg markers")
	}

	fileCount, err := readUint32()
	if err != nil {
		return nil, errors.Errorf("read wxapkg file count: %v", err)
	}
	if fileCount > maxWxapkgFileCount {
		return nil, errors.Errorf("wxapkg file count %d exceeds limit %d", fileCount, maxWxapkgFileCount)
	}

	entries := make([]wxapkgIndexEntry, 0, int(fileCount))
	for i := uint32(0); i < fileCount; i++ {
		nameLen, err := readUint32()
		if err != nil {
			return nil, errors.Errorf("read wxapkg entry %d name length: %v", i, err)
		}
		if nameLen > maxWxapkgNameLength {
			return nil, errors.Errorf("wxapkg entry %d name length %d exceeds limit %d", i, nameLen, maxWxapkgNameLength)
		}
		if uint64(offset)+uint64(nameLen)+8 > uint64(logicalSize) {
			return nil, errors.Errorf("wxapkg entry %d index exceeds file boundary", i)
		}

		nameBytes := make([]byte, int(nameLen))
		if err := readAtFull(reader, nameBytes, offset); err != nil {
			return nil, errors.Errorf("read wxapkg entry %d name: %v", i, err)
		}
		offset += int64(nameLen)

		entryOffset, err := readUint32()
		if err != nil {
			return nil, errors.Errorf("read wxapkg entry %d offset: %v", i, err)
		}
		entrySize, err := readUint32()
		if err != nil {
			return nil, errors.Errorf("read wxapkg entry %d size: %v", i, err)
		}
		if entrySize > maxEntrySize {
			return nil, errors.Errorf("wxapkg entry %q size %d exceeds limit %d", string(nameBytes), entrySize, maxEntrySize)
		}
		if uint64(entryOffset)+uint64(entrySize) > uint64(logicalSize) {
			return nil, errors.Errorf("wxapkg entry %q exceeds file boundary", string(nameBytes))
		}

		name, err := normalizeArchiveEntryName(string(nameBytes))
		if err != nil {
			return nil, errors.Errorf("invalid wxapkg entry name %q: %v", string(nameBytes), err)
		}
		entries = append(entries, wxapkgIndexEntry{name: name, offset: entryOffset, size: entrySize})
	}
	return entries, nil
}

func readWxapkgAppName(path, wxid string) (string, error) {
	result, err := readWxapkgAppNameResult(path, wxid)
	return result.name, err
}

func readWxapkgAppNameResult(path, wxid string) (appNameResult, error) {
	reader, err := openWxapkgDataReader(path, wxid)
	if err != nil {
		return appNameResult{}, err
	}
	defer reader.Close()

	entries, err := parseWxapkgIndex(reader, reader.logicalSize)
	if err != nil {
		return appNameResult{}, err
	}

	var navigationResult appNameResult
	var codeResult appNameResult
	codePriority := int(^uint(0) >> 1)
	for _, entry := range entries {
		if entry.name == "app-config.json" {
			if entry.size > maxMetadataEntrySize {
				continue
			}
			data, err := readWxapkgEntry(reader, entry, maxMetadataEntrySize)
			if err != nil {
				return appNameResult{}, err
			}
			result := appNameFromConfigResult(data)
			if result.source == appNameSourcePackageConfig {
				return result, nil
			}
			if result.source == appNameSourceNavigationTitle && navigationResult.name == "" {
				navigationResult = result
			}
			continue
		}

		priority, ok := codeEntryPriority(entry.name)
		if !ok || priority >= codePriority {
			continue
		}
		data, err := readWxapkgEntrySample(reader, entry, maxCodeMetadataSize)
		if err != nil {
			continue
		}
		result := appNameFromCodeContent(data)
		if result.name != "" {
			codeResult = result
			codePriority = priority
		}
	}
	if navigationResult.name != "" {
		return navigationResult, nil
	}
	if codeResult.name != "" {
		return codeResult, nil
	}
	return appNameResult{}, nil
}

func appNameFromPackageData(data []byte, files []*WxapkgFileItemStructure) string {
	return appNameFromPackageDataResult(data, files).name
}

func appNameFromPackageDataResult(data []byte, files []*WxapkgFileItemStructure) appNameResult {
	var navigationResult appNameResult
	var codeResult appNameResult
	codePriority := int(^uint(0) >> 1)

	for _, file := range files {
		if file.name == "app-config.json" {
			configData, ok := packageEntryData(data, file, maxMetadataEntrySize)
			if !ok {
				continue
			}
			result := appNameFromConfigResult(configData)
			if result.source == appNameSourcePackageConfig {
				return result
			}
			if result.source == appNameSourceNavigationTitle && navigationResult.name == "" {
				navigationResult = result
			}
			continue
		}

		priority, ok := codeEntryPriority(file.name)
		if !ok || priority >= codePriority {
			continue
		}
		codeData, ok := packageEntryDataSample(data, file, maxCodeMetadataSize)
		if !ok {
			continue
		}
		result := appNameFromCodeContent(codeData)
		if result.name != "" {
			codeResult = result
			codePriority = priority
		}
	}
	if navigationResult.name != "" {
		return navigationResult
	}
	if codeResult.name != "" {
		return codeResult
	}
	return appNameResult{}
}

func readWxapkgEntry(reader io.ReaderAt, entry wxapkgIndexEntry, maxSize uint32) ([]byte, error) {
	if entry.size > maxSize {
		return nil, errors.Errorf("wxapkg entry %q is too large: %d bytes", entry.name, entry.size)
	}
	data := make([]byte, int(entry.size))
	if err := readAtFull(reader, data, int64(entry.offset)); err != nil {
		return nil, err
	}
	return data, nil
}

func readWxapkgEntrySample(reader io.ReaderAt, entry wxapkgIndexEntry, maxSize uint32) ([]byte, error) {
	if entry.size <= maxSize {
		return readWxapkgEntry(reader, entry, maxSize)
	}

	headSize := int(maxSize / 2)
	tailSize := int(maxSize) - headSize
	head := make([]byte, headSize)
	tail := make([]byte, tailSize)
	if err := readAtFull(reader, head, int64(entry.offset)); err != nil {
		return nil, err
	}
	tailOffset := int64(entry.offset) + int64(entry.size) - int64(tailSize)
	if err := readAtFull(reader, tail, tailOffset); err != nil {
		return nil, err
	}

	data := make([]byte, 0, len(head)+1+len(tail))
	data = append(data, head...)
	data = append(data, '\n')
	data = append(data, tail...)
	return data, nil
}

func packageEntryData(data []byte, file *WxapkgFileItemStructure, maxSize uint32) ([]byte, bool) {
	if file == nil || file.size > maxSize {
		return nil, false
	}
	start := uint64(file.offset)
	end := start + uint64(file.size)
	if end < start || end > uint64(len(data)) {
		return nil, false
	}
	return data[start:end], true
}

func packageEntryDataSample(data []byte, file *WxapkgFileItemStructure, maxSize uint32) ([]byte, bool) {
	if file == nil {
		return nil, false
	}
	if file.size <= maxSize {
		return packageEntryData(data, file, maxSize)
	}

	headSize := int(maxSize / 2)
	tailSize := int(maxSize) - headSize
	start := uint64(file.offset)
	end := start + uint64(file.size)
	if end < start || end > uint64(len(data)) {
		return nil, false
	}
	headEnd := start + uint64(headSize)
	tailStart := end - uint64(tailSize)
	if headEnd > end || tailStart < start {
		return nil, false
	}

	result := make([]byte, 0, headSize+1+tailSize)
	result = append(result, data[start:headEnd]...)
	result = append(result, '\n')
	result = append(result, data[tailStart:end]...)
	return result, true
}

func appNameFromConfig(data []byte) string {
	return appNameFromConfigResult(data).name
}

func appNameFromConfigResult(data []byte) appNameResult {
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return appNameResult{}
	}

	if result := appNameFromObject(config, appConfigNameFields, appNameSourcePackageConfig); result.name != "" {
		return result
	}
	for _, key := range []string{"global", "app", "appInfo", "miniProgram"} {
		if nested, ok := objectField(config, key); ok {
			if result := appNameFromObject(nested, appConfigNameFields, appNameSourcePackageConfig); result.name != "" {
				return result
			}
		}
	}

	global, ok := objectField(config, "global")
	if !ok {
		return appNameResult{}
	}
	window, ok := objectField(global, "window")
	if !ok {
		return appNameResult{}
	}
	return makeAppNameResult(stringField(window, []string{"navigationBarTitleText"}), appNameSourceNavigationTitle)
}

func appNameFromLocalConfig(data []byte) appNameResult {
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return appNameResult{}
	}

	if result := appNameFromObject(config, localNameFields, appNameSourceLocalMetadata); result.name != "" {
		return result
	}
	for _, key := range []string{"app", "appInfo", "miniProgram", "program", "data", "info"} {
		if nested, ok := objectField(config, key); ok {
			if result := appNameFromObject(nested, localNameFields, appNameSourceLocalMetadata); result.name != "" {
				return result
			}
		}
	}
	return appNameResult{}
}

func appNameFromObject(object map[string]interface{}, fields []string, source string) appNameResult {
	if object == nil {
		return appNameResult{}
	}
	return makeAppNameResult(stringField(object, fields), source)
}

func objectField(object map[string]interface{}, field string) (map[string]interface{}, bool) {
	if object == nil {
		return nil, false
	}
	for key, value := range object {
		if !strings.EqualFold(key, field) {
			continue
		}
		nested, ok := value.(map[string]interface{})
		return nested, ok
	}
	return nil, false
}

func stringField(object map[string]interface{}, fields []string) string {
	if object == nil {
		return ""
	}
	for _, field := range fields {
		for key, value := range object {
			if !strings.EqualFold(key, field) {
				continue
			}
			if name, ok := value.(string); ok {
				return name
			}
		}
	}
	return ""
}

func appNameFromCodeContent(data []byte) appNameResult {
	for _, pattern := range codeNamePatterns {
		match := pattern.FindSubmatch(data)
		if len(match) < 2 {
			continue
		}
		result := makeAppNameResult(string(match[1]), appNameSourceCodeCandidate)
		if result.name == "" {
			continue
		}
		if reWxId.MatchString(strings.ToLower(result.name)) ||
			strings.ContainsAny(result.name, "/\\") ||
			strings.EqualFold(result.name, "netscape") ||
			strings.EqualFold(result.name, "mozilla") {
			continue
		}
		if !isLikelyCodeName(result.name) {
			continue
		}
		return result
	}
	return appNameResult{}
}

func isLikelyCodeName(name string) bool {
	runes := []rune(name)
	if len(runes) >= 4 {
		return true
	}
	for _, value := range runes {
		if value >= 0x4e00 && value <= 0x9fff {
			return true
		}
	}
	return false
}

func codeEntryPriority(name string) (int, bool) {
	name = strings.ToLower(strings.TrimPrefix(name, "/"))
	if priority, ok := codeNameFilePriority[name]; ok {
		return priority, true
	}
	if separator := strings.LastIndexByte(name, '/'); separator >= 0 {
		if priority, ok := codeNameFilePriority[name[separator+1:]]; ok {
			return priority + 10, true
		}
	}
	return 0, false
}

func isGenericAppTitle(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "\u9996\u9875", "\u4e3b\u9875", "home", "index", "main", "welcome",
		"default", "login", "signin", "\u6211\u7684", "\u4e2a\u4eba\u4e2d\u5fc3",
		"\u8bbe\u7f6e", "\u6d88\u606f", "\u53d1\u73b0", "loading", "cart",
		"profile", "settings", "search", "orders", "order",
		"\u52a0\u8f7d\u4e2d", "\u52a0\u8f7d\u4e2d...", "\u652f\u4ed8",
		"\u5fae\u4fe1\u652f\u4ed8", "\u8d2d\u7269\u8f66", "\u63d0\u793a",
		"\u767b\u5f55", "\u4f1a\u5458\u767b\u5f55", "\u641c\u7d22", "\u8ba2\u5355",
		"\u8ba2\u5355\u5217\u8868", "\u8ba2\u5355\u8be6\u60c5", "\u6211\u7684\u8ba2\u5355",
		"\u6536\u8d27\u5730\u5740", "\u65b0\u589e\u6536\u8d27\u5730\u5740",
		"\u9009\u62e9\u95e8\u5e97", "\u9009\u62e9\u57ce\u5e02", "\u70b9\u9910",
		"\u624b\u673a\u70b9\u9910", "\u4e0b\u5355\u9875", "\u4f18\u60e0\u5238\u9875",
		"\u4f1a\u5458\u4e2d\u5fc3", "\u6d4b\u8bd5\u7ed3\u679c", "\u6d4b\u8bd5\u8be6\u60c5",
		"\u6d4b\u8bd5\u5386\u53f2", "\u66f4\u591a\u7ed3\u679c":
		return true
	default:
		return false
	}
}

func makeAppNameResult(name, source string) appNameResult {
	name = cleanAppName(name)
	if name == "" || isGenericAppTitle(name) {
		return appNameResult{}
	}
	return appNameResult{name: name, source: source}
}

func cleanAppName(name string) string {
	name = strings.Join(strings.Fields(strings.TrimSpace(name)), " ")
	if name == "" {
		return ""
	}
	runes := []rune(name)
	if len(runes) > maxAppNameLength {
		return string(runes[:maxAppNameLength])
	}
	return name
}

func wxIDFromPath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	for _, component := range strings.Split(path, "/") {
		if reWxId.MatchString(component) {
			return component
		}
	}
	return ""
}

func readLocalAppName(path, wxid string) appNameResult {
	root := appIDRootForPath(path, wxid)
	if root == "" {
		return appNameResult{}
	}
	localRoot := filepath.Join(filepath.Dir(filepath.Dir(root)), "local")
	localDir := filepath.Join(localRoot, wxid)
	if info, err := os.Stat(localDir); err != nil || !info.IsDir() {
		entries, err := os.ReadDir(localRoot)
		if err != nil {
			return appNameResult{}
		}
		for _, entry := range entries {
			if entry.IsDir() && strings.EqualFold(entry.Name(), wxid) {
				localDir = filepath.Join(localRoot, entry.Name())
				break
			}
		}
	}

	metadataEntries, err := os.ReadDir(localDir)
	if err != nil {
		return appNameResult{}
	}
	for _, wanted := range knownLocalMetadataFiles {
		for _, entry := range metadataEntries {
			if entry.IsDir() || !strings.EqualFold(entry.Name(), wanted) {
				continue
			}
			info, err := entry.Info()
			if err != nil || info.Size() < 0 || info.Size() > int64(maxMetadataEntrySize) {
				continue
			}
			data, err := os.ReadFile(filepath.Join(localDir, entry.Name()))
			if err != nil {
				continue
			}
			if result := appNameFromLocalConfig(data); result.name != "" {
				return result
			}
		}
	}
	return appNameResult{}
}

func appIDRootForPath(path, wxid string) string {
	if wxid == "" {
		wxid = wxIDFromPath(path)
	}
	if wxid == "" {
		return ""
	}
	current := filepath.Clean(path)
	if info, err := os.Stat(current); err == nil && !info.IsDir() {
		current = filepath.Dir(current)
	}
	for {
		if strings.EqualFold(filepath.Base(current), wxid) {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

func listDirectWxapkgFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".wxapkg") {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}
	return files, nil
}

func wxapkgFilePriority(path string) int {
	switch strings.ToLower(filepath.Base(path)) {
	case "__app__.wxapkg":
		return 0
	case "__full__.wxapkg":
		return 1
	case "__subpackage__.wxapkg":
		return 2
	case "__plugincode__.wxapkg":
		return 4
	default:
		return 3
	}
}

func fileModTimeNano(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.ModTime().UnixNano()
}

func chooseMainWxapkg(files []string) string {
	if len(files) == 0 {
		return ""
	}
	candidates := append([]string(nil), files...)
	sort.SliceStable(candidates, func(i, j int) bool {
		pi, pj := wxapkgFilePriority(candidates[i]), wxapkgFilePriority(candidates[j])
		if pi != pj {
			return pi < pj
		}
		ti, tj := fileModTimeNano(candidates[i]), fileModTimeNano(candidates[j])
		if ti != tj {
			return ti > tj
		}
		return strings.ToLower(candidates[i]) < strings.ToLower(candidates[j])
	})
	return candidates[0]
}

func compareVersionNames(left, right string) int {
	leftDigits := left != ""
	rightDigits := right != ""
	for _, value := range left {
		if value < '0' || value > '9' {
			leftDigits = false
			break
		}
	}
	for _, value := range right {
		if value < '0' || value > '9' {
			rightDigits = false
			break
		}
	}
	if leftDigits && rightDigits {
		leftTrimmed := strings.TrimLeft(left, "0")
		rightTrimmed := strings.TrimLeft(right, "0")
		if leftTrimmed == "" {
			leftTrimmed = "0"
		}
		if rightTrimmed == "" {
			rightTrimmed = "0"
		}
		if len(leftTrimmed) != len(rightTrimmed) {
			if len(leftTrimmed) > len(rightTrimmed) {
				return 1
			}
			return -1
		}
		if leftTrimmed != rightTrimmed {
			if leftTrimmed > rightTrimmed {
				return 1
			}
			return -1
		}
		return 0
	}
	if left == right {
		return 0
	}
	if left > right {
		return 1
	}
	return -1
}

func findMainWxapkg(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		return path
	}

	if files, err := listDirectWxapkgFiles(path); err == nil && len(files) > 0 {
		return chooseMainWxapkg(files)
	}

	type versionCandidate struct {
		path    string
		version string
		modTime int64
	}
	entries, err := os.ReadDir(path)
	if err == nil {
		candidates := make([]versionCandidate, 0)
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			versionPath := filepath.Join(path, entry.Name())
			files, err := listDirectWxapkgFiles(versionPath)
			if err != nil || len(files) == 0 {
				continue
			}
			mainPath := chooseMainWxapkg(files)
			candidates = append(candidates, versionCandidate{
				path:    mainPath,
				version: entry.Name(),
				modTime: fileModTimeNano(mainPath),
			})
		}
		if len(candidates) > 0 {
			sort.SliceStable(candidates, func(i, j int) bool {
				if candidates[i].modTime != candidates[j].modTime {
					return candidates[i].modTime > candidates[j].modTime
				}
				if versionCompare := compareVersionNames(candidates[i].version, candidates[j].version); versionCompare != 0 {
					return versionCompare > 0
				}
				return strings.ToLower(candidates[i].path) < strings.ToLower(candidates[j].path)
			})
			return candidates[0].path
		}
	}

	files, err := ListFilesWithExtension(path, ".wxapkg")
	if err != nil || len(files) == 0 {
		return ""
	}
	return chooseMainWxapkg(files)
}

func appNameForPathResult(path, wxid string) appNameResult {
	if result := readLocalAppName(path, wxid); result.name != "" {
		return result
	}
	packagePath := findMainWxapkg(path)
	if packagePath == "" {
		return appNameResult{}
	}
	result, err := readWxapkgAppNameResult(packagePath, wxid)
	if err != nil {
		return appNameResult{}
	}
	return result
}

func appNameForPath(path, wxid string) string {
	return appNameForPathResult(path, wxid).name
}
