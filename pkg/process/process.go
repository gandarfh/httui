package process

import (
	"fmt"
	"strings"

	"github.com/gandarfh/httui-repl/pkg/errors"
	"github.com/gandarfh/httui-repl/pkg/repl"
	"github.com/gandarfh/httui-repl/pkg/utils"
)

// [ command ] [ sub-command ] [ values... ]
// create      workspace  name=api-prd baseUrl=localhost:5000
func Start(args []string, commands repl.CommandList) error {
	key := ""

	for _, item := range args {
		if !utils.IsKeyValue(item) {
			key = strings.Trim(fmt.Sprintf("%s %s", key, item), " ")
		}
	}

	// Try find command, if dont find return a error to print a not found message
	// Will render `errors` command created into methods folder
	command, error := commands.Find(key)
	if error != nil {
		return errors.NotFoundError()
	}

	err := command.Repl.Run(args...)
	if err != nil {
		return err
	}

	return nil
}
