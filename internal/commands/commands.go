package commands

import (
	"github.com/gandarfh/httui-repl/internal/commands/clear"
	"github.com/gandarfh/httui-repl/internal/commands/exit"
	"github.com/gandarfh/httui-repl/internal/commands/workspace"
	"github.com/gandarfh/httui-repl/pkg/repl"
)

var Commands = map[string]repl.Repl{
	"workspace": workspace.Init(),
	"exit":      exit.Init(),
	"clear":     clear.Init(),
}
