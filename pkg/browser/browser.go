package browser

import (
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

func OpenPage(url string) tea.Cmd {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	}

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return nil
	})
}
