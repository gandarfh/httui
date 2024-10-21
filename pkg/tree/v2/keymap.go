package tree

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap holds the key bindings for the table.
type KeyMap struct {
	PrevPage     key.Binding
	NextPage     key.Binding
	Bottom       key.Binding
	Top          key.Binding
	Down         key.Binding
	Up           key.Binding
	ToggleExpand key.Binding
}

// DefaultKeyMap is the default key bindings for the table.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		PrevPage: key.NewBinding(),
		NextPage: key.NewBinding(),
		ToggleExpand: key.NewBinding(
			key.WithKeys("enter", "o"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("bottom", "G"),
		),
		Top: key.NewBinding(
			key.WithKeys("top", "g"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
		),
	}
}
