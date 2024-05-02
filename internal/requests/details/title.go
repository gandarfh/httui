package details

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) Title() string {
  titleWidth := m.Width-m.Width/3

	title := title_style.Width(m.Width).Border(lipgloss.RoundedBorder()).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Width(titleWidth).Bold(true).Render(" Name: "+m.Request.Name),
		lipgloss.NewStyle().Width(titleWidth/4).String(),
		lipgloss.NewStyle().Width(titleWidth/4).Bold(true).Render(fmt.Sprint(" ID: ", m.Request.ID)),
	))

	return title
}
