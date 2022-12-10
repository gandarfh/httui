package envs

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
	width    int
	height   int
	env_list list.Model
	env_repo *repositories.EnvsRepo
}

func New() Model {
	env_repo, _ := repositories.NewEnvs()

	list := list.New(nil, Delegate{}, 0, 0)
	list.Title = "All Environments"
	list.SetShowPagination(false)
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)

	list.Styles.Title = titleStyle
	list.Styles.NoItems = noItemsStyle

	return Model{env_repo: env_repo, env_list: list}
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
		m.env_list.SetHeight(14)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "d":
			index := m.env_list.Index()
			m.env_repo.Delete(common.ListOfEnvs[index].ID)

		case "c":
			data := repositories.Env{}
			term := terminal.NewPreview(&data)
			return m, tea.Batch(term.OpenVim("Create"))

		case "r":
			index := m.env_list.Index()
			common.CurrEnv = common.ListOfEnvs[index]

			term := terminal.NewPreview(&common.CurrEnv)
			return m, tea.Batch(term.OpenVim("Update"))

		}

	case terminal.Finish:
		switch msg.Category {
		case "Update":
			data := repositories.Env{}
			msg.Preview.Execute(&data)
			m.env_repo.Update(&common.CurrEnv, &data)

		case "Create":
			data := repositories.Env{}
			msg.Preview.Execute(&data)
			m.env_repo.Create(&data)
		}

		defer msg.Preview.Close()
		if msg.Err != nil {
			return m, nil
		}
	}

	m.env_list.SetItems(m.ItemsOfList())
	m.env_list, cmd = m.env_list.Update(msg)
	cmds = append(cmds, cmd)

	return m, nil
}

func (m Model) View() string {
	return m.env_list.View()
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
	key   string
	value string
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

  value := utils.Truncate(i.value, 20)

	if index == m.Index() {
		fmt.Fprint(w,
			lipgloss.NewStyle().
				Border(item_border).
				BorderTop(false).
				BorderForeground(styles.DefaultTheme.SecondaryBorder).
				Foreground(styles.DefaultTheme.PrimaryText).
				Render(
					fmt.Sprintf(
						"> %s %s",
						lipgloss.NewStyle().
							Bold(true).
							Render(utils.AddWhiteSpace(i.key, 30, 27)),
						lipgloss.NewStyle().
							Foreground(styles.DefaultTheme.SecondaryText).
							Render(value),
					),
				),
		)
	} else {
		fmt.Fprint(w, itemStyle.Render(fmt.Sprintf("  %s %s", lipgloss.NewStyle().Bold(true).Render(utils.AddWhiteSpace(i.key, 30, 27)), lipgloss.NewStyle().Render(value))))
	}
}

func (m Model) ItemsOfList() []list.Item {
	list := []list.Item{}
	common.ListOfEnvs, _ = m.env_repo.List()

	w := m.width - (m.width / 10)

	for _, i := range common.ListOfEnvs {
		list = append(list, Item{i.Key, i.Value, w})
	}

	return list
}
