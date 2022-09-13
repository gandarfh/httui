package validate

import (
	"fmt"

	"github.com/gandarfh/maid-san/pkg/convert"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/utils"
)

func ValidateInputErrors(args []string, decode any) error {
	mappedArgs, err := utils.ArgsFormat(args[1:])
	if err != nil {
		return errors.BadRequest()
	}

	err = convert.MapToStruct(mappedArgs, decode)
	if err != nil {
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
