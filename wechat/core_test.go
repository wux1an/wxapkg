package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func testUnpacker(t *testing.T) *Unpacker {
	t.Helper()
	root := t.TempDir()
	item := &WxapkgItem{UnpackSavePath: root}
	options := &UnpackOptions{SavePath: root}
	return NewUnpacker(item, options)
}

func testArchive(t *testing.T, name string, body []byte) ([]byte, uint32) {
	t.Helper()
	var data bytes.Buffer
	write := func(value interface{}) {
		t.Helper()
		if err := binary.Write(&data, binary.BigEndian, value); err != nil {
			t.Fatalf("write archive field: %v", err)
		}
	}

	write(uint8(0xBE))
	write(uint32(0))
	write(uint32(0))
	write(uint32(0))
	write(uint8(0xED))
	write(uint32(1))
	nameBytes := []byte(name)
	write(uint32(len(nameBytes)))
	if _, err := data.Write(nameBytes); err != nil {
		t.Fatalf("write archive name: %v", err)
	}
	offset := uint32(data.Len() + 8)
	write(offset)
	write(uint32(len(body)))
	if _, err := data.Write(body); err != nil {
		t.Fatalf("write archive body: %v", err)
	}
	return data.Bytes(), offset
}

func TestIsDecryptedWxapkgFileRejectsShortData(t *testing.T) {
	u := &Unpacker{}
	short := make([]byte, 13)
	short[0] = 0xBE
	if u.isDecryptedWxapkgFile(short) {
		t.Fatal("13-byte input must not be accepted as a decrypted wxapkg")
	}

	valid := make([]byte, 14)
	valid[0] = 0xBE
	valid[13] = 0xED
	if !u.isDecryptedWxapkgFile(valid) {
		t.Fatal("valid wxapkg markers were not accepted")
	}
}

func TestDecryptWxapkgFileRejectsShortData(t *testing.T) {
	for _, size := range []int{0, 6, 1029} {
		if _, err := decryptWxapkgFile("wx1234567890abcdef", make([]byte, size)); err == nil {
			t.Fatalf("decrypting %d-byte input must return an error", size)
		}
	}
}

func TestAnalyzeRejectsOutOfBoundsEntry(t *testing.T) {
	u := testUnpacker(t)
	data, offset := testArchive(t, "/app.js", []byte("test"))
	binary.BigEndian.PutUint32(data[int(offset)-8:], offset+1)

	if _, err := u.analyze(data, "archive.wxapkg", "archive.wxapkg"); err == nil {
		t.Fatal("out-of-bounds entry was accepted")
	}
}

func TestAnalyzeAllowsLeadingSlashEntry(t *testing.T) {
	u := testUnpacker(t)
	data, _ := testArchive(t, "/app-config.json", []byte("{}"))

	files, err := u.analyze(data, "archive.wxapkg", "archive.wxapkg")
	if err != nil {
		t.Fatalf("analyze returned an error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one file, got %d", len(files))
	}
	if files[0].name != "app-config.json" {
		t.Fatalf("unexpected normalized name %q", files[0].name)
	}
	relative, err := filepath.Rel(u.item.UnpackSavePath, files[0].savePath)
	if err != nil {
		t.Fatalf("compute relative output path: %v", err)
	}
	if relative == ".." || len(relative) >= 3 && relative[:3] == ".."+string(filepath.Separator) {
		t.Fatalf("output path escaped the unpack directory: %s", files[0].savePath)
	}
}

func TestAnalyzeRejectsTraversalEntry(t *testing.T) {
	u := testUnpacker(t)
	data, _ := testArchive(t, "../escape.js", []byte("test"))

	if _, err := u.analyze(data, "archive.wxapkg", "archive.wxapkg"); err == nil {
		t.Fatal("directory traversal entry was accepted")
	}
}

