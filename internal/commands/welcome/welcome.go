package welcome

import (
	"fmt"

	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/logrusorgru/aurora/v3"
)

const (
	version = "v0.1.0"
)

type Wellcome struct{}

func (w *Wellcome) Read(args ...string) error {
	return nil
}

func (w *Wellcome) Eval() error {

	return nil
}

func (w *Wellcome) Print() error {
	fmt.Println(aurora.Yellow(art))
	fmt.Print("\n")
	fmt.Printf("Okaerinasaimase, Goshujin-sama >.<  ───  %s.\n", aurora.Green(version))
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

var art = `
 _____     _   _
|     |___|_|_| |    ___ ___ ___
| | | | .'| | . | _ |_ -| .'|   |
|_|_|_|__,|_|___|   |___|__,|_|_|`
