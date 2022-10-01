package repository

import (
	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
)

type Resources struct {
	gorm.Model
}

type ResourceRepo struct {
	Sql *gorm.DB
}

func NewResourcesRepo() (*ResourceRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Resources{})

	if err != nil {
		return nil, err
	}

	return &ResourceRepo{
		Sql: db,
	}, nil
}

func (repo *ResourceRepo) Create(value *Resources) {
	repo.Sql.Create(value)
}

func (repo *ResourceRepo) List() *[]Resources {
	list := []Resources{}
	repo.Sql.Find(&list)

	return &list
}
