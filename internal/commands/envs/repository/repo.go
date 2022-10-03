package repository

import (
	"fmt"

	"github.com/gandarfh/maid-san/external/database"
	"github.com/gandarfh/maid-san/internal/commands/envs/dtos"
	"github.com/gandarfh/maid-san/pkg/errors"
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

func (repo *EnvsRepo) Update(resource *Envs, value *dtos.InputUpdate) {
	data := Envs{
		Key:   value.Key,
		Value: value.Value,
	}

	db := repo.Sql.Model(resource).Session(&gorm.Session{FullSaveAssociations: true})
	db.Updates(data)
}

func (repo *EnvsRepo) Find(id uint) *Envs {
	value := Envs{}
	db := repo.Sql.Model(&Envs{})
	db.First(&value, id)

	return &value
}

func (repo *EnvsRepo) FindByKey(key string) (*Envs, error) {
	value := Envs{}
	db := repo.Sql.Model(&Envs{})
	result := db.Where("key = ?", key).First(&value)

	if result.Error != nil {
		return nil, errors.NotFoundError()
	}

	return &value, nil
}

func (repo *EnvsRepo) List() *[]Envs {
	envs := []Envs{}
	result := repo.Sql.Find(&envs)

	if result.Error != nil {
		fmt.Println("Deu ruim criar dado")
	}

	return &envs
}
