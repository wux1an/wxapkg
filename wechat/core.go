package wechat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var reWxId = regexp.MustCompile(`^wx[0-9a-f]{16}$`)

const (
	maxWxapkgFileCount  uint32 = 102400
	maxWxapkgNameLength uint32 = 1024
	maxEntrySize        uint32 = 10 * 1024 * 1024
	maxUnpackWorkers           = 4
)

type StatusType string

var (
	StatusTypeRunning  StatusType = "running"
	StatusTypeFinished StatusType = "finished"
	StatusTypeError    StatusType = "error"
)

type UnpackOptions struct {
	EnableDecrypt      bool
	EnableJsBeautify   bool
	EnableHtmlBeautify bool
	EnableJsonBeautify bool
	OutputDir          string
	SavePath           string
}

type WxapkgItem struct {
	UUID            string
	WxId            string
	AppName         string
	AppNameSource   string
	Location        string
	EncryptKey      string
	Size            int64
	IsDir           bool
	LastModifyTime  int64
	WxapkgFilePaths []string

	UnpackStatus       StatusType
	UnpackCurrent      int64
	UnpackTotal        int64
	UnpackProgress     float64
	UnpackCurrentFile  string
	UnpackSavePath     string
	UnpackErrorMessage string
}

func (u *WxapkgItem) SetErrorState(msg string) {
	u.UnpackErrorMessage = msg
	u.UnpackStatus = StatusTypeError
}

func (u *WxapkgItem) IncreaseProgress(currentFile string) {
	u.UnpackCurrent++
	u.UnpackCurrentFile = currentFile
	if u.UnpackTotal != 0 {
		u.UnpackProgress = float64(u.UnpackCurrent) / float64(u.UnpackTotal) * 100
	}
}

type WxapkgFileItemStructure struct {
	name        string
	offset      uint32
	size        uint32
	savePath    string
	sourcePath  string
	archivePath string
}

// PathScanResult GetDefaultPaths 的返回类型
type PathScanResult struct {
	Paths []string // 检测到的小程序安装目录
	Logs  string   // 检测过程的日志详情
}

type platform interface {
	GetDefaultPaths() PathScanResult
}

var Platform = newPlatform()

type Unpacker struct {
	item    *WxapkgItem
	options *UnpackOptions
	locker  sync.Mutex

	files        []*WxapkgFileItemStructure
	fileNameMap  map[string]bool
	tempFilePath []string
}

func NewUnpacker(item *WxapkgItem, options *UnpackOptions) *Unpacker {
	return &Unpacker{
		item:         item,
		options:      options,
		locker:       sync.Mutex{},
		files:        []*WxapkgFileItemStructure{},
		fileNameMap:  map[string]bool{},
		tempFilePath: []string{},
	}
}

func (u *Unpacker) init() error {
	if u == nil || u.item == nil {
		return errors.New("解包任务未提供有效的小程序项目")
	}
	if u.options == nil {
		return errors.New("解包任务未提供解包配置")
	}

	outputDir, err := filepath.Abs(u.options.OutputDir)
	if err != nil {
		return errors.Errorf("检查输出目录 '%s' 出错，%v", u.options.OutputDir, err)
	}
	savePath := strings.TrimSpace(u.options.SavePath)
	if savePath == "" {
		return errors.New("解包输出目录不能为空")
	}
	savePath, err = filepath.Abs(savePath)
	if err != nil {
		return errors.Errorf("检查解包输出目录 '%s' 出错，%v", u.options.SavePath, err)
	}
	u.options.OutputDir = outputDir
	u.options.SavePath = savePath
	u.item.UnpackSavePath = savePath

	if u.item.IsDir {
		files, err := ListFilesWithExtension(u.item.Location, ".wxapkg")
		if err != nil {
			return errors.Errorf("扫描目录 '%s' 下 wxapkg 文件失败，%v", u.item.Location, err)
		}
		if len(files) == 0 {
			return errors.Errorf("目录 '%s' 下没有 wxapkg 文件", u.item.Location)
		}
		u.item.WxapkgFilePaths = files
	} else {
		u.item.WxapkgFilePaths = []string{u.item.Location}
	}

	return nil
}

