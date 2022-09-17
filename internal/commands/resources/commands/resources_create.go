package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/resources/dtos"
	"github.com/gandarfh/maid-san/internal/commands/resources/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
)

type Create struct {
	inpt dtos.InputCreate
}

func (c *Create) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.inpt); err != nil {
		return err
	}

	return nil
}

func (c *Create) Eval() error {
	repo, err := repository.NewResourcesRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	value := repository.Resources{}

	repo.Create(&value)

	return nil
}

func (c *Create) Print() error {
	msg := fmt.Sprintf("[%s] Resource success created!\n", c.inpt.Name)

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
