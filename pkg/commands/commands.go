package commands

type Command interface {
	Read(args ...string) error
	Eval() error
	Print() error
}
