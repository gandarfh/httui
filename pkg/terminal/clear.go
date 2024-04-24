package terminal

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func ClearScreen() tea.Cmd {
	cmd := exec.Command("clear")

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return nil
	})
}
