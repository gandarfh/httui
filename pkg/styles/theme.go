package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	SingleRuneWidth    = 4
	MainContentPadding = 1
)

type Theme struct {
	SelectedBackground lipgloss.AdaptiveColor
	PrimaryBorder      lipgloss.AdaptiveColor
	FaintBorder        lipgloss.AdaptiveColor
	SecondaryBorder    lipgloss.AdaptiveColor
	FaintText          lipgloss.AdaptiveColor
	PrimaryText        lipgloss.AdaptiveColor
	SecondaryText      lipgloss.AdaptiveColor
	InvertedText       lipgloss.AdaptiveColor
	SuccessText        lipgloss.AdaptiveColor
	WarningText        lipgloss.AdaptiveColor
}

var theme *Theme

var DefaultTheme = func() Theme {
	theme = &Theme{
		SelectedBackground: lipgloss.AdaptiveColor{Light: "000", Dark: "000"},
		PrimaryBorder:      lipgloss.AdaptiveColor{Light: "003", Dark: "003"},
		SecondaryBorder:    lipgloss.AdaptiveColor{Light: "244", Dark: "244"},
		FaintBorder:        lipgloss.AdaptiveColor{Light: "254", Dark: "000"},
		PrimaryText:        lipgloss.AdaptiveColor{Light: "003", Dark: "003"},
		SecondaryText:      lipgloss.AdaptiveColor{Light: "094", Dark: "094"},
		FaintText:          lipgloss.AdaptiveColor{Light: "007", Dark: "249"},
		InvertedText:       lipgloss.AdaptiveColor{Light: "015", Dark: "236"},
		SuccessText:        lipgloss.AdaptiveColor{Light: "002", Dark: "002"},
		WarningText:        lipgloss.AdaptiveColor{Light: "001", Dark: "001"},
	}

	return *theme
}()
