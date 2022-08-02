package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width  int
	Height int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)
	}

	return m, nil
}

func (m Model) View() string {
	style := lipgloss.
		NewStyle().
		Bold(true).
		Height(1).
		MarginRight(4).
		Border(lipgloss.RoundedBorder(), true, true).
		Foreground(lipgloss.Color("#4ff")).
		BorderForeground(lipgloss.Color("#4ff")).
		Padding(0, 1)

	return style.Render("httui")

}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.Width = msg.Width
	m.Height = msg.Height
}
