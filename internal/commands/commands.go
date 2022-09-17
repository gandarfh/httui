package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/clear"
	"github.com/gandarfh/maid-san/internal/commands/envs"
	"github.com/gandarfh/maid-san/internal/commands/exit"
	"github.com/gandarfh/maid-san/internal/commands/generator"
	"github.com/gandarfh/maid-san/internal/commands/welcome"
	"github.com/gandarfh/maid-san/internal/commands/workspace"
	"github.com/gandarfh/maid-san/pkg/repl"
)

func Cmds() []repl.Command {
	var (
		subs repl.CommandList
	)

	commands := []repl.Command{
		{Key: "generate", Repl: generator.Init()},
		{Key: "workspace", Repl: workspace.Init()},
		{Key: "envs", Repl: envs.Init()},
		{Key: "exit", Repl: exit.Init()},
		{Key: "clear", Repl: clear.Init()},
		{Key: "welcome", Repl: welcome.Init()},
	}

	subs = generator.SubCommands()
	commands = appendSubs(commands, subs, "generate")

	subs = workspace.SubCommands()
	commands = appendSubs(commands, subs, "workspace")

	subs = envs.SubCommands()
	commands = appendSubs(commands, subs, "envs")

	return commands
}

func appendSubs(commands repl.CommandList, subs repl.CommandList, parrent string) []repl.Command {
	for _, key := range subs {
		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf("%s %s", parrent, key.Key),
			Repl: key.Repl,
		})
	}

	return commands
}
