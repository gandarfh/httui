package process

import (
	"github.com/gandarfh/httui-repl/pkg/commands"
	"github.com/gandarfh/httui-repl/pkg/errors"
)

func Run(process commands.Command, tokens []string) error {
	var (
		err error
	)

	err = process.Read(tokens...)
	if err != nil {
		return errors.ReadError()
	}

	err = process.Eval()
	if err != nil {
		return errors.EvalError()
	}

	err = process.Print()
	if err != nil {
		return errors.PrintError()
	}

	return nil
}
