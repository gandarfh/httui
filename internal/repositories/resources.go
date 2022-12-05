package repositories

import (
	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Resource struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	TagId       uint           `json:"tagId"`
	Endpoint    string         `json:"endpoint"`
	Method      string         `json:"method"`
	QueryParams datatypes.JSON `json:"queryParams"`
	Headers     datatypes.JSON `json:"headers"`
	Body        datatypes.JSON `json:"body"`
}

type ResourcesRepo struct {
	Sql *gorm.DB
}

func NewResource() (*ResourcesRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Resource{})

	return &ResourcesRepo{db}, err
}

func (repo *ResourcesRepo) Create(value *Resource) error {
	db := repo.Sql.Create(value)
	return db.Error
}

func (repo *ResourcesRepo) Update(resource *Resource, value *Resource) error {
	db := repo.Sql.Model(resource)
	db.Updates(value)

	return db.Error
}

func (repo *ResourcesRepo) FindOne(id uint) (*Resource, error) {
	workspace := Resource{}

	db := repo.Sql.Model(&workspace).Where("id IS ?", id).First(&workspace)
	return &workspace, db.Error
}

func (repo *ResourcesRepo) List(tagId uint) ([]Resource, error) {
	resources := []Resource{}

	db := repo.Sql.Model(&resources).
		Where("tag_id IS ?", tagId).
		Find(&resources)

	return resources, db.Error
}

func (repo *ResourcesRepo) Delete(id uint) error {
	db := repo.Sql.Model(&Resource{}).Where("id IS ?", id).Delete(&Resource{})
	return db.Error
}
