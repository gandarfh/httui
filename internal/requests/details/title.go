package details

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) Title() string {
	title := title_style.Width(m.Width).Border(lipgloss.RoundedBorder()).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Width(m.Width-m.Width/3).Bold(true).Render(" Name: "+m.Request.Name),
		lipgloss.NewStyle().Width(m.Width/6).String(),
		lipgloss.NewStyle().Width(m.Width/6).Bold(true).Render(fmt.Sprint(" ID: ", m.Request.ID)),
	))

	return title
}
