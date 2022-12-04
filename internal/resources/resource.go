package resources

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/maid-san/pkg/common"
)

type Model struct {
	width int
	hight int
}

func New() Model {

	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.hight = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// case "ctrl+c", "q":
		// return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {

	return fmt.Sprint("resources ", common.CurrWorkspace.Name)
}
