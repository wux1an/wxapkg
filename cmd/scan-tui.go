package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/wux1an/wxapkg/util"
	"regexp"
	"strconv"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("240"))

type scanTui struct {
	table    table.Model
	raw      []util.WxidInfo
	selected *util.WxidInfo
	progress progress.Model
}

func newScanTui(wxidInfo []util.WxidInfo) *scanTui {
	var prog = progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

	var rows = make([]table.Row, 0, len(wxidInfo))
	for _, info := range wxidInfo {
		rows = append(rows, []string{
			info.Nickname,
			info.PrincipalName,
			info.Description,
		})
	}

	var title = color.New(color.FgMagenta, color.Bold).Sprint
	columns := []table.Column{
		{Title: title("Name"), Width: 20},
		{Title: title("Developer"), Width: 30},
		{Title: title("Description"), Width: 40},
	}
	prog.Width = 0
	for _, c := range columns {
		prog.Width += c.Width
	}
	prog.Width += len(columns) * 2

	var height = 10
	if len(rows) < height {
		height = len(rows)
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)
	t.Rows()

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &scanTui{
		table:    t,
		raw:      wxidInfo,
		progress: prog,
	}
}

func (s *scanTui) Init() tea.Cmd {
	return nil
}

func (s *scanTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if s.table.Focused() {
				s.table.Blur()
			} else {
				s.table.Focus()
			}
		case "q", "ctrl+c":
			return s, tea.Quit
		case "enter":
			s.selected = &s.raw[s.table.Cursor()]
			return s, tea.Quit
		}
	}
	s.table, cmd = s.table.Update(msg)
	return s, cmd
}

func (s *scanTui) renderProgress() string {
	var prog = s.progress.ViewAs(float64(s.table.Cursor()+1) / float64(len(s.raw)))
	var p = regexp.MustCompile(`\d{1,3}%`).FindString(prog)
	format := "%" + strconv.Itoa(len(p)) + "s"
	newStr := fmt.Sprintf(format, fmt.Sprintf("%d/%d", s.table.Cursor()+1, len(s.raw)))
	prog = strings.Replace(prog, p, newStr, 1)

	return prog
}

func (s *scanTui) renderDetail() string {
	var result = ""

	var info = s.raw[s.table.Cursor()]

	var link = color.New(color.Italic, color.Underline).Sprint
	var title = color.New(color.FgMagenta, color.Bold).Sprint
	var content = color.CyanString

	if info.Error != "" {
		result += title("  error: ") + color.RedString(info.Error) + "\n"
	}

	if info.Error == "" {
		result += title("  wxid: ") + content(info.Wxid) + "\n"
		result += title("  Name: ") + content(info.Nickname) + "\n"
		result += title("  Developer: ") + content(info.PrincipalName) + "\n"
		result += title("  Description: ") + content(info.Description) + "\n"
	}

	result += title("  Location: ") + content(link(info.Location)) + "\n"

	if info.Error == "" {
		result += title("  Avatar: ") + content(link(info.Avatar)) + "\n"
	}

	result += title("  All information see '") + content(".\\"+util.CachePath) + title("'")

	return result
}

func (s *scanTui) renderHelp() string {
	return help.New().ShortHelpView([]key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "unpack"),
		),
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "exit"),
		),
	})
}

func (s *scanTui) View() string {
	var result = ""
	result += "" + s.renderProgress() + "\n"
	result += baseStyle.Render(s.table.View()) + "\n"
	result += s.renderDetail() + "\n"
	result += "\n  " + s.renderHelp() + "\n"

	return result
}
