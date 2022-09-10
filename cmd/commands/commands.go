package commands

type Command interface {
	Read(tokens []string) error
	Eval() error
	Print() error
}
