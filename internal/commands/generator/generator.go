package generator

import (
	"github.com/gandarfh/maid-san/internal/commands/generator/commands"
	"github.com/gandarfh/maid-san/pkg/repl"
)

type Commands struct{}

func (w *Commands) Read(args ...string) error {
	return nil
}

func (w *Commands) Eval() error {
	return nil
}

func (w *Commands) Print() error {
	return nil
}

func (w *Commands) Run(args ...string) error {
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
	return &Commands{}
}

func SubCommands() repl.CommandList {
	return repl.CommandList{
		{Key: "create", Repl: commands.CreateInit()},
	}
}
