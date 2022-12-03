package repositories

import (
	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Resources   []Resource `json:"resources"`
	WorkspaceId uint       `json:"workspaceId"`
}

type TagsRepo struct {
	Sql *gorm.DB
}

func NewTag() (*TagsRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Tag{})

	return &TagsRepo{db}, err
}

func (repo *TagsRepo) List() ([]Tag, error) {
	tags := []Tag{}

	db := repo.Sql.Model(&tags).Find(&tags)
	return tags, db.Error
}

func (repo *TagsRepo) Create(value *Tag) error {
	db := repo.Sql.Create(value)
	return db.Error
}
