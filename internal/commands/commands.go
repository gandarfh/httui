package commands

import (
	"fmt"

	"github.com/gandarfh/httui-repl/internal/commands/clear"
	"github.com/gandarfh/httui-repl/internal/commands/exit"
	"github.com/gandarfh/httui-repl/internal/commands/welcome"
	"github.com/gandarfh/httui-repl/internal/commands/workspace"
	"github.com/gandarfh/httui-repl/pkg/repl"
)

func Cmds() []repl.Command {
	commands := []repl.Command{
		{Key: "workspace", Repl: workspace.Init(), SubCommands: workspace.SubCommands()},
		{Key: "exit", Repl: exit.Init()},
		{Key: "clear", Repl: clear.Init()},
		{Key: "welcome", Repl: welcome.Init()},
	}

	subs := workspace.SubCommands()
	commands = appendSubs(commands, subs)

	return commands
}

func appendSubs(commands []repl.Command, subs repl.SubCommands) []repl.Command {
	for _, key := range subs {
		commands = append(commands, repl.Command{
			Key:    fmt.Sprintf("%s %s", key.Parent, key.Key),
			Parent: key.Parent,
			Repl:   key.Repl,
		})
	}

	return commands
}
