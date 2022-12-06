package tabs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/maid-san/pkg/common"
)

var (
	wrapperStyle     = lipgloss.NewStyle().Align(lipgloss.Left)
	tableBorderStyle = lipgloss.Border{
		Top:         "",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "│",
		TopRight:    "│",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}
)

type Content struct {
	Tab     string
	Content common.Component
}

type Contents []Content

func New(items Contents, active int, w, h int, loading common.Loading) string {
	content := items[active].Content

	table := wrapperStyle.
		Border(tableBorderStyle, true).
		BorderTop(false).
		Width(w).
		Height(h).
		Render(
			lipgloss.JoinVertical(0, loading.Msg, content.View()),
		)

	tabs := Tabs(items, active, w)

	return lipgloss.JoinVertical(0, tabs, table)
}
