package login

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/services"
)

type Model struct {
	Width  int
	Height int
	keys   KeyMap
	url    string
}

func New() Model {
	return Model{
		keys: keys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.url = services.AuthCLI("Jaum")

	return fmt.Sprintf("Press Enter to open the browser or visit %s", m.url)
}
