package commands

import (
	"strconv"

	"github.com/gandarfh/maid-san/internal/commands/workspace/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/table"
)

type List struct {
	wks *[]repository.Workspaces
}

func (c *List) Read(args ...string) error {
	return nil
}

func (c *List) Eval() error {
	repo, err := repository.NewWorkspaceRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	c.wks = repo.List()

	return nil
}

func (c *List) Print() error {
	tbl := table.NewTable([]string{"id", "name", "uri", "created at"})
	rows := []table.Row{}

	for _, item := range *c.wks {
		row := table.Row{
			strconv.FormatUint(uint64(item.ID), 10),
			item.Name,
			item.Uri,
			item.CreatedAt.Format("2006/01/02"),
		}
		rows = append(rows, row)
	}

	tbl.SetRows(rows)

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
