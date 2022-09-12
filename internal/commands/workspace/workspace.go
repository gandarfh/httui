package workspace

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/welcome"
	"github.com/gandarfh/maid-san/internal/commands/workspace/repository"
	"github.com/gandarfh/maid-san/pkg/convert"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/utils"
	"github.com/gandarfh/maid-san/pkg/validate"
)

type Workspaces struct {
	workspace repository.Workspaces
	wks       *[]repository.Workspaces
	Repo      *repository.WorkspaceRepo
}

func SubCommands() repl.SubCommands {
	return repl.SubCommands{
		{Key: "run", Parent: "workspace", Repl: welcome.Init()},
	}

}

func (w *Workspaces) Read(args ...string) error {
	mappedArgs, err := utils.ArgsFormat(args[1:])
	if err != nil {
		return errors.BadRequest()
	}

	err = convert.MapToStruct(mappedArgs, &w.workspace)
	if err != nil {
		return errors.BadRequest()
	}

	validator := validate.NewValidator()

	if err := validator.Struct(w.workspace); err != nil {
		errorList := []string{"Unprocessable Entity!\n"}

		// TODO: throw a error
		for key, value := range validate.ValidatorErrors(err) {
			er := fmt.Sprintf("[%s] - %s", key, value)
			errorList = append(errorList, er)
		}

		return errors.UnprocessableEntity(errorList...)

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

func (w *Workspaces) Run(args ...string) error {
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

func Init() repl.Repl {
	repo, err := repository.NewWorkspaceRepo()
	if err != nil {
		fmt.Println("erro init repo")
	}

	return &Workspaces{
		Repo: repo,
	}
}
