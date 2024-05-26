package details

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/utils"
)

func (m Model) Preview() string {
	maxPreviewTextSize := m.Width - m.Width/3
	preview := utils.Truncate(fmt.Sprintf("%s - %s", strings.ToUpper(string(m.Request.Method)), string(m.Request.Endpoint)), maxPreviewTextSize)

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
