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
	repo.Sql.Create(ws)
}

func (repo *WorkspaceRepo) List() *[]Workspaces {
	ws := []Workspaces{}
	repo.Sql.Find(&ws)

	return &ws
}
