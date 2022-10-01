package dtos

import "encoding/json"

type InputCreate struct {
	ParentId int    `json:"parent_id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Endpoint string `json:"endpoint" validate:"required"`
	Method   string `json:"method" validate:"required"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type InputUpdate struct {
	WorkspacesId uint            `json:"workspace_id"`
	Name         string          `json:"name"`
	Endpoint     string          `json:"endpoint"`
	Method       string          `json:"method"`
	Params       []KeyValue      `json:"params"`
	Headers      []KeyValue      `json:"headers"`
	Body         json.RawMessage `json:"body"`
}
