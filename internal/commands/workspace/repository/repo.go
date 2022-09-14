package repository

import (
	"fmt"

	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
)

type Workspaces struct {
	gorm.Model
	Name string `db:"name" validate:"required"`
	Uri  string `db:"uri" validate:"required"`
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

func (repo *WorkspaceRepo) Create(ws *Workspaces) {
	result := repo.Sql.Create(ws)

	if result.Error != nil {
		fmt.Println("Deu ruim criar dado")
	}
}

func (repo *WorkspaceRepo) List() *[]Workspaces {
	ws := []Workspaces{}
	result := repo.Sql.Find(&ws)

	if result.Error != nil {
		fmt.Println("Deu ruim criar dado")
	}

	return &ws
}
