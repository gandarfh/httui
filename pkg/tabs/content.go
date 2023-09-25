package tabs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/common"
)

var (
	wrapperStyle     = lipgloss.NewStyle().Align(lipgloss.Left)
	tableBorderStyle = lipgloss.RoundedBorder()
)

type Content struct {
	Tab     string
	Content common.Component
}

type Contents []Content

func New(items Contents, active int, w, h int, loading_component common.Loading, env_component common.Environment) string {
	content := items[active].Content

	table := wrapperStyle.
		Border(tableBorderStyle, true).
		BorderTop(false).
		Width(w).
		Height(h).
		Render(content.View())

	tabs := Tabs(items, active, w, loading_component.Msg, env_component.Name)

	return lipgloss.JoinVertical(0, tabs, table)
}
