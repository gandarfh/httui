package requests

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Filter     key.Binding
	Workspace  key.Binding
	OpenGroup  key.Binding
	CloseGroup key.Binding
	Delete     key.Binding
	Create     key.Binding
	Edit       key.Binding
	Exec       key.Binding
	Envs       key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Filter,
		k.Workspace,
		k.OpenGroup,
		k.CloseGroup,
		k.Delete,
		k.Create,
		k.Edit,
		k.Exec,
		k.Envs,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = KeyMap{
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
	Workspace: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "Select Workspaces"),
	),
	OpenGroup: key.NewBinding(
		key.WithKeys("o", "enter", "l"),
		key.WithHelp("o/l/enter", "Open"),
	),
	CloseGroup: key.NewBinding(
		key.WithKeys("O", "h"),
		key.WithHelp("O/h", "Close"),
	),
	Edit: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "Edit Request"),
	),
	Exec: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Exec Request"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete"),
	),
	Create: key.NewBinding(
		key.WithKeys("c", "a"),
		key.WithHelp("c/a", "Create"),
	),
	Envs: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "Open Envs"),
	),
}

func (m Model) Help() string {
	return m.help.View(m.keys)
}
