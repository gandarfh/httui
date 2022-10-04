package resources

import (
	"github.com/gandarfh/maid-san/internal/commands/resources/commands"
	"github.com/gandarfh/maid-san/pkg/repl"
)

type Resources struct{}

func (w *Resources) Read(args ...string) error {
	return nil
}

func (w *Resources) Eval() error {
	return nil
}

func (w *Resources) Print() error {
	return nil
}

func (w *Resources) Run(args ...string) error {
	if err := w.Read(args...); err != nil {
		return err
	}

	if err := w.Eval(); err != nil {
		return err
	}

	if err := w.Print(); err != nil {
		return err
	}

	return nil
}

func Init() repl.Repl {
	return &Resources{}
}

func SubCommands() repl.CommandList {
	subs := repl.CommandList{
		{Key: "create", Repl: commands.CreateInit()},
		{Key: "list", Repl: commands.ListInit()},
		{Key: "delete", Repl: commands.DeleteInit()},
		{Key: "exec", Repl: commands.ExecInit()},
		{Key: "vim exec", Repl: commands.ExecInit()},
	}

	// sub commands from update command
	// [resource] [update] [resourceId]
	subs = append(subs, commands.UpdateSubs()...)

	// sub commands to delete resource
	// [resource] [delete] [resourceId]
	subs = append(subs, commands.DeleteSubs()...)

	// sub commands to execute resource
	// [resource] [exec] [resourceId]
	// [resource] [exec] [vim] [resourceId]
	subs = append(subs, commands.ExecSubs()...)

	return subs
}
