package requests

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/styles"
	"github.com/gandarfh/httui/pkg/utils"
)

func NewRequestList() list.Model {
	list := list.New(nil, RequestDelegate{}, 0, 0)

	list.SetShowStatusBar(false)
	list.SetShowPagination(false)
	list.SetFilteringEnabled(false)
	list.SetShowFilter(false)
	list.SetShowHelp(false)

	list.Styles.Title = titleStyle
	list.Styles.NoItems = noItemsStyle

	return list
}

var (
	item_border = lipgloss.HiddenBorder()
)

var (
	noItemsStyle = lipgloss.NewStyle().
			MarginLeft(2).MarginRight(12).
			Foreground(styles.DefaultTheme.SecondaryBorder)

	titleStyle = lipgloss.NewStyle().MarginTop(1).Bold(true)

	itemStyle = lipgloss.NewStyle().
			Border(item_border).
			BorderTop(false)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Border(item_border).
				BorderTop(false)
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
		if common.CurrTab == common.Tab_Requests {
			fmt.Fprint(
				w,
				selectedItemStyle.
					Foreground(styles.DefaultTheme.PrimaryText).
					Width(i.width).Render("> "+str),
			)
		} else {
			fmt.Fprint(
				w,
				selectedItemStyle.
					Foreground(styles.DefaultTheme.SecondaryText).
					Width(i.width).Render("> "+str),
			)
		}

	} else {
		fmt.Fprint(w, itemStyle.Width(i.width).Render("  "+str))
	}
}

func (m Model) RequestOfList() []list.Item {
	list := []list.Item{}
	common.ListOfRequests, _ = repositories.NewRequest().List(m.parentId, m.filter)

	w := (m.width / 7)

	for _, i := range common.ListOfRequests {
		list = append(list, RequestItem{i.Name, i.Type, w})
	}

	return list
}
