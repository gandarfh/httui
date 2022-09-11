package workspace

import (
	"fmt"

	"github.com/gandarfh/httui-repl/internal/commands/welcome"
	"github.com/gandarfh/httui-repl/internal/commands/workspace/repository"
	"github.com/gandarfh/httui-repl/pkg/convert"
	"github.com/gandarfh/httui-repl/pkg/errors"
	"github.com/gandarfh/httui-repl/pkg/repl"
	"github.com/gandarfh/httui-repl/pkg/utils"
	"github.com/gandarfh/httui-repl/pkg/validate"
)

type Workspaces struct {
	workspace repository.Workspaces
	wks       *[]repository.Workspaces
	Repo      *repository.WorkspaceRepo
}

func SubCommands() repl.SubCommands {
	subs := repl.SubCommands{
		{Parent: "workspace", Key: "help", Repl: welcome.Init()},
	}

	return subs
}

func (w *Workspaces) Read(args ...string) error {
	mappedArgs, err := utils.ArgsFormat(args[1:])

	if err != nil {
		return fmt.Errorf("jaum, ta dando merda aqui carai")
	}

	err = convert.MapToStruct(mappedArgs, &w.workspace)

	if err != nil {
		return err
	}

	validator := validate.NewValidator()

	if err := validator.Struct(w.workspace); err != nil {

		// TODO: throw a error
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

func Init() repl.Repl {
	repo, err := repository.NewWorkspaceRepo()
	if err != nil {
		fmt.Println("erro init repo")
	}

	return &Workspaces{
		Repo: repo,
	}
}
