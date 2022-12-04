package tabs

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/maid-san/pkg/styles"
)

var (
	activeFirstTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "│",
		BottomRight: "╰",
	}
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╯",
		BottomRight: "╰",
	}
	firstTabBorder = lipgloss.Border{
		Top:         "",
		Bottom:      "─",
		Left:        "",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "╭",
		BottomRight: "─",
	}
	tabBorder = lipgloss.Border{
		Top:         "",
		Bottom:      "─",
		Left:        "",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "─",
		BottomRight: "─",
	}

	gapBorder = lipgloss.Border{
		Top:         "",
		Bottom:      "─",
		Left:        "",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "─",
		BottomRight: "╮",
	}

	tabStyle = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(styles.DefaultTheme.SecondaryBorder).
			Padding(0, 1)

	firstTabStyle = lipgloss.NewStyle().
			Border(firstTabBorder, true).
			BorderForeground(styles.DefaultTheme.SecondaryBorder).
			Padding(0, 1)

	activeFirstTabStyle = tabStyle.Copy().
				Border(activeFirstTabBorder, true).
				Foreground(styles.DefaultTheme.PrimaryBorder).
				BorderForeground(styles.DefaultTheme.PrimaryBorder)

	activeTabStyle = tabStyle.Copy().
			Border(activeTabBorder, true).
			Foreground(styles.DefaultTheme.PrimaryBorder).
			BorderForeground(styles.DefaultTheme.PrimaryBorder)

	tabGapStyle = tabStyle.Copy().Padding(0, 1).
			Border(gapBorder, true)
)

func Tabs(items Contents, active int, size int) string {
	var listoftabs []string = []string{}

	for i, item := range items {
		listoftabs = append(listoftabs, tab(item.Tab, i == active, i))
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listoftabs...,
	)

	gap := tabGapStyle.Render(strings.Repeat(" ", max(0, (size)-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	return row
}

func tab(text string, active bool, page int) string {
	if active && page == 0 {
		return activeFirstTabStyle.Render(text)
	}

	if page == 0 {
		return firstTabStyle.Render(text)
	}

	if active {
		return activeTabStyle.Render(text)
	}

	return tabStyle.Render(text)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
