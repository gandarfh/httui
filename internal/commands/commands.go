package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/clear"
	"github.com/gandarfh/maid-san/internal/commands/exit"
	"github.com/gandarfh/maid-san/internal/commands/welcome"
	"github.com/gandarfh/maid-san/internal/commands/workspace"
	"github.com/gandarfh/maid-san/pkg/repl"
)

func Cmds() []repl.Command {
	var (
		subs repl.SubCommands
	)

	commands := []repl.Command{
		{Key: "workspace", Repl: workspace.Init(), SubCommands: workspace.SubCommands()},
		{Key: "exit", Repl: exit.Init()},
		{Key: "clear", Repl: clear.Init()},
		{Key: "welcome", Repl: welcome.Init()},
	}

	subs = workspace.SubCommands()
	commands = appendSubs(commands, subs, "workspace")

	return commands
}

func appendSubs(commands []repl.Command, subs repl.SubCommands, parrent string) []repl.Command {
	for _, key := range subs {
		commands = append(commands, repl.Command{
			Key:    fmt.Sprintf("%s %s", parrent, key.Key),
			Parent: key.Parent,
			Repl:   key.Repl,
		})
	}

	return commands
}
