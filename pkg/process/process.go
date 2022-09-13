package process

import (
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
)

// [ command ] [ sub-command ] [ values... ]
// create      workspace  name=api-prd baseUrl=localhost:5000
func Start(args []string, commands repl.CommandList) error {
	// Try find command, if dont find return a error to print a not found message
	// Will render `errors` command created into methods folder
	command, error := commands.Find(args[0])
	if error != nil {
		return errors.NotFoundError()
	}

	err := command.Repl.Run(args...)
	if err != nil {
		return err
	}

	return nil
}
