package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/envs/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
)

type List struct {
	envs *[]repository.Envs
}

func (c *List) Read(args ...string) error {
	return nil
}

func (c *List) Eval() error {
	repo, err := repository.NewEnvsRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	c.envs = repo.List()

	return nil
}

func (c *List) Print() error {
	for i, env := range *c.envs {
		msg := fmt.Sprintf("[%d] key: %s, value: %s", i, env.Key, env.Value)

		fmt.Println(msg)
	}

	return nil
}

func (w *List) Run(args ...string) error {
	if err := w.Read(args...); err != nil {
		return err
	}

	if err := w.Eval(); err != nil {
		return err
	}

	if err := w.Print(); err != nil {
		return err
	}

	return nil
}

func ListInit() repl.Repl {
	return &List{}
}
