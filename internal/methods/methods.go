package methods

import (
	"github.com/gandarfh/httui-repl/internal/methods/exit"
	"github.com/gandarfh/httui-repl/internal/methods/workspace"
	"github.com/gandarfh/httui-repl/pkg/commands"
)

var Commands = map[string]commands.Command{
	"workspace": workspace.Init(),
	"exit":      exit.Init(),
}
