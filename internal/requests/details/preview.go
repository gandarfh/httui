package details

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/utils"
)

func (m Model) Preview() string {
	maxPreviewTextSize := m.Width - m.Width/3
	preview := utils.Truncate(fmt.Sprintf("%s - %s", m.Request.Method, m.Request.Endpoint), maxPreviewTextSize)

  preview = ""

	text := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(
		lipgloss.Place(
			m.Width,
			1,
			lipgloss.Left,
			lipgloss.Top,
			lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Preview: ", preview)),
		))

	return text
}
