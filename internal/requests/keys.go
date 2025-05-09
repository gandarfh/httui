package requests

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Filter          key.Binding
	SetWorkspace    key.Binding
	CreateWorkspace key.Binding
	OpenGroup       key.Binding
	Delete          key.Binding
	Create          key.Binding
	Edit            key.Binding
	Exec            key.Binding
	Envs            key.Binding
	Detail          key.Binding
	Next            key.Binding
	Save            key.Binding
	Quit            key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Filter,
		k.SetWorkspace,
		k.CreateWorkspace,
		k.OpenGroup,
		k.Delete,
		k.Create,
		k.Edit,
		k.Exec,
		k.Envs,
		k.Next,
		k.Save,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
	),
	Detail: key.NewBinding(
		key.WithKeys("j", "k"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
	SetWorkspace: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "Select Workspaces"),
	),
	CreateWorkspace: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("shift+s", "Create Workspaces"),
	),
	OpenGroup: key.NewBinding(
		key.WithKeys("o", "esc", "enter", "l"),
		key.WithHelp("o/l/enter", "Open"),
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
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab/shift+tab", "Next or Prev"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
	),
}

func (k *KeyMap) DisableKeysForInputs(v bool) {
	k.Detail.SetEnabled(v)
	k.Filter.SetEnabled(v)
	k.SetWorkspace.SetEnabled(v)
	k.CreateWorkspace.SetEnabled(v)
	k.OpenGroup.SetEnabled(v)
	k.Delete.SetEnabled(v)
	k.Create.SetEnabled(v)
	k.Edit.SetEnabled(v)
	k.Exec.SetEnabled(v)
	k.Envs.SetEnabled(v)
	k.Quit.SetEnabled(v)
}

func (m Model) Help() string {
	return m.help.View(m.keys)
}
