package errors

type ProcessErrors struct {
	Status  int
	Message string
}

func (r *ProcessErrors) Error() string {
	return ""
}

func NotFoundError(msg ...string) error {
	if len(msg) > 0 {
		return &ProcessErrors{
			Status:  404,
			Message: msg[0],
		}
	}

	return &ProcessErrors{
		Status:  404,
		Message: "Command not found!",
	}
}

func ReadError(msg ...string) error {

	return nil
}

func EvalError(msg ...string) error {

	return nil
}

func PrintError(msg ...string) error {

	return nil
}
