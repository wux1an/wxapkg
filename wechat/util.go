package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	uuidlib "github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
)

func getPathSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	if !info.IsDir() {
		return info.Size()
	}
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

func uuid() string {
	newUUID, _ := uuidlib.NewRandom()
	return newUUID.String()
}

// ListFilesWithExtension extension eg: .wxapkg
func ListFilesWithExtension(dir, extension string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, extension) {
			paths = append(paths, path)
		}
		return nil
	})

	return paths, err
}

func decryptWxapkgFile(wxid string, dataByte []byte) ([]byte, error) {
	var (
		salt = "saltiest"
		iv   = "the iv: 16 bytes"
	)
	const encryptedHeaderSize = 6
	const encryptedBlockSize = 1024
	if len(dataByte) < encryptedHeaderSize+encryptedBlockSize {
		return nil, errors.Errorf("加密 wxapkg 文件长度 %d 小于最小长度 %d", len(dataByte), encryptedHeaderSize+encryptedBlockSize)
	}

	dk := pbkdf2.Key([]byte(wxid), []byte(salt), 1000, 32, sha1.New)
	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, errors.Errorf("创建 wxapkg 解密器失败: %v", err)
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	originData := make([]byte, 1024)
	blockMode.CryptBlocks(originData, dataByte[encryptedHeaderSize:encryptedHeaderSize+encryptedBlockSize])

	afData := make([]byte, len(dataByte)-encryptedBlockSize-encryptedHeaderSize) // remove first 6 + 1024 byte
	var xorKey = byte(0x66)
	if len(wxid) >= 2 {
		xorKey = wxid[len(wxid)-2]
	}
	for i, b := range dataByte[encryptedHeaderSize+encryptedBlockSize:] { // from 6 + 1024 byte
		afData[i] = b ^ xorKey
	}

	originData = append(originData[:1023], afData...)

	return originData, nil
}

func ScanWxapkgItem(path string, scan bool) ([]WxapkgItem, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() || !scan {
		wxid := wxIDFromPath(path)
		itemWxID := wxid
		if itemWxID == "" {
			itemWxID = "unknown"
		}
		appName := appNameForPathResult(path, wxid)
		return []WxapkgItem{{
			UUID:           uuid(),
			WxId:           itemWxID,
			Location:       path,
			AppName:        appName.name,
			AppNameSource:  appName.source,
			EncryptKey:     wxid,
			Size:           getPathSize(path),
			IsDir:          stat.IsDir(),
			LastModifyTime: stat.ModTime().Unix(),
		}}, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []WxapkgItem

	for _, entry := range entries {
		if entry.IsDir() {
			dirName := entry.Name()
			if reWxId.MatchString(dirName) {
				localPath := filepath.Join(path, dirName)
				appName := appNameForPathResult(localPath, dirName)

				item := WxapkgItem{
					UUID:          uuid(),
					WxId:          dirName,
					Location:      localPath,
					EncryptKey:    dirName,
					AppName:       appName.name,
					AppNameSource: appName.source,
					Size:          getPathSize(localPath),
					IsDir:         true,
				}
				if info, err := os.Stat(localPath); err == nil {
					item.LastModifyTime = info.ModTime().Unix()
				}

				result = append(result, item)
			}
		}
	}

	return result, nil
}
