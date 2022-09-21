package validate

import (
	"fmt"

	"github.com/gandarfh/maid-san/pkg/convert"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/utils"
)

func InputErrors(args []string, decode any) error {
	mappedArgs, err := utils.ArgsFormat(args[1:])
	if err != nil {
		fmt.Println("\nArgs Format:")
		fmt.Println(err.Error())
		fmt.Println()
		return errors.BadRequest()
	}

	fmt.Printf("%T\n", mappedArgs["parent_id"])

	err = convert.MapToStruct(mappedArgs, decode)
	if err != nil {
		fmt.Println("\nMap to struct:")
		fmt.Println(err.Error())
		fmt.Println()
		return errors.BadRequest()
	}

	validator := NewValidator()
	if err := validator.Struct(decode); err != nil {
		errorList := []string{"Unprocessable Entity!\n"}

		for key, value := range ValidatorErrors(err) {
			er := fmt.Sprintf("[%s] - %s", key, value)
			errorList = append(errorList, er)
		}

		return errors.UnprocessableEntity(errorList...)
	}

	return nil
}
