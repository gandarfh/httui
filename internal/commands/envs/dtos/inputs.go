package dtos

type InputCreate struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type InputUpdate struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
