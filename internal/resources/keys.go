package resources

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/gandarfh/httui/pkg/common"
)

type KeyMapTag struct {
	Filter       key.Binding
	Delete       key.Binding
	Create       key.Binding
	FullCreate   key.Binding
	Enter        key.Binding
	Next         key.Binding
	Close        key.Binding
	Exec         key.Binding
	FastRename   key.Binding
	CustomRename key.Binding
	Move         key.Binding
	Back         key.Binding
}

func (k KeyMapTag) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Delete,
		k.Create,
		k.FullCreate,
		k.Enter,
		k.Back,
		k.FastRename,
		k.CustomRename,
	}
}

func (k KeyMapTag) FullHelp() [][]key.Binding {
	return nil
}

var keys_tags = KeyMapTag{
	Next: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l/tab", "Next Tab"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "l"),
		key.WithHelp("enter/l", "Select"),
	),
	Exec: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Exec resource"),
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
	FullCreate: key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "Full tag Create"),
	),
	Move: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "Move to Tag"),
	),
	Back: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "Back to Workspace"),
	),
	Close: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Exit current resource"),
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

type KeyMapResource struct {
	Filter     key.Binding
	Delete     key.Binding
	Next       key.Binding
	Create     key.Binding
	Close      key.Binding
	Exec       key.Binding
	FastRename key.Binding
	Move       key.Binding
}

func (k KeyMapResource) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Filter,
		k.Delete,
		k.Create,
		k.Close,
		k.Next,
		k.Exec,
		k.Move,
		k.FastRename,
	}
}

func (k KeyMapResource) FullHelp() [][]key.Binding {
	return nil
}

var keys_resources = KeyMapResource{
	Next: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l/tab", "Next Tab"),
	),
	Exec: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Exec resource"),
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
	Move: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "Move to Tag"),
	),
	Close: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Exit current resource"),
	),
	FastRename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Update Values"),
	),
}

func (m Model) Help() string {
	if common.CurrTab == common.Tab_Tags {
		return m.help_resource.View(keys_tags)
	}

	return m.help_resource.View(keys_resources)
}
