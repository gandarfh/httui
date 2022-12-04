package styles

import "github.com/charmbracelet/lipgloss"

type container struct {
	Base     lipgloss.Style
	Loading  lipgloss.Style
	Resource lipgloss.Style
}

var Container = func() container {
	base := lipgloss.NewStyle().MarginTop(2).Padding(0, 3)
	loading := lipgloss.NewStyle().Padding(0, 3).PaddingTop(1)
	resource := lipgloss.NewStyle().Padding(1)

	return container{
		Base:     base,
		Loading:  loading,
		Resource: resource,
	}
}()
