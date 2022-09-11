package clear

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/gandarfh/httui-repl/pkg/repl"
)

type Clear struct{}

func (w *Clear) Read(args ...string) error {
	return nil
}

func (w *Clear) Eval() error {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}

	return nil
}

func (w *Clear) Print() error {
	return nil
}

func (w *Clear) Run(args ...string) error {
	w.Read(args...)
	w.Print()
	defer w.Eval()

	return nil
}

func Init() repl.Repl {
	return &Clear{}
}

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func clearTerminal() {
}
