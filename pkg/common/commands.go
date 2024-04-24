package common

import tea "github.com/charmbracelet/bubbletea"

type Command struct {
	Active   bool
	Value    string
	Category string
	Prefix   string
}

type CommandClose struct {
	Command
}

var command = Command{}

func SetCommand(value string) tea.Cmd {
	return func() tea.Msg {
		command.Value = value
		return command
	}
}

func OpenCommand(category, prefix string) tea.Cmd {
	return func() tea.Msg {
		command.Category = category
		command.Prefix = prefix
		command.Value = ""
		command.Active = true
		return command
	}
}

func CloseCommand() tea.Cmd {
	return func() tea.Msg {
		command.Active = false
		return CommandClose{command}
	}
}

func ClearCommand() tea.Cmd {
	return func() tea.Msg {
		command.Active = false
		command.Prefix = ""
		command.Value = ""
		command.Category = ""
		return command
	}
}
