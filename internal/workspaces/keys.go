package workspaces

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Filter       key.Binding
	Delete       key.Binding
	Create       key.Binding
	Enter        key.Binding
	FastRename   key.Binding
	CustomRename key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		// k.Filter,
		k.Delete,
		k.Create,
		k.Enter,
		k.FastRename,
		k.CustomRename,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = KeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter", "l"),
		key.WithHelp("enter/l", "Select"),
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
	return m.help.View(m.keys)
}
