package common

import tea "github.com/charmbracelet/bubbletea"

type Loading struct {
	Msg   string
	Value bool
}

func SetLoading(loading bool, msg ...string) tea.Cmd {
	return func() tea.Msg {
		if len(msg) == 0 {
			return Loading{Msg: "", Value: loading}
		}

		return Loading{Msg: msg[0], Value: loading}
	}
}
