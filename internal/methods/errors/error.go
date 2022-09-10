package errors

import (
	"fmt"

	"github.com/gandarfh/httui-repl/pkg/commands"
	errs "github.com/gandarfh/httui-repl/pkg/errors"
)

type Error struct {
	Command string
	Err     error
	message *errs.ProcessErrors
}

func (w *Error) Read(tokens ...string) error {
	w.Command = tokens[0]

	return nil
}

func (w *Error) Eval() error {
	if err, ok := w.Err.(*errs.ProcessErrors); ok {
		w.message = err
		// w.message.Command = w.Command
		return nil
	}

	w.message = &errs.ProcessErrors{
		Status:  500,
		Message: w.Err.Error(),
	}

	return nil
}

func (w *Error) Print() error {
	var msg string
	msg = fmt.Sprintf("Error when execute [%s] command, pls try again.\n", w.Command)
	fmt.Println(msg)

	msg = fmt.Sprintf("[%d] | %s", w.message.Status, w.message.Message)
	fmt.Println(msg)

	fmt.Println("For more information type: [help]")

	return nil
}

func Init(err error) commands.Command {
	return &Error{
		Err: err,
	}
}
