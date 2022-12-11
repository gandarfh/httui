package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/pkg/common"
)

type Model struct {
	width     int
	height    int
	textInput textinput.Model
}

func New() Model {
	ti := textinput.New()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Focus()

	return Model{textInput: ti}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case common.Command:
		switch msg.Category {
		case "FILTER":
			m.textInput.Prompt = "/"
		default:
			m.textInput.Prompt = "> "
		}

		m.textInput.SetValue(msg.Value)

	case tea.KeyMsg:
		switch msg.Type.String() {
		case "enter", "ctrl+c", "esc":
			value := m.textInput.Value()
			m.textInput.Reset()

			return m, tea.Batch(
				common.CloseCommand(),
				common.SetCommand(value),
			)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.textInput.View()
}
