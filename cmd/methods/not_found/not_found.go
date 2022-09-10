package notfound

import (
	"fmt"

	"github.com/gandarfh/httui-repl/cmd/commands"
)

type NotFound struct {
	Name string
}

func (w *NotFound) Read(tokens []string) error {
	w.Name = tokens[0]

	return nil
}

func (w *NotFound) Eval() error {
	return nil
}

func (w *NotFound) Print() error {
	msg := fmt.Sprintf("Command [%s] not exist, pls select another.\n\nFor more information type: [help]", w.Name)

	fmt.Println(msg)
	fmt.Println()

	return nil
}

func Init() commands.Command {
	return &NotFound{}
}
