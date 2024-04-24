package details

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/utils"
)

func (m Model) Headers() string {
	paramrenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styleConfig),
		glamour.WithWordWrap(m.Width/3-10),
	)

	rawheader := utils.GetAllParentsHeaders(m.Request.ParentID, m.Request.Headers.Data())
	header, _ := paramrenderer.Render("```json\n" + DataToString(rawheader, 10, m.Workspace.ID) + "\n  ```")

	header_box := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Headers:")),
		header,
	)

	return header_box
}
