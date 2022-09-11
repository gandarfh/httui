package errors

type ProcessErrors struct {
	Status  int
	Message []string
}

func (r *ProcessErrors) Error() string {
	return ""
}

func NotFoundError(msg ...string) error {
	if len(msg) > 0 {
		return &ProcessErrors{
			Status:  404,
			Message: msg,
		}
	}

	return &ProcessErrors{
		Status:  404,
		Message: []string{"Command not found!"},
	}
}

func ReadError(msg ...string) error {

	return &ProcessErrors{
		Status:  404,
		Message: msg,
	}

}

func EvalError(msg ...string) error {

	return &ProcessErrors{
		Status:  404,
		Message: msg,
	}
}

func PrintError(msg ...string) error {

	return &ProcessErrors{
		Status:  404,
		Message: msg,
	}
}
