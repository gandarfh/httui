package dtos

type InputCreate struct {
	Name string `json:"name" validate:"required"`
	Uri  string `json:"uri" validate:"required"`
}