func TestNormalizeArchiveEntryNameIsPlatformIndependent(t *testing.T) {
	rejected := []string{
		"../escape.js",
		"..\\escape.js",
		"C:\\escape.js",
		"C:/escape.js",
		"\\\\server\\share\\file.js",
		"//server/share/file.js",
		"/../escape.js",
		"/C:/escape.js",
	}
	for _, name := range rejected {
		if _, err := normalizeArchiveEntryName(name); err == nil {
			t.Errorf("path %q was accepted", name)
		}
	}

	accepted := map[string]string{
		"/app-config.json": "app-config.json",
		"pages\\index.js":  "pages/index.js",
	}
	for name, want := range accepted {
		got, err := normalizeArchiveEntryName(name)
		if err != nil {
			t.Errorf("path %q was rejected: %v", name, err)
		} else if got != want {
			t.Errorf("path %q normalized to %q, want %q", name, got, want)
		}
	}
}

func TestUnpackOneRejectsTruncatedSource(t *testing.T) {
	u := testUnpacker(t)
	data, _ := testArchive(t, "app.js", []byte("test"))
	sourcePath := filepath.Join(t.TempDir(), "source.wxapkg")
	if err := os.WriteFile(sourcePath, data, 0600); err != nil {
		t.Fatalf("write source: %v", err)
	}
	files, err := u.analyze(data, sourcePath, sourcePath)
	if err != nil {
		t.Fatalf("analyze returned an error: %v", err)
	}
	if err := os.WriteFile(sourcePath, data[:len(data)-1], 0600); err != nil {
		t.Fatalf("truncate source: %v", err)
	}

	if err := u.unpackOne(files[0]); err == nil {
		t.Fatal("truncated source was accepted")
	}
}

func TestUnpackReturnsAfterWorkerError(t *testing.T) {
	u := testUnpacker(t)
	root := u.item.UnpackSavePath
	u.files = []*WxapkgFileItemStructure{{
		name:        "missing.js",
		savePath:    filepath.Join(root, "missing.js"),
		sourcePath:  filepath.Join(root, "missing.wxapkg"),
		archivePath: filepath.Join(root, "missing.wxapkg"),
		size:        1,
	}}

	result := make(chan bool, 1)
	go func() {
		result <- u.unpack(20, func(*WxapkgItem) {})
	}()

	select {
	case hasError := <-result:
		if !hasError {
			t.Fatal("worker failure was not propagated")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("unpack did not return after a worker failure")
	}
	if u.item.UnpackStatus != StatusTypeError {
		t.Fatalf("expected error status, got %q", u.item.UnpackStatus)
	}
}

func TestAppNameFromConfig(t *testing.T) {
	config := []byte("{\"global\":{\"window\":{\"navigationBarTitleText\":\"  Demo\\n App \"}}}")
	if got := appNameFromConfig(config); got != "Demo App" {
		t.Fatalf("unexpected app name %q", got)
	}
}

func TestAppNameFromConfigUsesExplicitIdentityField(t *testing.T) {
	config := []byte("{\"appName\":\"Official app\",\"global\":{\"window\":{\"navigationBarTitleText\":\"首页\"}}}")
	result := appNameFromConfigResult(config)
	if result.name != "Official app" {
		t.Fatalf("unexpected app name %q", result.name)
	}
	if result.source != appNameSourcePackageConfig {
		t.Fatalf("unexpected app name source %q", result.source)
	}
}

func TestAppNameFromConfigUsesGlobalTitle(t *testing.T) {
	config := []byte("{\"global\":{\"title\":\"Global app\",\"window\":{\"navigationBarTitleText\":\"首页\"}}}")
	result := appNameFromConfigResult(config)
	if result.name != "Global app" || result.source != appNameSourcePackageConfig {
		t.Fatalf("unexpected global title result: %+v", result)
	}
}

func TestAppNameFromConfigDoesNotUsePageTitle(t *testing.T) {
	config := []byte("{\"entryPagePath\":\"pages/index/index.html\",\"page\":{\"pages/index/index.html\":{\"window\":{\"navigationBarTitleText\":\"Entry title\"}}}}")
	if got := appNameFromConfig(config); got != "" {
		t.Fatalf("page title must not be used as app name, got %q", got)
	}
}

func TestAppNameFromConfigRejectsGenericTitle(t *testing.T) {
	config := []byte("{\"global\":{\"window\":{\"navigationBarTitleText\":\"\\u9996\\u9875\"}}}")
	if got := appNameFromConfig(config); got != "" {
		t.Fatalf("generic page title must not be used as app name, got %q", got)
	}
}

func TestReadWxapkgAppNameFromPlainPackage(t *testing.T) {
	config := []byte("{\"global\":{\"window\":{\"navigationBarTitleText\":\"Plain app\"}}}")
	data, _ := testArchive(t, "/app-config.json", config)
	path := filepath.Join(t.TempDir(), "main.wxapkg")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write package: %v", err)
	}

	got, err := readWxapkgAppName(path, "")
	if err != nil {
		t.Fatalf("read app name: %v", err)
	}
	if got != "Plain app" {
		t.Fatalf("unexpected app name %q", got)
	}
	result, err := readWxapkgAppNameResult(path, "")
	if err != nil {
		t.Fatalf("read app name result: %v", err)
	}
	if result.source != appNameSourceNavigationTitle {
		t.Fatalf("unexpected app name source %q", result.source)
	}
}

