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

func SubCommands() repl.CommandList {
	subs := repl.CommandList{
		{Key: "create", Repl: commands.CreateInit()},
		{Key: "list", Repl: commands.ListInit()},
		{Key: "delete", Repl: commands.DeleteInit()},
	}

	// sub commands to update
	// [envs] [update] [envId]
	subs = append(subs, commands.UpdateSubs()...)

	// sub commands to delete
	// [envs] [delete] [envId]
	subs = append(subs, commands.DeleteSubs()...)

	return subs
}
