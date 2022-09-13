package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/envs/dtos"
	"github.com/gandarfh/maid-san/internal/commands/envs/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
)

type Create struct {
	env dtos.InputCreate
}

func (c *Create) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.env); err != nil {
		return err
	}

	return nil
}

func (c *Create) Eval() error {
	repo, err := repository.NewEnvsRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	env := repository.Envs{
		Key:   c.env.Key,
		Value: c.env.Value,
	}

	repo.Create(&env)

	return nil
}

func (c *Create) Print() error {
	msg := fmt.Sprintf("[%s] Env success created!\n", c.env.Key)

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
