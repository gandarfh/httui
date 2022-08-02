package requests

import (
	"fmt"
	"github/gandarfh/httui/config"
	"github/gandarfh/httui/data"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (Model, tea.Cmd) {
	config.Connect[data.Data]() // connect

	m = Model{}

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			fmt.Println("requests")
			return m, nil

		}

	}

	return m, nil
}

func (m Model) View() string {
	style := lipgloss.NewStyle()

	return style.Render("Requests")

}

func (m Model) Help() string {

	return "help requests"
}
