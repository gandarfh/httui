package welcome

import (
	"fmt"

	"github.com/gandarfh/httui-repl/pkg/repl"
)

const (
	version = "0.1.0"
)

type Wellcome struct{}

func (w *Wellcome) Read(args ...string) error {
	return nil
}

func (w *Wellcome) Eval() error {

	return nil
}

func (w *Wellcome) Print() error {
	fmt.Printf("Welcome my love >.<  ───  v%s.\n", version)
	fmt.Println("Want more?")
	fmt.Print("\n")

	return nil
}

func (w *Wellcome) Run(args ...string) error {
	w.Read(args...)
	w.Print()
	defer w.Eval()

	return nil
}

func Init() repl.Repl {
	return &Wellcome{}
}
