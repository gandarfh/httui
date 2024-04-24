package details

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Body() string {

	bodyrenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styleConfig),
		glamour.WithWordWrap(m.Width-m.Width/3-2),
	)

	rawbody := m.Request.Body.Data()
	body, _ := bodyrenderer.Render("```json\n" + DataToString(rawbody, 20, m.Workspace.ID) + "\n```")

	body_box := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Bodys:")),
		body,
	)

	return body_box
}
