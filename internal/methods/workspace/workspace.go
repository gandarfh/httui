package workspace

import (
	"fmt"

	"github.com/gandarfh/httui-repl/internal/methods/workspace/repository"
	"github.com/gandarfh/httui-repl/pkg/commands"
	"github.com/gandarfh/httui-repl/pkg/convert"
	"github.com/gandarfh/httui-repl/pkg/errors"
	"github.com/gandarfh/httui-repl/pkg/utils"
	"github.com/gandarfh/httui-repl/pkg/validate"
)

type Workspaces struct {
	workspace repository.Workspaces
	wks       *[]repository.Workspaces
	Repo      *repository.WorkspaceRepo
}

func (w *Workspaces) Read(tokens ...string) error {

	mappedArgs, err := utils.ArgsFormat(tokens[1:])

	if err != nil {
		return fmt.Errorf("jaum, ta dando merda aqui carai")
	}

	err = convert.MapToStruct(mappedArgs, &w.workspace)

	if err != nil {
		return err
	}

	validator := validate.NewValidator()

	if err := validator.Struct(w.workspace); err != nil {

		for key, value := range validate.ValidatorErrors(err) {
			fmt.Println(key, value)
		}

	}

	return nil
}

func (w *Workspaces) Eval() error {
	return nil
}

func (w *Workspaces) Print() error {
	// fmt.Println()

	// for _, item := range *w.wks {
	// 	fmt.Println(item)
	// }

	return nil
}

func (w *Workspaces) Help() error {
	return nil
}

func (w *Workspaces) Run(args ...string) error {
	if err := w.Read(args...); err != nil {
		return errors.ReadError("jaum, ta dando merda aqui carai")
	}

	if err := w.Eval(); err != nil {
		return fmt.Errorf("jaum, ta dando merda aqui carai")
	}

	if err := w.Print(); err != nil {
		return fmt.Errorf("jaum, ta dando merda aqui carai")
	}

	return nil
}

func Init() commands.Command {
	repo, err := repository.NewWorkspaceRepo()

	if err != nil {
		fmt.Println("erro init repo")
	}

	return &Workspaces{
		Repo: repo,
	}
}
