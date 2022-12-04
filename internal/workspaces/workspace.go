package workspaces

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/maid-san/internal/repositories"
	"github.com/gandarfh/maid-san/pkg/common"
	"github.com/gandarfh/maid-san/pkg/terminal"
)

type Model struct {
	width          int
	hight          int
	workspace_list list.Model
	workspace_repo *repositories.WorkspacesRepo
}

func New() Model {
	repo, _ := repositories.NewWorkspace()

	list := list.New(nil, list.NewDefaultDelegate(), 34, 10)
	list.SetShowPagination(false)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)

	return Model{workspace_repo: repo, workspace_list: list}
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
		m.hight = msg.Height
		m.width = msg.Width
		m.workspace_list.SetSize(msg.Width-10, msg.Height-10)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
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
			workspaces, _ := m.workspace_repo.List()
			common.CurrWorkspace = workspaces[index]

			term := terminal.NewPreview(&common.CurrWorkspace)
			return m, tea.Batch(term.OpenVim("Update"))

		case "enter":
			index := m.workspace_list.Index()
			workspaces, _ := m.workspace_repo.List()
			common.CurrWorkspace = workspaces[index]

			return m, tea.Batch(common.SetPage(common.Page_Resource))
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

	m.workspace_list.SetItems(m.ItemsOfList())
	m.workspace_list, cmd = m.workspace_list.Update(msg)
	cmds = append(cmds, cmd)

	return m, nil
}

func (m Model) View() string {
	return m.workspace_list.View()
}

type Item struct {
	title, desc string
}

func NewItem(title, desc string) Item {
	return Item{
		title,
		desc,
	}
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

func (m Model) ItemsOfList() []list.Item {
	list := []list.Item{}

	workspaces, _ := m.workspace_repo.List()
	for _, i := range workspaces {
		list = append(list, NewItem(i.Name, i.Uri))
	}

	return list
}
