package workspaces

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/maid-san/internal/repositories"
	"github.com/gandarfh/maid-san/pkg/common"
	"github.com/gandarfh/maid-san/pkg/styles"
	"github.com/gandarfh/maid-san/pkg/terminal"
	"github.com/gandarfh/maid-san/pkg/utils"
)

type Model struct {
	width          int
	height         int
	workspace_list list.Model
	workspace_repo *repositories.WorkspacesRepo
	default_repo   *repositories.DefaultsRepo
	tags_repo      *repositories.TagsRepo
}

func New() Model {
	workspace_repo, _ := repositories.NewWorkspace()
	tags_repo, _ := repositories.NewTag()
	default_repo, _ := repositories.NewDefault()

	list := list.New(nil, Delegate{}, 0, 0)
	list.Title = "All workspaces"
	list.SetShowPagination(false)
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)

	list.Styles.Title = titleStyle
	list.Styles.NoItems = noItemsStyle

	return Model{
		workspace_repo: workspace_repo,
		tags_repo:      tags_repo,
		default_repo:   default_repo,
		workspace_list: list,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height - (msg.Height / 3)
		m.width = msg.Width
		m.workspace_list.SetHeight(msg.Height/2 - 2)

	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			index := m.workspace_list.Index()
			workspaces, _ := m.workspace_repo.List()
			common.CurrWorkspace = workspaces[index]

			m.workspace_repo.Delete(common.CurrWorkspace.ID)

		case "c":
			data := repositories.Workspace{}
			term := terminal.NewPreview(&data)
			return m, tea.Batch(term.OpenVim("Create"))

		case "r":
			index := m.workspace_list.Index()
			common.CurrWorkspace = common.ListOfWorkspaces[index]

			term := terminal.NewPreview(&common.CurrWorkspace)
			return m, tea.Batch(term.OpenVim("Update"))

		case "enter":
			index := m.workspace_list.Index()
			common.CurrWorkspace = common.ListOfWorkspaces[index]

			data := repositories.Default{
				WorkspaceId: common.CurrWorkspace.ID,
			}

			m.default_repo.Update(&data)

			return m, tea.Batch(
				common.SetPage(common.Page_Resource),
				common.SetTab(common.Tab_Tags),
				common.ListTags(common.CurrWorkspace.ID),
			)
		}

	case terminal.Finish:
		switch msg.Category {
		case "Update":
			data := repositories.Workspace{}
			msg.Preview.Execute(&data)
			m.workspace_repo.Update(&common.CurrWorkspace, &data)

		case "Create":
			data := repositories.Workspace{}
			msg.Preview.Execute(&data)
			m.workspace_repo.Create(&data)
		}

		defer msg.Preview.Close()
		if msg.Err != nil {
			return m, nil
		}
	}

	common.ListOfWorkspaces, _ = m.workspace_repo.List()
	m.workspace_list.SetItems(m.ItemsOfList())
	m.workspace_list, cmd = m.workspace_list.Update(msg)
	cmds = append(cmds, cmd)

	return m, nil
}

func (m Model) View() string {
	return m.workspace_list.View()
}
func NewTagList() list.Model {
	list := list.New(nil, Delegate{}, 0, 0)

	list.SetShowStatusBar(false)
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
	noItemsStyle = lipgloss.NewStyle().MarginLeft(2).
			Foreground(styles.DefaultTheme.SecondaryBorder)
	titleStyle = lipgloss.NewStyle().MarginTop(1).Bold(true)
	itemStyle  = lipgloss.NewStyle().
			Border(item_border).
			BorderTop(false).
			BorderForeground(styles.DefaultTheme.SecondaryBorder)
)

type Item struct {
	title string
	host  string
	width int
}

func (i Item) FilterValue() string { return "" }

type Delegate struct{}

func (d Delegate) Height() int                               { return 1 }
func (d Delegate) Spacing() int                              { return 0 }
func (d Delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d Delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprint(w,
			lipgloss.NewStyle().
				Border(item_border).
				BorderTop(false).
				BorderForeground(styles.DefaultTheme.SecondaryBorder).
				Foreground(styles.DefaultTheme.PrimaryText).
				Render(
					fmt.Sprintf(
						"%s %s",
						lipgloss.NewStyle().
							Bold(true).
							Render("> "+utils.AddWhiteSpace(i.title, 30, 27)),
						lipgloss.NewStyle().
							Foreground(styles.DefaultTheme.SecondaryText).
							Render(i.host),
					),
				),
		)
	} else {
		fmt.Fprint(w, itemStyle.Render(fmt.Sprintf("  %s %s", lipgloss.NewStyle().Bold(true).Render(utils.AddWhiteSpace(i.title, 30, 27)), lipgloss.NewStyle().Render(i.host))))
	}
}

func (m Model) ItemsOfList() []list.Item {
	list := []list.Item{}
	common.ListOfWorkspaces, _ = m.workspace_repo.List()

	w := m.width - (m.width / 10)

	for _, i := range common.ListOfWorkspaces {
		list = append(list, Item{i.Name, i.Host, w})
	}

	return list
}
