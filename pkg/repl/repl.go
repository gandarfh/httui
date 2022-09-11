package repl

type Repl interface {
	Read(args ...string) error
	Eval() error
	Print() error
	Run(args ...string) error
}
