package details

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/utils"
)

func (m Model) Params() string {
	paramrenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styleConfig),
		glamour.WithWordWrap(m.Width/3-10),
	)

	rawparams := utils.GetAllParentsParams(m.Request.ParentID, m.Request.QueryParams.Data())
	query, _ := paramrenderer.Render("```json\n" + DataToString(rawparams, 10, m.Workspace.ID) + "\n  ```")

	query_box := lipgloss.NewStyle().Height(m.Height / 2).Render(lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Params:")),
		query,
	))

	return query_box
}
