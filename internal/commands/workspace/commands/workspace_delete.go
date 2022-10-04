package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gandarfh/maid-san/internal/commands/workspace/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
)

type Delete struct {
	WorkspaceId uint
}

func (c *Delete) Read(args ...string) error {
	args = strings.Split(args[0], " ")
	envId, err := strconv.Atoi(args[2])

	if err != nil {
		return errors.UnprocessableEntity("Id provided isn't uint")
	}

	c.WorkspaceId = uint(envId)
	return nil
}

func (c *Delete) Eval() error {
	repo, err := repository.NewWorkspaceRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	repo.Delete(c.WorkspaceId)

	return nil
}

func (c *Delete) Print() error {
	msg := fmt.Sprintf("Workspace success deleted!\n")

	fmt.Println(msg)

	return nil
}

func (w *Delete) Run(args ...string) error {
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

func DeleteSubs() repl.CommandList {
	repo, _ := repository.NewWorkspaceRepo()
	list := repo.List()

	commands := repl.CommandList{}
	commands = append(commands, repl.Command{
		Key:  "delete",
		Repl: DeleteInit(),
	})

	for _, item := range *list {
		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf("%s %d", "delete", item.ID),
			Repl: DeleteInit(),
		})
	}

	return commands
}

func DeleteInit() repl.Repl {
	return &Delete{}
}