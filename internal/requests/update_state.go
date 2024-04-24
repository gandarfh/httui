package requests

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/terminal"
)

func (m Model) StateActions(msg common.State) (Model, tea.Cmd) {
	m.state = msg

	switch msg {
	case common.Start_state:
		return m, common.SetState(common.Loaded_state)

	case common.Exit_state:
		return m, tea.Sequence(terminal.ClearScreen())

	case common.Loaded_state:
		return m, tea.Sequence(tea.ClearScreen, tea.EnterAltScreen)
	}

	return m, nil
}
