package process

import (
	"github.com/gandarfh/httui-repl/pkg/errors"
	"github.com/gandarfh/httui-repl/pkg/repl"
)

// [ command ] [ method ] [ values... ]
// create      workspace  name=api-prd baseUrl=localhost:5000
func Start(args []string, commands map[string]repl.Repl) error {
	command := (args[0])

	// Try find method, if dont find return a error to print a not found message
	// Will render `errors` method created into methods folder
	method, ok := commands[command]
	if !ok {
		return errors.NotFoundError()
	}

	err := method.Run(args...)
	if err != nil {
		return errors.NotFoundError()
	}

	return nil
}
