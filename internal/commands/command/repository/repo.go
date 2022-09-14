package repository

import (
	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
)

type Commands struct {
	gorm.Model
}

type CommandRepo struct {
	Sql *gorm.DB
}

func NewCommandRepo() (*CommandRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Commands{})

	if err != nil {
		return nil, err
	}

	return &CommandRepo{
		Sql: db,
	}, nil
}

func (repo *CommandRepo) Create(ws *Commands) {
	repo.Sql.Create(ws)
}

func (repo *CommandRepo) List() *[]Commands {
	ws := []Commands{}
	repo.Sql.Find(&ws)

	return &ws
}
