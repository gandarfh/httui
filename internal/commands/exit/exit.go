package exit

import (
	"fmt"
	"os"

	"github.com/gandarfh/httui-repl/pkg/repl"
)

type Exit struct{}

func (w *Exit) Read(args ...string) error {
	return nil
}

func (w *Exit) Eval() error {
	os.Exit(0)

	return nil
}

func (w *Exit) Print() error {
	fmt.Println("By my love <3 >.<")

	return nil
}

func (w *Exit) Run(args ...string) error {
	w.Read(args...)
	w.Print()
	defer w.Eval()

	return nil
}

func Init() repl.Repl {
	return &Exit{}
}
