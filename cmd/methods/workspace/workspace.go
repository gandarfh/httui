package workspace

import (
	"fmt"

	"github.com/gandarfh/httui-repl/cmd/commands"
)

type Workspaces struct {
	Id      string
	Name    string
	BaseUrl string
}

func (w *Workspaces) Read(tokens []string) error {
	w.Name = tokens[0]
	return nil
}

func (w *Workspaces) Eval() error {
	return nil
}

func (w *Workspaces) Print() error {

	fmt.Println(w.Name)
	return nil
}

func Init() commands.Command {
	return &Workspaces{}
}
