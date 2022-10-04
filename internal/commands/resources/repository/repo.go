package repository

import (
	"encoding/json"

	"github.com/gandarfh/maid-san/external/database"
	"github.com/gandarfh/maid-san/internal/commands/resources/dtos"
	"gorm.io/gorm"
)

type Params struct {
	gorm.Model
	ResourcesId uint
	Key         string `db:"key" json:"key"`
	Value       string `db:"value" json:"value"`
}

type Headers struct {
	gorm.Model
	ResourcesId uint
	Key         string `db:"key" json:"key"`
	Value       string `db:"value" json:"value"`
}

type workspace struct {
	gorm.Model
	Name string `db:"name"`
	Uri  string `db:"uri"`
}

type Resources struct {
	gorm.Model
	WorkspacesId uint            `db:"name" json:"workspace_id"`
	Name         string          `db:"name" json:"name"`
	Endpoint     string          `db:"endpoint" json:"endpoint"`
	Method       string          `db:"method" json:"method"`
	Params       []Params        `db:"params" json:"params"`
	Headers      []Headers       `db:"headers" json:"headers"`
	Body         json.RawMessage `db:"body" json:"body"`
}

func (re *Resources) AfterUpdate(tx *gorm.DB) error {
	tx.Model(&Params{}).Where("resources_id IS ?", nil).Unscoped().Delete(&Params{})
	tx.Model(&Headers{}).Where("resources_id IS ?", nil).Unscoped().Delete(&Headers{})

	return nil
}

func (re *Resources) Parent() *workspace {
	wk := workspace{}
	repo, err := NewResourcesRepo()

	if err != nil {
		return nil
	}

	repo.Sql.First(&wk, re.WorkspacesId)

	return &wk
}

type ResourceRepo struct {
	Sql *gorm.DB
}

func NewResourcesRepo() (*ResourceRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Resources{})
	db.AutoMigrate(&Headers{})
	db.AutoMigrate(&Params{})

	if err != nil {
		return nil, err
	}

	return &ResourceRepo{
		Sql: db,
	}, nil
}

func (repo *ResourceRepo) Update(resource *Resources, value *dtos.InputUpdate) {

	params := []Params{}
	for key, value := range value.Params {
		params = append(params, Params{ResourcesId: resource.ID, Value: value.(string), Key: key})
	}

	headers := []Headers{}
	for key, value := range value.Headers {
		headers = append(headers, Headers{ResourcesId: resource.ID, Value: value.(string), Key: key})
	}

	data := Resources{
		Name:     value.Name,
		Endpoint: value.Endpoint,
		Method:   value.Method,
		Body:     value.Body,
		Params:   params,
		Headers:  headers,
	}

	db := repo.Sql.Model(resource).Session(&gorm.Session{FullSaveAssociations: true})

	db.Association("Headers").Replace(headers)
	db.Association("Params").Replace(params)

	db.Updates(data)
}

func (repo *ResourceRepo) Create(value *dtos.InputCreate) {
	resource := Resources{
		WorkspacesId: uint(value.ParentId),
		Name:         value.Name,
		Endpoint:     value.Endpoint,
		Method:       value.Method,
		Params:       []Params{},
		Headers:      []Headers{},
		Body:         nil,
	}

	repo.Sql.Create(&resource)
}

func (repo *ResourceRepo) Find(id uint) *Resources {
	value := Resources{}

	db := repo.Sql.Model(&Resources{})

	db.Preload("Headers")
	db.Preload("Params")

	db.First(&value, id)

	return &value
}

func (repo *ResourceRepo) Delete(id uint) {
	repo.Sql.Model(&Params{}).Where("resources_id IS ?", nil).Unscoped().Delete(&Params{})
	repo.Sql.Model(&Headers{}).Where("resources_id IS ?", nil).Unscoped().Delete(&Headers{})

	db := repo.Sql.Model(&Resources{})
	db.Where("id IS ?", id).Unscoped().Delete(&Resources{})
}

func (repo *ResourceRepo) FindByName(name string) *Resources {
	value := Resources{}

	db := repo.Sql.Model(&Resources{})

	db.Preload("Headers")
	db.Preload("Params")

	db.Where("name = ?", name)

	db.First(&value)

	return &value
}

func (repo *ResourceRepo) List() *[]Resources {
	list := []Resources{}
	repo.Sql.Find(&list)

	return &list
}
