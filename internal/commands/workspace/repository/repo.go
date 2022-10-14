package repository

import (
	"fmt"

	"github.com/gandarfh/maid-san/external/database"
	resourcerepo "github.com/gandarfh/maid-san/internal/commands/resources/repository"
	"github.com/gandarfh/maid-san/internal/commands/workspace/dtos"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Workspaces struct {
	gorm.Model
	Name      string `db:"name" validate:"required"`
	Uri       string `db:"uri" validate:"required"`
	Resources []resourcerepo.Resources
}

func (re *Workspaces) AfterUpdate(tx *gorm.DB) error {
	tx.Model(&resourcerepo.Resources{}).Where("workspaces_id IS ?", nil).Unscoped().Delete(&resourcerepo.Resources{})

	return nil
}

type WorkspaceRepo struct {
	Sql *gorm.DB
}

func NewWorkspaceRepo() (*WorkspaceRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Workspaces{})

	if err != nil {
		fmt.Println("Deu ruim database")
		return nil, err
	}

	return &WorkspaceRepo{
		Sql: db,
	}, nil
}

func (repo *WorkspaceRepo) Create(ws *dtos.InputCreate) {
	value := Workspaces{
		Uri:  ws.Uri,
		Name: ws.Name,
	}

	repo.Sql.Omit(clause.Associations).Create(&value)
}

func (repo *WorkspaceRepo) Delete(id uint) {
	db := repo.Sql.Model(&Workspaces{})
	db.Where("id IS ?", id).Unscoped().Delete(&Workspaces{})
}

func (repo *WorkspaceRepo) List() *[]Workspaces {
	ws := []Workspaces{}
	db := repo.Sql.Model(&Workspaces{})
	db.Preload("Resources")
	db.Find(&ws)

	return &ws
}

func (repo *WorkspaceRepo) FindByName(name string) *Workspaces {
	value := Workspaces{}

	db := repo.Sql.Model(&Workspaces{})

	db.Preload("Resources")
	db.Preload("Resources.Headers")
	db.Preload("Resources.Params")

	db.Where("name = ?", name)

	db.First(&value)

	return &value
}

func (repo *WorkspaceRepo) Update(workspace *Workspaces, value *dtos.InputUpdate) {
	resources := []resourcerepo.Resources{}
	for _, item := range value.Resources {
		params := []resourcerepo.Params{}
		for _, param := range item.Params {
			for key, value := range param {
				params = append(params, resourcerepo.Params{Key: key, Value: value.(string)})
			}
		}

		headers := []resourcerepo.Headers{}
		for key, value := range item.Headers {
			headers = append(headers, resourcerepo.Headers{Key: key, Value: value.(string)})
		}

		resources = append(
			resources,
			resourcerepo.Resources{
				WorkspacesId: item.WorkspacesId,
				Name:         item.Name,
				Endpoint:     item.Endpoint,
				Method:       item.Method,
				Body:         item.Body,
				Params:       params,
				Headers:      headers,
			})
	}

	data := Workspaces{
		Name:      value.Name,
		Uri:       value.Uri,
		Resources: resources,
	}

	db := repo.Sql.Model(workspace).Session(&gorm.Session{FullSaveAssociations: true})

	db.Association("Resources").Replace(resources)

	db.Updates(data)
}
