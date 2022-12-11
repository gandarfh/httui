package resources

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/styles"
	"github.com/gandarfh/httui/pkg/utils"
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

type TagItem struct {
	title string
	desc  string
	width int
}

func (i TagItem) FilterValue() string { return i.title }

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
		if common.CurrTab == common.Tab_Tags {
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

	common.ListOfResources, _ = m.resources_repo.List(common.CurrTag.ID, m.filter)
	w := m.width - (m.width / 3) - 1

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

	str := fmt.Sprintf("%s %s %s", utils.AddWhiteSpace(i.name, 30, 26), utils.AddWhiteSpace(strings.ToUpper(i.method), 10, 10), i.endpoint)
	if len(str) > i.width-10 {
		str = utils.Truncate(str, i.width-10)
	}

	if index == m.Index() {
		if common.CurrTab == common.Tab_Resources {
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
