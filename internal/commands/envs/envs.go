package envs

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/envs/commands"
	"github.com/gandarfh/maid-san/pkg/repl"
)

type Envs struct{}

func (w *Envs) Read(args ...string) error {
	fmt.Println(args)
	return nil
}

func (w *Envs) Eval() error {
	return nil
}

func (w *Envs) Print() error {
	return nil
}

func (w *Envs) Run(args ...string) error {
	w.Read(args...)
	w.Eval()
	w.Print()

	return nil
}

func Init() repl.Repl {
	return &Envs{}
}

func SubCommands() repl.SubCommands {
	return repl.SubCommands{
		{Key: "create", Repl: commands.CreateInit()},
		{Key: "list", Repl: commands.ListInit()},
	}
}
