package tabs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/maid-san/pkg/styles"
	"golang.org/x/term"
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
	Content string
}

type Contents []Content

func New(items Contents, active int, loading string) string {
	content := items[active].Content
	w, h, _ := term.GetSize(0)

	teste := w - styles.Container.Base.GetHorizontalPadding() - 2
	teste2 := h - styles.Container.Base.GetVerticalPadding() - 8

	table := wrapperStyle.
		Border(tableBorderStyle, true).
		BorderTop(false).
		Width(teste).
		Height(teste2).
		Render(
			lipgloss.JoinVertical(0, loading, content),
		)

	tabs := Tabs(items, active, teste)

	return lipgloss.JoinVertical(0, tabs, table)
}
