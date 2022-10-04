package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/workspace/dtos"
	"github.com/gandarfh/maid-san/internal/commands/workspace/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
	"github.com/gandarfh/maid-san/pkg/vim"
)

type Create struct {
	inpt    dtos.InputCreate
	withvim bool
}

func (c *Create) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.inpt); err != nil {
		c.withvim = true
	}

	return nil
}

func (c *Create) Eval() error {
	repo, err := repository.NewWorkspaceRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	if c.withvim {
		preview := vim.NewPreview(c.inpt)
		defer preview.Close()

		if err := preview.Open(); err != nil {
			return err
		}

		if err := preview.Execute(&c.inpt); err != nil {
			return err
		}
	}

	repo.Create(&c.inpt)

	return nil
}

func (c *Create) Print() error {
	msg := fmt.Sprintf("[%s] workspace created with success!\n", c.inpt.Name)

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
