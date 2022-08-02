package home

import (
	"fmt"
	"github/gandarfh/httui/config"
	"github/gandarfh/httui/data"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	resources []data.Resources
	endpoints []data.Endpoints
	selected  string
	Width     int
	Height    int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (Model, tea.Cmd) {
	connect := config.Connect[data.Data]()

	m = Model{
		resources: connect.Resources,
		endpoints: connect.GetAllEndpoints(),
		selected:  "resources",
	}

	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			fmt.Println("jaum")

			return m, nil

		}

	}

	return m, nil
}

func (m Model) View() string {
	style := lipgloss.NewStyle()

	return style.Render("Home")

}

func (m Model) Help() string {

	return "help home"
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.Width = msg.Width
	m.Height = msg.Height
}
