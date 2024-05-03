package login

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return nil
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
	),
}
