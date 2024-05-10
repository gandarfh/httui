package details

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/pkg/common"
)

var (
	styleConfig   = glamour.DarkStyleConfig
	title_style   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	preview_style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
)

type Model struct {
	Width     int
	Height    int
	bodyVP    viewport.Model
	headerVP  viewport.Model
	paramsVP  viewport.Model
	Request   offline.Request
	Workspace offline.Workspace
}

func New() Model {
	return Model{}
}

func (m *Model) SetWorkspace(w offline.Workspace) tea.Cmd {
	return func() tea.Msg {
		return w
	}
}

func (m *Model) SetRequest(r offline.Request) tea.Cmd {
	return func() tea.Msg {
		return r
	}
}

type Width int
type Height int

func (m *Model) SetWidth(w int) tea.Cmd {
	return func() tea.Msg {
		return Width(w)
	}
}

func (m *Model) SetHeight(h int) tea.Cmd {
	return func() tea.Msg {
		return Height(h)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case common.State:
		if msg == common.Start_state {
			m.bodyVP = viewport.New(m.Width-m.Width/3-4, m.Height)
			m.headerVP = viewport.New(m.Width/3-2, m.Height/2-2)
			m.paramsVP = viewport.New(m.Width/3-2, m.Height/2-2)
		}

	case offline.Request:
		m.Request = msg

	case offline.Workspace:
		m.Workspace = msg
	}

	m.bodyVP, cmd = m.bodyVP.Update(msg)
	cmds = append(cmds, cmd)

	m.headerVP, cmd = m.headerVP.Update(msg)
	cmds = append(cmds, cmd)

	m.paramsVP, cmd = m.paramsVP.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.bodyVP.SetContent(m.Body())
	m.headerVP.SetContent(m.Headers())
	m.paramsVP.SetContent(m.Params())

	paramsContent := lipgloss.
		NewStyle().
		Height(m.Height).
		Width(m.Width / 3).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			m.paramsVP.View(),
			lipgloss.NewStyle().Width(m.Width/3).Border(lipgloss.NormalBorder(), false).BorderBottom(true).String(),
			m.headerVP.View(),
		))

	content_row := lipgloss.NewStyle().Width(m.Width).Height(m.Height).Border(lipgloss.RoundedBorder()).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.bodyVP.View(),
		lipgloss.NewStyle().Height(m.Height).Width(3).Border(lipgloss.NormalBorder(), false).BorderRight(true).String(),
		paramsContent,
	))

	container := lipgloss.NewStyle().Padding(0, 1).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.Title(),
		m.Preview(),
		content_row,
	))

	return container
}
