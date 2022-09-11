package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return validate
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		key := strings.Join(strings.Split(err.Namespace(), ".")[1:], ".")

		fields[key] = fmt.Sprintf("failed on the '%s' tag", err.Tag())
	}

	return fields
}
