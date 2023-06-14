package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
	"github.com/wux1an/wxapkg/util"
	"github.com/yosssi/gohtml"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
)

var programName = filepath.Base(os.Args[0])
var unpackCmd = &cobra.Command{
	Use:     "unpack",
	Short:   "Decrypt wechat mini program",
	Example: "  " + programName + "unpack -o unpack -r \"%USERPROFILE%\\Documents\\WeChat Files\\Applet\\wx00000000000000\"",
	Run: func(cmd *cobra.Command, args []string) {
		root, _ := cmd.Flags().GetString("root")
		output, _ := cmd.Flags().GetString("output")
		thread, _ := cmd.Flags().GetInt("thread")

		wxid, err := parseWxid(root)
		util.Fatal(err)

		dirs, err := os.ReadDir(root)
		util.Fatal(err)

		color.Cyan("[+] unpack root '%s' with %d threads\n", root, thread)

		var allFileCount = 0
		for _, subDir := range dirs {
			subOutput := filepath.Join(output, subDir.Name())

			files, err := scanFiles(filepath.Join(root, subDir.Name()))
			util.Fatal(err)

			for _, file := range files {
				var decryptedData = decryptFile(wxid, file)
				fileCount, err := unpack(decryptedData, subOutput, thread)
				util.Fatal(err)
				allFileCount += fileCount

				rel, _ := filepath.Rel(filepath.Dir(root), file)
				color.Yellow("\r[+] unpacked %5d files from '%s'", fileCount, rel)
			}
		}

		color.Cyan("[+] all %d files saved to '%s'\n", allFileCount, output)
		if len(args) == 2 && "detailFilePath" == args[0] {
			color.Cyan("[+] mini program detail info saved to '%s'\n", args[1])
		}

		color.Cyan("[+] extension statistics:\n")

		var keys [][]interface{}
		for k, v := range exts {
			keys = append(keys, []interface{}{k, v})
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i][1].(int) > keys[j][1].(int)
		})

		for _, kk := range keys {
			color.Cyan("  - %-5s %5d\n", kk[0], kk[1])
		}
	},
}

type wxapkgFile struct {
	nameLen uint32
	name    []byte
	offset  uint32
	size    uint32
}

func unpack(decryptedData []byte, unpackRoot string, thread int) (int, error) {
	var f = bytes.NewReader(decryptedData)

	// Read header
	var (
		firstMark       uint8
		info1           uint32
		indexInfoLength uint32
		bodyInfoLength  uint32
		lastMark        uint8
	)
	_ = binary.Read(f, binary.BigEndian, &firstMark)
	_ = binary.Read(f, binary.BigEndian, &info1)
	_ = binary.Read(f, binary.BigEndian, &indexInfoLength)
	_ = binary.Read(f, binary.BigEndian, &bodyInfoLength)
	_ = binary.Read(f, binary.BigEndian, &lastMark)

	if firstMark != 0xBE || lastMark != 0xED {
		return 0, errors.New("failed to unpack, it's not a valid wxapkg file")
	}

	var fileCount uint32
	_ = binary.Read(f, binary.BigEndian, &fileCount)

	// Read index
	var fileList = make([]*wxapkgFile, fileCount)
	for i := uint32(0); i < fileCount; i++ {
		data := &wxapkgFile{}
		_ = binary.Read(f, binary.BigEndian, &data.nameLen)

		if data.nameLen > 10<<20 { // 10 MB
			return 0, errors.New("invalid decrypted wxapkg file")
		}

		data.name = make([]byte, data.nameLen)
		_, _ = io.ReadAtLeast(f, data.name, int(data.nameLen))
		_ = binary.Read(f, binary.BigEndian, &data.offset)
		_ = binary.Read(f, binary.BigEndian, &data.size)

		fileList[i] = data
	}

	// Save files
	var chFiles = make(chan *wxapkgFile)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, d := range fileList {
			chFiles <- d
		}
		close(chFiles)
	}()

	wg.Add(thread)
	var locker = sync.Mutex{}
	var count = 0
	var colorPrint = color.New()
	for i := 0; i < thread; i++ {
		go func() {
			defer wg.Done()

			for d := range chFiles {
				d.name = []byte(filepath.Join(unpackRoot, string(d.name)))
				outputFilePath := string(d.name)
				dir := filepath.Dir(outputFilePath)

				err := os.MkdirAll(dir, os.ModePerm)
				util.Fatal(err)

				data := decryptedData[d.offset : d.offset+d.size]

				beautify := fileBeautify(outputFilePath, data)
				err = os.WriteFile(outputFilePath, beautify, 0600)
				util.Fatal(err)

				locker.Lock()
				count++
				_, _ = colorPrint.Print(color.GreenString("\runpack %d/%d", count, fileCount))
				locker.Unlock()
			}
		}()
	}

	wg.Wait()

	return int(fileCount), nil
}

var exts = make(map[string]int)
var extsLocker = sync.Mutex{}

func fileBeautify(name string, data []byte) (result []byte) {
	defer func() {
		if err := recover(); err != nil {
			result = data
		}
	}()

	result = data
	var ext = filepath.Ext(name)

	extsLocker.Lock()
	exts[ext] = exts[ext] + 1
	extsLocker.Unlock()

	switch ext {
	case ".json":
		result = pretty.Pretty(data)
	case ".html": // todo beautify js code in html
		result = gohtml.FormatBytes(bytes.TrimSpace(data)) // remove leading whitespace
	case ".js":
		var code = string(bytes.TrimSpace(data)) // remove leading whitespace
		options := jsbeautifier.DefaultOptions()
		beautify, err := jsbeautifier.Beautify(&code, options)
		if err == nil {
			result = []byte(beautify)
		}
	}

	return result
}

func parseWxid(root string) (string, error) {
	var regAppId = regexp.MustCompile(`(wx[0-9a-f]{16})`)
	if !regAppId.MatchString(filepath.Base(root)) {
		return "", errors.New("the path is not a mini program path")
	}

	return regAppId.FindStringSubmatch(filepath.Base(root))[1], nil
}

func scanFiles(root string) ([]string, error) {
	paths, err := util.GetDirAllFilePaths(root, "", ".wxapkg")
	util.Fatal(err)

	if len(paths) == 0 {
		return nil, errors.New(fmt.Sprintf("no '.wxapkg' file found in '%s'", root))
	}

	return paths, nil
}

func decryptFile(wxid, wxapkgPath string) []byte {
	var (
		salt = "saltiest"
		iv   = "the iv: 16 bytes"
	)

	dataByte, err := os.ReadFile(wxapkgPath)
	if err != nil {
		log.Fatal(err)
	}

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

	return originData
}

func init() {
	RootCmd.AddCommand(unpackCmd)

	var homeDir, _ = os.UserHomeDir()
	var defaultRoot = filepath.Join(homeDir, "Documents/WeChat Files/Applet", "wx00000000000000")

	unpackCmd.Flags().StringP("root", "r", "", "the mini progress path you want to decrypt, see: "+defaultRoot)
	unpackCmd.Flags().StringP("output", "o", "unpack", "the output path to save result")
	unpackCmd.Flags().IntP("thread", "n", 30, "the thread number")
	_ = unpackCmd.MarkFlagRequired("root")
}
