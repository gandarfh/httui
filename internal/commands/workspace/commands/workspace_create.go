package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/workspace/dtos"
	"github.com/gandarfh/maid-san/internal/commands/workspace/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
)

type Create struct {
	ws dtos.InputCreate
}

func (c *Create) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.ws); err != nil {
		return err
	}

	return nil
}

func (c *Create) Eval() error {
	repo, err := repository.NewWorkspaceRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	ws := repository.Workspaces{
		Uri:  c.ws.Uri,
		Name: c.ws.Name,
	}

	repo.Create(&ws)

	return nil
}

func (c *Create) Print() error {
	msg := fmt.Sprintf("[%s] workspace created with success!\n", c.ws.Name)

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
