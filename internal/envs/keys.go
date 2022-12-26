package envs

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Filter       key.Binding
	Delete       key.Binding
	Create       key.Binding
	Left         key.Binding
	FastRename   key.Binding
	CustomRename key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Delete,
		k.Create,
		k.Left,
		k.FastRename,
		k.CustomRename,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = KeyMap{
	Left: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h/shift+tab", "Back"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete"),
	),
	Create: key.NewBinding(
		key.WithKeys("c", "a"),
		key.WithHelp("c/a", "Create"),
	),
	FastRename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Fast Rename"),
	),
	CustomRename: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "Custom rename"),
	),
}

func (m Model) Help() string {
	return m.help.View(keys)
}
