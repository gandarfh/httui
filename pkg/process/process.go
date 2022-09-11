package process

import (
	"github.com/gandarfh/httui-repl/pkg/errors"
	"github.com/gandarfh/httui-repl/pkg/repl"
)

// [ command ] [ method ] [ values... ]
// create      workspace  name=api-prd baseUrl=localhost:5000
func Start(args []string, commands map[string]repl.Repl) error {
	name := (args[0])

	// Try find command, if dont find return a error to print a not found message
	// Will render `errors` command created into methods folder
	command, ok := commands[name]
	if !ok {
		return errors.NotFoundError()
	}

	err := command.Run(args...)
	if err != nil {
		return errors.NotFoundError()
	}

	return nil
}
