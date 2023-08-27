package repositories

import (
	"github.com/gandarfh/httui/external/database"
	"gorm.io/gorm"
)

type DefaultsRepo struct {
	Sql *gorm.DB
}

type Default struct {
	gorm.Model
	WorkspaceId uint `json:"workspaceId"`
	TagId       uint `json:"tagId"`
}

func NewDefault() (*DefaultsRepo, error) {
	db := database.Client
	db.AutoMigrate(&Default{})

	items := []Default{}
	db.Find(&items)

	if len(items) == 0 {
		db.Create(&Default{})
	}

	return &DefaultsRepo{db}, nil
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
