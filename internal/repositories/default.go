package repositories

import (
	"gorm.io/gorm"
)

type DefaultsRepo struct {
	Sql *gorm.DB
}

type Default struct {
	gorm.Model
	WorkspaceId uint `json:"workspaceId"`
	RequestId   uint `json:"requestId"`
}

func NewDefault() *DefaultsRepo {
	items := []Default{}
	Database.Find(&items)

	if len(items) == 0 {
		Database.Create(&Default{})
	}

	return &DefaultsRepo{Database}
}

func (repo *DefaultsRepo) Update(value *Default) error {
	db := repo.Sql.Model(&Default{}).Where("id IS ?", 1).Updates(value)
	return db.Error
}

func (repo *DefaultsRepo) First() (*Default, error) {
	value := Default{}

	db := repo.Sql.Model(&value).First(&value)
	return &value, db.Error
}
