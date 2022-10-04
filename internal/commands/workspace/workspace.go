package workspace

import (
	"github.com/gandarfh/maid-san/internal/commands/workspace/commands"
	"github.com/gandarfh/maid-san/internal/commands/workspace/repository"
	"github.com/gandarfh/maid-san/pkg/repl"
)

type Workspaces struct {
	workspace repository.Workspaces
	wks       *[]repository.Workspaces
}

func (w *Workspaces) Read(args ...string) error {
	return nil
}

func (w *Workspaces) Eval() error {
	return nil
}

func (w *Workspaces) Print() error {
	return nil
}

func (w *Workspaces) Run(args ...string) error {
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
	return &Workspaces{}
}

func SubCommands() repl.CommandList {
	subs := repl.CommandList{
		{Key: "update", Repl: commands.UpdateInit()},
		{Key: "create", Repl: commands.CreateInit()},
		{Key: "list", Repl: commands.ListInit()},
	}

	// sub commands from update command
	// [workspace] [update] [resourceId]
	subs = append(subs, commands.UpdateSubs()...)

	return subs
}
