package errors

import (
	"fmt"

	"github.com/gandarfh/httui-repl/pkg/commands"
	errs "github.com/gandarfh/httui-repl/pkg/errors"
)

type Error struct {
	Command string
	Err     error
	msgs    *errs.ProcessErrors
}

func (w *Error) Read(tokens ...string) error {
	w.Command = tokens[0]

	return nil
}

func (w *Error) Eval() error {
	if err, ok := w.Err.(*errs.ProcessErrors); ok {
		w.msgs = err
		return nil
	}

	w.msgs = &errs.ProcessErrors{
		Status:  500,
		Message: []string{w.Err.Error()},
	}

	return nil
}

func (w *Error) Print() error {
	var msg string
	msg = fmt.Sprintf("[%d] | Error when execute [%s] command, pls try again.\n", w.msgs.Status, w.Command)
	fmt.Println(msg)

	for _, msg := range w.msgs.Message {
		fmt.Println(msg)
	}

	fmt.Println("For more information type: [help]")

	return nil
}

func (w *Error) Help() error {
	return nil
}

func (w *Error) Run(args ...string) error {
	if err := w.Read(args...); err != nil {
		return err
	}

	if err := w.Eval(); err != nil {
		return err
	}

	w.Print()
	return nil
}

func Init(err error) commands.Command {
	return &Error{
		Err: err,
	}
}
