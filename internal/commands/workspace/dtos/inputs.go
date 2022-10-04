package dtos

import (
	"github.com/gandarfh/maid-san/internal/commands/resources/dtos"
)

type InputCreate struct {
	Name string `json:"name" validate:"required"`
	Uri  string `json:"uri" validate:"required"`
}

type InputUpdate struct {
	Name      string             `json:"name"`
	Uri       string             `json:"uri"`
	Resources []dtos.InputUpdate `json:"resources"`
}
