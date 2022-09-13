package errors

type ProcessErrors struct {
	Status  int
	Message []string
}

func (r *ProcessErrors) Error() string {
	return ""
}

func BadRequest(msg ...string) error {
	if len(msg) > 0 {
		return &ProcessErrors{
			Status:  400,
			Message: msg,
		}
	}

	return &ProcessErrors{
		Status:  400,
		Message: []string{"Bad Request!"},
	}
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
		Message: []string{"Not Found!"},
	}
}

func UnprocessableEntity(msg ...string) error {
	if len(msg) > 0 {
		return &ProcessErrors{
			Status:  422,
			Message: msg,
		}
	}

	return &ProcessErrors{
		Status:  422,
		Message: []string{"Unprocessable Entity!"},
	}
}

func InternalServer(msg ...string) error {
	if len(msg) > 0 {
		return &ProcessErrors{
			Status:  500,
			Message: msg,
		}
	}

	return &ProcessErrors{
		Status:  500,
		Message: []string{"Internal Server Error!"},
	}
}
