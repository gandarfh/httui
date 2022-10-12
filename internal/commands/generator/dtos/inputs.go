package dtos

type InputCreate struct {
	Template string `json:"template" validate:"required"`
	Path     string `json:"path" validate:"required"`
	Type     string `json:"type" validate:"required" default:".go"`
}

type InputSwagger struct {
	Link   string `json:"link" validate:"required"`
	Parent string `json:"Parent" validate:"required"`
}