func (u *Unpacker) analyzeAll() error {
	u.files = nil
	u.fileNameMap = make(map[string]bool)
	u.cleanupTempFiles()

	for _, wxapkgFile := range u.item.WxapkgFilePaths {
		fileData, err := os.ReadFile(wxapkgFile)
		if err != nil {
			return err
		}
		sourcePath := wxapkgFile

		if !u.isDecryptedWxapkgFile(fileData) {
			if !u.options.EnableDecrypt {
				return errors.Errorf("小程序文件 '%s' 为加密文件，请在解包配置中启用解密", wxapkgFile)
			}
			if u.item.EncryptKey == "" {
				return errors.Errorf("小程序文件 '%s' 为加密文件，未设置解密密钥，秘钥为小程序的 wxid，格式：^wx[0-9a-f]{16}$", wxapkgFile)
			}

			fileData, err = decryptWxapkgFile(u.item.EncryptKey, fileData)
			if err != nil {
				return errors.Errorf("解密小程序文件 '%s' 失败, %v", wxapkgFile, err)
			}
			sourcePath, err = u.writeTempSource(fileData)
			if err != nil {
				return errors.Errorf("保存解密后的小程序文件 '%s' 失败, %v", wxapkgFile, err)
			}
		}

		files, err := u.analyze(fileData, wxapkgFile, sourcePath)
		if err != nil {
			return errors.Errorf("解析小程序文件 '%s' 失败, %v", wxapkgFile, err)
		}

		if u.item.AppName == "" {
			result := appNameFromPackageDataResult(fileData, files)
			u.item.AppName = result.name
			u.item.AppNameSource = result.source
		}
		u.files = append(u.files, files...)
	}

	return nil
}

func (u *Unpacker) writeTempSource(data []byte) (string, error) {
	file, err := os.CreateTemp("", "wxapkg-decrypted-*")
	if err != nil {
		return "", err
	}
	tempPath := file.Name()
	n, writeErr := file.Write(data)
	if writeErr == nil && n != len(data) {
		writeErr = io.ErrShortWrite
	}
	closeErr := file.Close()
	if writeErr != nil {
		_ = os.Remove(tempPath)
		return "", writeErr
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return "", closeErr
	}
	u.tempFilePath = append(u.tempFilePath, tempPath)
	return tempPath, nil
}

func (u *Unpacker) cleanupTempFiles() {
	for _, tempPath := range u.tempFilePath {
		_ = os.Remove(tempPath)
	}
	u.tempFilePath = nil
}

func (u *Unpacker) isDecryptedWxapkgFile(data []byte) bool {
	if len(data) < 14 {
		return false
	}

	return data[0] == 0xBE && data[13] == 0xED // the firstMark and lastMark
}

func (u *Unpacker) analyze(data []byte, archivePath, sourcePath string) ([]*WxapkgFileItemStructure, error) {
	if u == nil || u.item == nil {
		return nil, errors.New("解包任务未提供有效的小程序项目")
	}
	if u.fileNameMap == nil {
		u.fileNameMap = make(map[string]bool)
	}

	f := bytes.NewReader(data)
	var (
		firstMark       uint8
		info1           uint32
		indexInfoLength uint32
		bodyInfoLength  uint32
		lastMark        uint8
	)
	if err := binary.Read(f, binary.BigEndian, &firstMark); err != nil {
		return nil, errors.Errorf("读取 wxapkg 头失败: %v", err)
	}
	if err := binary.Read(f, binary.BigEndian, &info1); err != nil {
		return nil, errors.Errorf("读取 wxapkg 头失败: %v", err)
	}
	if err := binary.Read(f, binary.BigEndian, &indexInfoLength); err != nil {
		return nil, errors.Errorf("读取 wxapkg 头失败: %v", err)
	}
	if err := binary.Read(f, binary.BigEndian, &bodyInfoLength); err != nil {
		return nil, errors.Errorf("读取 wxapkg 头失败: %v", err)
	}
	if err := binary.Read(f, binary.BigEndian, &lastMark); err != nil {
		return nil, errors.Errorf("读取 wxapkg 头失败: %v", err)
	}

	if firstMark != 0xBE || lastMark != 0xED {
		return nil, errors.New("wxapkg 文件结构不合法")
	}

	var fileCount uint32
	if err := binary.Read(f, binary.BigEndian, &fileCount); err != nil {
		return nil, errors.Errorf("读取 wxapkg 文件数量失败: %v", err)
	}
	if fileCount > maxWxapkgFileCount {
		return nil, errors.Errorf("文件总数量 %d 超出上限 %d", fileCount, maxWxapkgFileCount)
	}

	result := make([]*WxapkgFileItemStructure, fileCount)
	for i := uint32(0); i < fileCount; i++ {
		var nameLen uint32
		if err := binary.Read(f, binary.BigEndian, &nameLen); err != nil {
			return nil, errors.Errorf("读取第 %d 个文件名长度失败: %v", i, err)
		}
		if nameLen > maxWxapkgNameLength {
			return nil, errors.Errorf("文件名长度 %d 超出上限 %d 字节", nameLen, maxWxapkgNameLength)
		}

		nameBytes := make([]byte, int(nameLen))
		if _, err := io.ReadFull(f, nameBytes); err != nil {
			return nil, errors.Errorf("读取第 %d 个文件名失败: %v", i, err)
		}

		item := &WxapkgFileItemStructure{
			sourcePath:  sourcePath,
			archivePath: archivePath,
		}
		if err := binary.Read(f, binary.BigEndian, &item.offset); err != nil {
			return nil, errors.Errorf("读取文件 '%s' 偏移量失败: %v", string(nameBytes), err)
		}
		if err := binary.Read(f, binary.BigEndian, &item.size); err != nil {
			return nil, errors.Errorf("读取文件 '%s' 大小失败: %v", string(nameBytes), err)
		}
		if item.size > maxEntrySize {
			return nil, errors.Errorf("文件 '%s' 标记长度 %d 超出上限 %d MB", string(nameBytes), item.size, maxEntrySize/(1024*1024))
		}
		if uint64(item.offset)+uint64(item.size) > uint64(len(data)) {
			return nil, errors.Errorf("文件 '%s' 的数据范围超出 wxapkg 文件边界", string(nameBytes))
		}

		name, err := normalizeArchiveEntryName(string(nameBytes))
		if err != nil {
			return nil, errors.Errorf("文件名 '%s' 不合法: %v", string(nameBytes), err)
		}

		// If a package contains duplicate names, keep all entries without overwriting.
		item.name = name
		for j := 1; u.fileNameMap[item.name]; j++ {
			dot := strings.LastIndex(name, ".")
			if dot == -1 {
				item.name = fmt.Sprintf("%s-%d", name, j)
			} else {
				item.name = fmt.Sprintf("%s-%d%s", name[:dot], j, name[dot:])
			}
		}
		u.fileNameMap[item.name] = true

		item.savePath, err = joinSafeOutputPath(u.item.UnpackSavePath, item.name)
		if err != nil {
			return nil, errors.Errorf("文件名 %s 会导致目录穿越: %v", item.name, err)
		}
		result[i] = item
	}

	return result, nil
}

