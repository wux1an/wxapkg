package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the wechat mini program",
	Run: func(cmd *cobra.Command, args []string) {
		root, err := cmd.Flags().GetString("root")
		if err != nil {
			color.Red("%v", err)
			return
		}

		var regAppId = regexp.MustCompile(`(wx[0-9a-f]{16})`)

		var files []os.DirEntry
		if files, err = os.ReadDir(root); err != nil {
			color.Red("%v", err)
			return
		}

		for _, file := range files {
			if !file.IsDir() || !regAppId.MatchString(file.Name()) {
				continue
			}

			var id = regAppId.FindStringSubmatch(file.Name())[1]
			fmt.Printf("%s %s\n", color.GreenString(id), color.CyanString(filepath.Join(root, file.Name()))) // todo scan file
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	var homeDir, _ = os.UserHomeDir()
	var defaultRoot = filepath.Join(homeDir, "Documents/WeChat Files/Applet")

	scanCmd.Flags().StringP("root", "r", defaultRoot, "the mini app path")
}
