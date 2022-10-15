package dtos

type InputCreate struct {
	Template string `json:"template" validate:"required"`
	Path     string `json:"path" validate:"required"`
	Type     string `json:"type" validate:"required" default:".go"`
}

type InputSwagger struct {
	Path   string `json:"path" validate:"required"`
	Parent string `json:"Parent" validate:"required"`
}