func normalizeArchiveEntryName(name string) (string, error) {
	if name == "" {
		return "", errors.New("文件名不能为空")
	}
	if strings.IndexByte(name, 0) >= 0 {
		return "", errors.New("文件名不能包含 NUL 字节")
	}

	// wxapkg names are slash-separated virtual paths. Apply these checks
	// independently of the host OS before converting to filepath paths.
	name = strings.ReplaceAll(name, "\\", "/")
	if strings.HasPrefix(name, "//") {
		return "", errors.New("不允许使用 UNC 路径")
	}
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		return "", errors.New("文件名不能为空")
	}

	cleaned := path.Clean(name)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", errors.New("不允许使用目录穿越路径")
	}
	if len(cleaned) >= 2 && cleaned[1] == ':' {
		return "", errors.New("不允许使用盘符路径")
	}
	if path.IsAbs(cleaned) {
		return "", errors.New("不允许使用绝对路径")
	}
	return cleaned, nil
}

func joinSafeOutputPath(root, name string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	target, err := filepath.Abs(filepath.Join(rootAbs, filepath.FromSlash(name)))
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(rootAbs, target)
	if err != nil {
		return "", err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return "", errors.New("目标路径不在解包输出目录内")
	}
	return target, nil
}

func (u *Unpacker) UnpackWithStatusCallback(callback func(item *WxapkgItem)) {
	if callback == nil {
		callback = func(*WxapkgItem) {}
	}
	if u == nil || u.item == nil {
		return
	}
	defer u.cleanupTempFiles()

	if err := u.init(); err != nil {
		u.item.SetErrorState(err.Error())
		callback(u.item)
		return
	}

	if err := u.analyzeAll(); err != nil {
		u.item.SetErrorState(err.Error())
		callback(u.item)
		return
	}

	u.item.UnpackCurrent = 0
	u.item.UnpackTotal = int64(len(u.files))
	u.item.UnpackProgress = 0
	u.item.UnpackCurrentFile = ""
	u.item.UnpackErrorMessage = ""
	u.item.UnpackStatus = StatusTypeRunning
	callback(u.item)

	var hasError = u.unpack(defaultUnpackWorkerCount(), callback)

	if !hasError {
		u.item.UnpackStatus = StatusTypeFinished
		callback(u.item)
	}
}

func defaultUnpackWorkerCount() int {
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		return 1
	}
	if workers > maxUnpackWorkers {
		return maxUnpackWorkers
	}
	return workers
}

