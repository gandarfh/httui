package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/maid-san/internal/envs"
	"github.com/gandarfh/maid-san/internal/resources"
	"github.com/gandarfh/maid-san/internal/workspaces"
	"github.com/gandarfh/maid-san/pkg/common"
	"github.com/gandarfh/maid-san/pkg/styles"
	"github.com/gandarfh/maid-san/pkg/tabs"
)

type state int

const (
	start_state state = iota
	error_state
	loaded_state
)

type Model struct {
	pages  tabs.Contents
	width  int
	height int
	state  state
}

func New() Model {
	pages := tabs.Contents{
		{Tab: "Workspaces"},
		{Tab: "Resouces"},
		{Tab: "Environments"},
	}

	return Model{pages: pages, state: start_state}
}

func (m Model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	m.pages[common.Page_Workspace].Content = workspaces.New()
	m.pages[common.Page_Resource].Content = resources.New()
	m.pages[common.Page_Env].Content = envs.New()

	for _, p := range m.pages {
		cmds = append(cmds, p.Content.Init())
	}

	m.state = loaded_state

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.state = loaded_state

		for p, i := range m.pages {
			m.pages[p].Content, cmd = i.Content.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.Page:
		common.CurrPage = msg
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l":
			common.CurrPage = min(common.CurrPage+1, len(m.pages)-1)
			return m, nil
		case "left", "h":
			common.CurrPage = max(common.CurrPage-1, 0)
			return m, nil
		}
	}

	if m.state == loaded_state {
		m.pages[common.CurrPage].Content, cmd = m.pages[common.CurrPage].Content.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	content := ""

	w := m.width - (m.width / 12)
	h := m.height - (m.height / 3)

	switch m.state {
	case loaded_state:
		content = tabs.New(m.pages, common.CurrPage, w, h)
	default:
		content = tabs.New(m.pages, common.CurrPage, w, h)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, styles.Container.Base.Render(content))
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
