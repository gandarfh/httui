package process

import (
	"fmt"

	"github.com/gandarfh/httui-repl/cmd/commands"
	"github.com/gandarfh/httui-repl/cmd/methods/workspace"
)

var methods = map[string]commands.Command{
	"workspace": workspace.Init(),
}

// [ command ] [ method ] [ values... ]
// create      workspace  name=porcao-prd baseUrl=localhost:5000

func Start(tokens []string) (commands.Command, error) {
	command := (tokens[0])

	// Try find method, if dont find return a error to print a not found message
	// Will render `not_found` method created into methods folder
	method, ok := methods[command]
	if !ok {
		return method, fmt.Errorf("Process not found!")
	}

	return method, nil
}