func TestReadWxapkgAppNameFromCodeCandidate(t *testing.T) {
	code := []byte("var meta={appId:\"wx1234567890abcdef\",appName:\"Code app\",appVersion:\"1.0\"};")
	data, _ := testArchive(t, "/app-service.js", code)
	path := filepath.Join(t.TempDir(), "main.wxapkg")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write package: %v", err)
	}

	result, err := readWxapkgAppNameResult(path, "")
	if err != nil {
		t.Fatalf("read code candidate: %v", err)
	}
	if result.name != "Code app" || result.source != appNameSourceCodeCandidate {
		t.Fatalf("unexpected code candidate result: %+v", result)
	}
}

func TestReadWxapkgAppNameFromAuthFloatInfoCandidate(t *testing.T) {
	code := []byte("var meta={authFloatInfo:{name:\"Auth app\"}};")
	data, _ := testArchive(t, "/app-service.js", code)
	path := filepath.Join(t.TempDir(), "main.wxapkg")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write package: %v", err)
	}

	result, err := readWxapkgAppNameResult(path, "")
	if err != nil {
		t.Fatalf("read auth candidate: %v", err)
	}
	if result.name != "Auth app" || result.source != appNameSourceCodeCandidate {
		t.Fatalf("unexpected auth candidate result: %+v", result)
	}
}

