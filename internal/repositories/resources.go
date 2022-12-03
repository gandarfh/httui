package repositories

import (
	"encoding/json"

	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Param struct {
	gorm.Model
	ResourceId string
	Key        string `json:"key"`
	Value      string `json:"value"`
}
type Header struct {
	gorm.Model
	ResourceId string
	Key        string `json:"key"`
	Value      string `json:"value"`
}

type Resource struct {
	gorm.Model
	WorkspaceId uint            `json:"workspaceId"`
	TagId       uint            `json:"tagId"`
	Name        string          `json:"name"`
	Endpoint    string          `json:"endpoint"`
	Method      string          `json:"method"`
	Params      []Param         `json:"params"`
	Headers     []Header        `json:"headers"`
	Body        json.RawMessage `json:"body"`
}

type ResourcesRepo struct {
	Sql *gorm.DB
}

func NewResource() (*ResourcesRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Resource{})
	db.AutoMigrate(&Param{})
	db.AutoMigrate(&Header{})

	return &ResourcesRepo{db}, err
}

func (repo *ResourcesRepo) Create(value *Resource) error {
	db := repo.Sql.Create(value)
	return db.Error
}

func (repo *ResourcesRepo) Update(workspace *Resource, value *Resource) error {
	db := repo.Sql.Model(workspace).
		Session(&gorm.Session{FullSaveAssociations: true})

	db.Association("Params").Replace(value.Params)
	db.Association("Headers").Replace(value.Headers)

	db.Updates(value)

	return db.Error
}

func (repo *ResourcesRepo) FindOne(id uint) (*Resource, error) {
	workspace := Resource{}

	db := repo.Sql.Model(&workspace).Where("id IS ?", id).First(&workspace)
	return &workspace, db.Error
}

func (repo *ResourcesRepo) ListByTagAndWorkspace(workspace uint) ([]Resource, error) {
	resources := []Resource{}

	db := repo.Sql.Model(&Resource{}).
		Preload(clause.Associations).
		Where("workspace_id = ?", workspace).
		Find(&resources)

	return resources, db.Error
}

func (repo *ResourcesRepo) Delete(id uint) error {
	db := repo.Sql.Model(&Resource{}).Where("id IS ?", id).Delete(&Resource{})
	return db.Error
}
