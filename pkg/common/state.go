package common

import tea "github.com/charmbracelet/bubbletea"

type State int

const (
	Start_state State = iota
	Error_state
	Loaded_state
)

func SetState(value State) tea.Cmd {
	return func() tea.Msg {
		return value
	}
}