func TestAppNameForPathPrefersLocalMetadata(t *testing.T) {
	const wxid = "wx1234567890abcdef"
	root := t.TempDir()
	packagesRoot := filepath.Join(root, "applet", "packages", wxid)
	versionRoot := filepath.Join(packagesRoot, "939")
	localRoot := filepath.Join(root, "applet", "local", wxid)
	if err := os.MkdirAll(versionRoot, 0700); err != nil {
		t.Fatalf("create package directory: %v", err)
	}
	if err := os.MkdirAll(localRoot, 0700); err != nil {
		t.Fatalf("create local directory: %v", err)
	}

	packageData, _ := testArchive(t, "/app-config.json", []byte("{\"appName\":\"Package app\"}"))
	if err := os.WriteFile(filepath.Join(versionRoot, "__APP__.wxapkg"), packageData, 0600); err != nil {
		t.Fatalf("write package: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localRoot, "appinfo.json"), []byte("{\"appName\":\"Local app\"}"), 0600); err != nil {
		t.Fatalf("write local metadata: %v", err)
	}

	result := appNameForPathResult(packagesRoot, wxid)
	if result.name != "Local app" || result.source != appNameSourceLocalMetadata {
		t.Fatalf("unexpected local metadata result: %+v", result)
	}
}

func TestFindMainWxapkgUsesLatestVersion(t *testing.T) {
	const wxid = "wx1234567890abcdef"
	root := t.TempDir()
	packagesRoot := filepath.Join(root, "applet", "packages", wxid)
	oldRoot := filepath.Join(packagesRoot, "001")
	newRoot := filepath.Join(packagesRoot, "002")
	if err := os.MkdirAll(oldRoot, 0700); err != nil {
		t.Fatalf("create old package directory: %v", err)
	}
	if err := os.MkdirAll(newRoot, 0700); err != nil {
		t.Fatalf("create new package directory: %v", err)
	}
	oldData, _ := testArchive(t, "/app-config.json", []byte("{\"appName\":\"Old app\"}"))
	newData, _ := testArchive(t, "/app-config.json", []byte("{\"appName\":\"New app\"}"))
	oldPath := filepath.Join(oldRoot, "__APP__.wxapkg")
	newPath := filepath.Join(newRoot, "__APP__.wxapkg")
	if err := os.WriteFile(oldPath, oldData, 0600); err != nil {
		t.Fatalf("write old package: %v", err)
	}
	if err := os.WriteFile(newPath, newData, 0600); err != nil {
		t.Fatalf("write new package: %v", err)
	}
	if err := os.Chtimes(oldPath, time.Unix(100, 0), time.Unix(100, 0)); err != nil {
		t.Fatalf("set old package time: %v", err)
	}
	if err := os.Chtimes(newPath, time.Unix(200, 0), time.Unix(200, 0)); err != nil {
		t.Fatalf("set new package time: %v", err)
	}

	if got := findMainWxapkg(packagesRoot); got != newPath {
		t.Fatalf("selected package %q, want %q", got, newPath)
	}
	result := appNameForPathResult(packagesRoot, wxid)
	if result.name != "New app" {
		t.Fatalf("selected package returned app name %q", result.name)
	}
}

func TestReadWxapkgAppNameFromEncryptedPackage(t *testing.T) {
	const wxid = "wx1234567890abcdef"
	config := []byte("{\"global\":{\"window\":{\"navigationBarTitleText\":\"Encrypted app\"}}}")
	plain, _ := testArchive(t, "/app-config.json", config)
	encrypted := encryptTestWxapkg(t, plain, wxid)
	path := filepath.Join(t.TempDir(), "main.wxapkg")
	if err := os.WriteFile(path, encrypted, 0600); err != nil {
		t.Fatalf("write encrypted package: %v", err)
	}

	got, err := readWxapkgAppName(path, wxid)
	if err != nil {
		t.Fatalf("read encrypted app name: %v", err)
	}
	if got != "Encrypted app" {
		t.Fatalf("unexpected app name %q", got)
	}
}

func TestScanWxapkgItemIncludesAppName(t *testing.T) {
	config := []byte("{\"global\":{\"window\":{\"navigationBarTitleText\":\"Scanned app\"}}}")
	data, _ := testArchive(t, "/app-config.json", config)
	dir := filepath.Join(t.TempDir(), "wxd6d308853899ff77")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("create package directory: %v", err)
	}
	path := filepath.Join(dir, "__APP__.wxapkg")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write package: %v", err)
	}

	items, err := ScanWxapkgItem(path, false)
	if err != nil {
		t.Fatalf("scan package: %v", err)
	}
	if len(items) != 1 || items[0].AppName != "Scanned app" {
		t.Fatalf("unexpected scan result: %+v", items)
	}
	if items[0].AppNameSource != appNameSourceNavigationTitle {
		t.Fatalf("unexpected scan app name source %q", items[0].AppNameSource)
	}
	if items[0].WxId != "wxd6d308853899ff77" {
		t.Fatalf("wxid was not inferred from path: %q", items[0].WxId)
	}
}

func encryptTestWxapkg(t *testing.T, plain []byte, wxid string) []byte {
	t.Helper()
	firstBlock := make([]byte, encryptedBlockSize)
	copy(firstBlock, plain)
	key := pbkdf2.Key([]byte(wxid), []byte("saltiest"), 1000, 32, sha1.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("create test cipher: %v", err)
	}
	encryptedBlock := make([]byte, encryptedBlockSize)
	cipher.NewCBCEncrypter(block, []byte("the iv: 16 bytes")).CryptBlocks(encryptedBlock, firstBlock)

	result := make([]byte, encryptedHeaderSize, encryptedHeaderSize+encryptedBlockSize+len(plain))
	result = append(result, encryptedBlock...)
	xorKey := wxid[len(wxid)-2]
	if len(plain) > encryptedLogicalPrefixSize {
		for _, value := range plain[encryptedLogicalPrefixSize:] {
			result = append(result, value^xorKey)
		}
	}
	return result
}
