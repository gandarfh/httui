package requests

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/command"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/styles"
)

type Model struct {
	title            string
	filter           string
	parentId         *uint
	previousParentId *uint
	command_active   bool
	keys             KeyMap
	detail           ModelDetail
	help             help.Model
	List             list.Model
	spinner          spinner.Model
	command_bar      common.Component
	loading          common.Loading
	state            common.State
	Width            int
	Height           int
}

var (
	divider = lipgloss.NewStyle().MarginLeft(1).Border(lipgloss.NormalBorder(), false, true, false, false)
)

func New() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().MarginLeft(2).Foreground(styles.DefaultTheme.PrimaryText)



	m := Model{
		Width:          100,
		state:          common.Start_state,
		List:           NewRequestList(),
		detail:         NewDetail(),
		help:           help.New(),
		keys:           keys,
		spinner:        s,
		command_bar:    command.New(),
		command_active: false,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	cmds = append(cmds, m.command_bar.Init())

	cmd = common.SetState(common.Start_state)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	requestsPage := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.List.View(),
		divider.Height(m.Height-1).String(),
		m.detail.View(),
	)

	content := ""

	if m.command_active {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			requestsPage,
			m.command_bar.View(),
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Right,
			requestsPage,
			lipgloss.NewStyle().Render(m.Help()),
		)
	}

	return styles.Container.Base.Width(m.Width).MaxWidth(m.Width).Render(content)

}
