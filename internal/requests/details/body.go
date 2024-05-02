package details

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Body() string {
	bodyWidth := m.Width / 3
	maxTextValueSize := bodyWidth / 3

	renderTerm, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styleConfig),
		glamour.WithWordWrap(bodyWidth),
	)

	rawbody := m.Request.Body.Data()
	jsonString := DataToString(rawbody, maxTextValueSize, m.Height-m.Height/4)
	bodyJson, _ := renderTerm.Render("```json\n" + jsonString + "\n  ```")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Bodys:")),
		bodyJson,
	)

	return content
}
