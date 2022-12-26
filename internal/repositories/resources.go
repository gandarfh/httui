package repositories

import (
	"github.com/gandarfh/httui/external/database"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Resource struct {
	gorm.Model  `json:"-"`
	Name        string         `json:"name"`
	Method      string         `json:"method"`
	Endpoint    string         `json:"endpoint"`
	Description string         `json:"description"`
	QueryParams datatypes.JSON `json:"queryParams"`
	Headers     datatypes.JSON `json:"headers"`
	Body        datatypes.JSON `json:"body"`
	TagId       uint           `json:"tag_id"`
	Tag         Tag            `json:"-" gorm:"foreignKey:TagId"`
}

type ResourcesRepo struct {
	Sql *gorm.DB
}

func NewResource() (*ResourcesRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Resource{})

	return &ResourcesRepo{db}, err
}

func (repo *ResourcesRepo) Create(value *Resource) error {
	if err := repo.Sql.Create(value).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ResourcesRepo) Update(value *Resource) error {
	if err := repo.Sql.Model(&Resource{}).
		Where("id = ?", value.ID).
		Updates(value).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ResourcesRepo) FindOne(id uint) (*Resource, error) {
	resource := Resource{}

	if err := repo.Sql.Model(&resource).Where("id = ?", id).First(&resource).Error; err != nil {
		return &resource, err
	}

	return &resource, nil
}

func (repo *ResourcesRepo) List(tagId uint, filter string) ([]Resource, error) {
	resources := []Resource{}

	if err := repo.Sql.Model(&resources).
		Preload("Tag").
		Where("tag_id = ?", tagId).
		Where("name LIKE '%" + filter + "%'").
		Find(&resources).Error; err != nil {
		return resources, err
	}

	return resources, nil
}

func (repo *ResourcesRepo) Delete(id uint) error {
	if err := repo.Sql.Model(&Resource{}).
		Where("id = ?", id).
		Delete(&Resource{}).Error; err != nil {
		return err
	}
	return nil
}
