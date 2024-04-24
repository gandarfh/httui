package repositories

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Response struct {
	RequestId  uint                                    `json:"request_id"`
	ExternalId string                                  `json:"external_id"`
	Method     string                                  `json:"method"`
	Endpoint   string                                  `json:"endpoint"`
	Url        string                                  `json:"url"`
	Status     string                                  `json:"status"`
	Params     datatypes.JSONType[[]map[string]string] `json:"query_params"`
	Headers    datatypes.JSONType[[]map[string]string] `json:"headers"`
	Response   datatypes.JSONType[interface{}]         `json:"response"`
	Body       datatypes.JSONType[interface{}]         `json:"body"`
	Curl       string                                  `json:"curl"`
	gorm.Model `json:"-"`
}

type ResponsesRepo struct {
	Sql *gorm.DB
}

func NewResponse() *ResponsesRepo {
	db := Database

	return &ResponsesRepo{db}
}

func (repo *ResponsesRepo) Create(value *Response) error {
	return repo.Sql.Create(value).Error
}

func (repo *ResponsesRepo) FindOne(requestId string) (*Response, error) {
	request := Response{}
	err := repo.Sql.Model(&request).
		Where("request_id = ?", requestId).
		Order("created_at DESC").
		First(&request).Error

	return &request, err
}
