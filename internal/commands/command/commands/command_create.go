package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/command/dtos"
	"github.com/gandarfh/maid-san/internal/commands/command/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
)

type Create struct {
	cmd dtos.InputCreate
}

func (c *Create) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.cmd); err != nil {
		return err
	}

	return nil
}

func (c *Create) Eval() error {
	repo, err := repository.NewCommandRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	cmd := repository.Commands{}

	repo.Create(&cmd)

	return nil
}

func (c *Create) Print() error {
	msg := fmt.Sprintf("[%s] command created with success!\n", c.cmd.Name)

	fmt.Println(msg)

	return nil
}

func (w *Create) Run(args ...string) error {
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

func CreateInit() repl.Repl {
	return &Create{}
}
