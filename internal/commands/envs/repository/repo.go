package repository

import (
	"fmt"

	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
)

type Envs struct {
	gorm.Model
	Key   string `db:"key"`
	Value string `db:"value"`
}

type EnvsRepo struct {
	Sql *gorm.DB
}

func NewEnvsRepo() (*EnvsRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Envs{})

	if err != nil {
		fmt.Println("Deu ruim database")
		return nil, err
	}

	return &EnvsRepo{
		Sql: db,
	}, nil
}

func (repo *EnvsRepo) Create(env *Envs) {
	result := repo.Sql.Create(env)

	if result.Error != nil {
		fmt.Println("Deu ruim criar dado")
	}
}

func (repo *EnvsRepo) List() *[]Envs {
	envs := []Envs{}
	result := repo.Sql.Find(&envs)

	if result.Error != nil {
		fmt.Println("Deu ruim criar dado")
	}

	return &envs
}
