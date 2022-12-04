package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/maid-san/pkg/styles"
	"github.com/gandarfh/maid-san/pkg/tabs"
)

type active uint

const (
	page_workspace = iota
	page_resource
	page_env
)

type Model struct {
	Tabs  tabs.Contents
	width int
	hight int
	page  int
}

func New() Model {
	items := tabs.Contents{
		{Tab: "Workspaces", Content: "Workspaces"},
		{Tab: "Resouces", Content: "Workspaces"},
		{Tab: "Envs", Content: "Workspaces"},
	}

	return Model{Tabs: items, page: 0}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.hight = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.page = min(m.page+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.page = max(m.page-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	dash := tabs.New(m.Tabs, m.page, "carregando")
	return styles.Container.Base.Render(dash)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
