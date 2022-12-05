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

func (repo *TagsRepo) List(workspaceId uint) ([]Tag, error) {
	tags := []Tag{}

	db := repo.Sql.Model(&tags).
		Where("workspace_id IS ?", workspaceId).
		Find(&tags)

	return tags, db.Error
}

func (repo *TagsRepo) Create(value *Tag) error {
	db := repo.Sql.Create(value)
	return db.Error
}

func (repo *TagsRepo) Update(tag *Tag, value *Tag) error {
	db := repo.Sql.Model(tag)
	db.Updates(value)

	return db.Error
}

func (repo *TagsRepo) Delete(id uint) error {
	db := repo.Sql.Model(&Tag{}).Where("id IS ?", id)
	db.Delete(&Tag{})

	return db.Error
}
