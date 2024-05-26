package requests

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/styles"
	"github.com/gandarfh/httui/pkg/utils"
)

func NewRequestList() list.Model {
	component := list.New(nil, RequestDelegate{}, 0, 0)

	component.DisableQuitKeybindings()
	component.SetShowStatusBar(false)
	component.SetShowPagination(false)
	component.SetFilteringEnabled(false)
	component.SetShowFilter(false)
	component.SetShowHelp(false)

	component.Styles.Title = titleStyle
	component.Styles.NoItems = noItemsStyle

	return component
}

var (
	noItemsStyle = lipgloss.NewStyle().
			MarginLeft(2).MarginRight(12).
			Foreground(styles.DefaultTheme.SecondaryBorder)

	titleStyle = lipgloss.NewStyle().MarginTop(1).Bold(true)

	itemStyle = lipgloss.NewStyle().PaddingBottom(1)

	selectedItemStyle = lipgloss.NewStyle().PaddingBottom(1).Bold(true)
)

type RequestItem struct {
	title  string
	typeOf string
	width  int
}

func (i RequestItem) FilterValue() string { return i.title }

type RequestDelegate struct{}

func (d RequestDelegate) Height() int                               { return 1 }
func (d RequestDelegate) Spacing() int                              { return 0 }
func (d RequestDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d RequestDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(RequestItem)
	if !ok {
		return
	}

	prefix := ""

	switch i.typeOf {
	case "group":
		prefix = "-"
	case "request":
		prefix = "."
	}

	str := fmt.Sprintf("%s %s", prefix, utils.Truncate(i.title, i.width-10))

	if index == m.Index() {
		fmt.Fprint(
			w,
			selectedItemStyle.
				Foreground(styles.DefaultTheme.PrimaryText).
				Render("> "+str),
		)

	} else {
		fmt.Fprint(w, itemStyle.Render("  "+str))
	}
}
