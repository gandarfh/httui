package home

import (
	"github/gandarfh/httui/ui/styles"
	"strings"

	// tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(styles.DefaultTheme.SecondaryBorder).
		Padding(0, 2)

	activeTab = tab.Copy().
			Border(activeTabBorder, true).
			Foreground(styles.DefaultTheme.PrimaryBorder).
			BorderForeground(styles.DefaultTheme.PrimaryBorder)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)

func Tabs(width int, active int) string {
	tabs := strings.Builder{}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		Tab("Resources", active == 0),
		Tab("Endpoints", active == 1),
	)

	gap := tabGap.Render(strings.Repeat(" ", max(0, (width)-lipgloss.Width(row)-2)))

	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	tabs.WriteString(row)

	return tabs.String()

}

func Tab(text string, active bool) string {
	if active {
		return activeTab.Render(text)
	}

	return tab.Render(text)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
