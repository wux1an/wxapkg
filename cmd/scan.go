package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wux1an/wxapkg/util"
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

		var wxidInfos = make([]util.WxidInfo, 0, len(files))
		for _, file := range files {
			if !file.IsDir() || !regAppId.MatchString(file.Name()) {
				continue
			}

			var wxid = regAppId.FindStringSubmatch(file.Name())[1]
			info, err := util.WxidQuery.Query(wxid)
			info.Location = filepath.Join(root, file.Name())
			info.Wxid = wxid
			if err != nil {
				info.Error = fmt.Sprintf("%v", err)
			}

			wxidInfos = append(wxidInfos, info)
		}

		var tui = newScanTui(wxidInfos)
		if _, err := tea.NewProgram(tui, tea.WithAltScreen()).Run(); err != nil {
			color.Red("Error running program: %v", err)
			os.Exit(1)
		}

		if tui.selected == nil {
			return
		}

		output := tui.selected.Wxid
		_ = unpackCmd.Flags().Set("root", tui.selected.Location)
		_ = unpackCmd.Flags().Set("output", output)
		detailFilePath := filepath.Join(output, "detail.json")
		unpackCmd.Run(unpackCmd, []string{"detailFilePath", detailFilePath})
		_ = os.WriteFile(detailFilePath, []byte(tui.selected.Json()), 0600)
	},
}

func init() {
	RootCmd.AddCommand(scanCmd)

	var homeDir, _ = os.UserHomeDir()
	var defaultRoot = filepath.Join(homeDir, "Documents/WeChat Files/Applet")

	scanCmd.Flags().StringP("root", "r", defaultRoot, "the mini app path")
}
