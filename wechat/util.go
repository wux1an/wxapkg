package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	uuidlib "github.com/google/uuid"
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

	dk := pbkdf2.Key([]byte(wxid), []byte(salt), 1000, 32, sha1.New)
	block, _ := aes.NewCipher(dk)
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	originData := make([]byte, 1024)
	blockMode.CryptBlocks(originData, dataByte[6:1024+6])

	afData := make([]byte, len(dataByte)-1024-6) // remove first 6 + 1024 byte
	var xorKey = byte(0x66)
	if len(wxid) >= 2 {
		xorKey = wxid[len(wxid)-2]
	}
	for i, b := range dataByte[1024+6:] { // from 6 + 1024 byte
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
		return []WxapkgItem{
			{
				UUID:           uuid(),
				WxId:           "unknown",
				Location:       path,
				Size:           getPathSize(path),
				IsDir:          stat.IsDir(),
				LastModifyTime: stat.ModTime().Unix(),
			},
		}, nil
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

				item := WxapkgItem{
					UUID:       uuid(),
					WxId:       dirName,
					Title:      "",
					Location:   localPath,
					EncryptKey: dirName,
					Size:       getPathSize(localPath),
					IsDir:      true,
				}
				if info, err := os.Stat(localPath); err == nil {
					item.LastModifyTime = info.ModTime().Unix()
				}

				item.Title, err = SimpleGetAppConfig(&item)
				if err != nil {
					fmt.Printf("获取小程序配置失败，wxid: %s, 错误: %v\n", dirName, err)
				}

				result = append(result, item)
			}
		}
	}

	return result, nil
}
