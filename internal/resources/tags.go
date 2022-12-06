package resources

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/maid-san/pkg/common"
	"github.com/gandarfh/maid-san/pkg/styles"
	"github.com/gandarfh/maid-san/pkg/utils"
)

func NewTagList() list.Model {
	list := list.New(nil, TagDelegate{}, 0, 0)

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
	item_border = lipgloss.Border{
		Top:         " ",
		Bottom:      " ",
		Left:        " ",
		Right:       " ",
		TopLeft:     " ",
		TopRight:    " ",
		BottomLeft:  " ",
		BottomRight: " ",
	}
)

var (
	noItemsStyle = lipgloss.NewStyle().MarginLeft(2).
			Foreground(styles.DefaultTheme.SecondaryBorder)
	titleStyle = lipgloss.NewStyle().Bold(true)
	itemStyle  = lipgloss.NewStyle().
			Border(item_border).
			BorderTop(false).
			BorderForeground(styles.DefaultTheme.SecondaryBorder)
	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Border(item_border).
				BorderTop(false).
				BorderForeground(styles.DefaultTheme.SecondaryBorder).
				Foreground(styles.DefaultTheme.PrimaryText)
)

type TagItem struct {
	title string
	desc  string
	width int
}

func (i TagItem) FilterValue() string { return " " }

type TagDelegate struct{}

func (d TagDelegate) Height() int                               { return 1 }
func (d TagDelegate) Spacing() int                              { return 0 }
func (d TagDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d TagDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(TagItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.title)

	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Width(i.width).Render("> "+str))
	} else {
		fmt.Fprint(w, itemStyle.Width(i.width).Render("  "+str))
	}
}

func (m Model) TagsOfList() []list.Item {
	list := []list.Item{}
	common.ListOfTags, _ = m.tags_repo.List(common.CurrWorkspace.ID)

	w := (m.width / 7)

	for _, i := range common.ListOfTags {
		list = append(list, TagItem{i.Name, i.Description, w})
	}

	return list
}

func NewResourceList() list.Model {
	list := list.New(nil, ResourceDelegate{}, 0, 0)

	list.SetShowStatusBar(false)
	list.SetShowPagination(false)
	list.SetFilteringEnabled(false)
	list.SetShowFilter(false)
	list.SetShowHelp(false)

	list.Styles.Title = titleStyle
	list.Styles.NoItems = noItemsStyle

	return list
}

func (m Model) ResourcesOfList() []list.Item {
	list := []list.Item{}

	common.ListOfResources, _ = m.resources_repo.List(common.CurrTag.ID)
	w := m.width - (m.width / 2)

	for _, i := range common.ListOfResources {
		list = append(list, ResourceItem{i.Name, i.Method, i.Endpoint, w})
	}

	return list
}

type ResourceItem struct {
	name     string
	method   string
	endpoint string
	width    int
}

func (i ResourceItem) FilterValue() string { return " " }

type ResourceDelegate struct{}

func (d ResourceDelegate) Height() int                               { return 1 }
func (d ResourceDelegate) Spacing() int                              { return 0 }
func (d ResourceDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d ResourceDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ResourceItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s %s %s", utils.AddWhiteSpace(i.name, 30, 26), utils.AddWhiteSpace(strings.ToUpper(i.method), 10, 10), utils.Truncate(i.endpoint, 18))

	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Width(i.width).Render("> "+str))
	} else {
		fmt.Fprint(w, itemStyle.Width(i.width).Render("  "+str))
	}
}