func (u *Unpacker) unpackOne(d *WxapkgFileItemStructure) (err error) {
	entryName := "<unknown>"
	archivePath := "<unknown>"
	if d != nil {
		entryName = d.name
		archivePath = d.archivePath
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("处理文件 %q（来自 %q）时发生内部错误: %v", entryName, archivePath, recovered)
		}
	}()

	if d == nil {
		return errors.New("解包文件项不能为空")
	}
	file, err := os.Open(d.sourcePath)
	if err != nil {
		return errors.Errorf("打开 wxapkg 文件 '%s' 失败: %v", d.sourcePath, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return errors.Errorf("读取 wxapkg 文件 '%s' 信息失败: %v", d.sourcePath, err)
	}
	if info.IsDir() {
		return errors.Errorf("wxapkg 源文件 '%s' 是目录", d.sourcePath)
	}
	end := uint64(d.offset) + uint64(d.size)
	if end > uint64(info.Size()) {
		return errors.Errorf("文件 '%s' 的数据范围超出源文件边界", d.name)
	}

	if err := os.MkdirAll(filepath.Dir(d.savePath), os.ModePerm); err != nil {
		return errors.Errorf("创建目录 '%s' 失败: %v", filepath.Dir(d.savePath), err)
	}

	section := io.NewSectionReader(file, int64(d.offset), int64(d.size))
	ext := strings.ToLower(filepath.Ext(d.name))
	beautify := (ext == ".json" && u.options.EnableJsonBeautify) ||
		(ext == ".html" && u.options.EnableHtmlBeautify) ||
		(ext == ".js" && u.options.EnableJsBeautify)
	if beautify {
		data, err := io.ReadAll(section)
		if err != nil {
			return errors.Errorf("读取文件 '%s' 失败: %v", d.name, err)
		}
		if int64(len(data)) != int64(d.size) {
			return errors.Errorf("读取文件 '%s' 时遇到意外的文件末尾", d.name)
		}
		switch ext {
		case ".json":
			data = PrettyJson(data)
		case ".html":
			data = PrettyHtml(data)
		case ".js":
			data = PrettyJavaScript(data)
		}
		if err := os.WriteFile(d.savePath, data, 0600); err != nil {
			return errors.Errorf("写入文件 '%s' 失败: %v", d.savePath, err)
		}
		return nil
	}

	output, err := os.OpenFile(d.savePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Errorf("打开输出文件 '%s' 失败: %v", d.savePath, err)
	}
	written, copyErr := io.Copy(output, section)
	closeErr := output.Close()
	if copyErr != nil {
		return errors.Errorf("写入文件 '%s' 失败: %v", d.savePath, copyErr)
	}
	if closeErr != nil {
		return errors.Errorf("关闭输出文件 '%s' 失败: %v", d.savePath, closeErr)
	}
	if written != int64(d.size) {
		return errors.Errorf("写入文件 '%s' 时遇到意外的文件末尾", d.savePath)
	}
	return nil
}

func (u *Unpacker) unpack(thread int, callback func(item *WxapkgItem)) bool {
	if callback == nil {
		callback = func(*WxapkgItem) {}
	}
	if thread < 1 {
		thread = 1
	}
	if thread > maxUnpackWorkers {
		thread = maxUnpackWorkers
	}

	chFiles := make(chan *WxapkgFileItemStructure, thread)
	done := make(chan struct{})
	var stopOnce sync.Once
	stop := func() {
		stopOnce.Do(func() {
			close(done)
		})
	}

	var hasError bool
	var callbackLocker sync.Mutex
	reportError := func(err error) {
		callbackLocker.Lock()
		defer callbackLocker.Unlock()

		u.locker.Lock()
		if hasError {
			u.locker.Unlock()
			return
		}
		hasError = true
		u.item.SetErrorState(err.Error())
		snapshot := *u.item
		u.locker.Unlock()
		stop()
		callback(&snapshot)
	}
	reportProgress := func(name string) {
		callbackLocker.Lock()
		defer callbackLocker.Unlock()

		u.locker.Lock()
		if hasError {
			u.locker.Unlock()
			return
		}
		u.item.IncreaseProgress(name)
		snapshot := *u.item
		u.locker.Unlock()
		callback(&snapshot)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(chFiles)
		for _, d := range u.files {
			select {
			case chFiles <- d:
			case <-done:
				return
			}
		}
	}()

	wg.Add(thread)
	for i := 0; i < thread; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case d, ok := <-chFiles:
					if !ok {
						return
					}
					u.locker.Lock()
					stopped := hasError
					u.locker.Unlock()
					if stopped {
						return
					}
					if err := u.unpackOne(d); err != nil {
						if d == nil {
							reportError(err)
						} else {
							reportError(errors.Errorf("解包小程序文件 '%s'（来自 '%s'）时出错: %v", d.name, d.archivePath, err))
						}
						return
					}
					reportProgress(d.name)
				}
			}
		}()
	}

	wg.Wait()
	u.locker.Lock()
	result := hasError
	u.locker.Unlock()
	return result
}
