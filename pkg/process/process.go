package process

import (
	"github.com/gandarfh/httui-repl/pkg/commands"
	"github.com/gandarfh/httui-repl/pkg/errors"
)

// [ command ] [ method ] [ values... ]
// create      workspace  name=api-prd baseUrl=localhost:5000
func Start(tokens []string, commands map[string]commands.Command) error {
	command := (tokens[0])

	// Try find method, if dont find return a error to print a not found message
	// Will render `errors` method created into methods folder
	method, ok := commands[command]
	if !ok {
		return errors.NotFoundError()
	}

	err := Run(method, tokens)
	if err != nil {
		return errors.NotFoundError()
	}

	return nil
}
