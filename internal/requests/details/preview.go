package details

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/utils"
)

func (m Model) Preview() string {

	Preview := utils.ReplaceByOperator(
		utils.Truncate(fmt.Sprintf("%s - %s", m.Request.Method, m.Request.Endpoint), 100),
		m.Workspace.ID,
	)

	text := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(
		lipgloss.Place(
			m.Width,
			1,
			lipgloss.Left,
			lipgloss.Top,
			lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Preview: ", Preview)),
		))

	return text
}
