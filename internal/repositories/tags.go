package repositories

import (
	"github.com/gandarfh/httui/external/database"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Resources   []Resource `json:"resources" gorm:"foreignKey:TagId;constraint:OnUpdate:CASCADE;"`
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

func (repo *TagsRepo) FindOne(tagId uint) (Tag, error) {
	tag := Tag{}

	db := repo.Sql.Model(&tag).
		Where("id IS ?", tagId).
		First(&tag)

	return tag, db.Error
}

func (repo *TagsRepo) FindOneByname(name string, workspaceId uint) (Tag, error) {
	tag := Tag{
		Name:        name,
		WorkspaceId: workspaceId,
	}

	if err := repo.Sql.Model(&tag).
		Where("name = ? AND workspace_id = ?", name, workspaceId).
		FirstOrCreate(&tag).Error; err != nil {
		return tag, err
	}

	return tag, nil
}

func (repo *TagsRepo) List(workspaceId uint) ([]Tag, error) {
	tags := []Tag{}

	db := repo.Sql.Model(&tags).
		Preload("Resources").
		Where("workspace_id IS ?", workspaceId).
		Find(&tags)

	return tags, db.Error
}

func (repo *TagsRepo) Create(value *Tag) error {
	db := repo.Sql.Create(value)
	return db.Error
}

func (repo *TagsRepo) Update(value *Tag) error {
	if err := repo.Sql.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Where("id = ?", value.ID).
		Updates(value).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TagsRepo) Delete(id uint) error {
	db := repo.Sql.Model(&Tag{}).Where("id IS ?", id)
	db.Delete(&Tag{})

	return db.Error
}
