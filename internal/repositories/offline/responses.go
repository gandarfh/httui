package offline

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Response struct {
	ID          string                          `gorm:"primaryKey" json:"id"`
	Status      string                          `json:"status"`
	WorkspaceId uint                            `json:"workspaceId"`
	RequestId   uint                            `json:"requestId"`
	Response    datatypes.JSONType[interface{}] `json:"response"`
	Request     datatypes.JSONType[Request]     `json:"request,omitempty"`
	CreatedAt   time.Time                       `json:"createdAt,omitempty"`
	UpdatedAt   time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt                  `gorm:"index" json:"deletedAt,omitempty"`
}

type ResponsesRepo struct {
	Sql *gorm.DB
}

func NewResponse() *ResponsesRepo {
	db := Database

	return &ResponsesRepo{db}
}

func (repo *ResponsesRepo) Create(value *Response) error {
	value.ID = primitive.NewObjectID().Hex()
	return repo.Sql.Create(value).Error
}

func (repo *ResponsesRepo) FindOne(requestId, workspace_id uint) (*Response, error) {
	request := Response{}
	err := repo.Sql.Model(&request).
		Where("request_id = ?", requestId).
		Where("workspace_id = ?", workspace_id).
		Order("created_at DESC").
		First(&request).Error

	return &request, err
}
