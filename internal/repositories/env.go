package repositories

import (
	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
)

type Envs struct {
	gorm.Model
  Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvsRepo struct {
	Sql *gorm.DB
}

func NewEnvsRepo() (*EnvsRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Envs{})

	if err != nil {
		return nil, err
	}

	return &EnvsRepo{
		Sql: db,
	}, nil
}

func (repo *EnvsRepo) Create(env *Envs) {
	repo.Sql.Create(&Envs{Value: env.Value, Key: env.Key})
}

func (repo *EnvsRepo) Update(resource *Envs, value *Envs) {
	data := Envs{
		Key:   value.Key,
		Value: value.Value,
	}

	db := repo.Sql.Model(resource).Session(&gorm.Session{FullSaveAssociations: true})
	db.Updates(data)
}

func (repo *EnvsRepo) Delete(id uint) {
	db := repo.Sql.Model(&Envs{})
	db.Where("id IS ?", id).Delete(&Envs{})
}

func (repo *EnvsRepo) Find(id uint) *Envs {
	value := Envs{}
	db := repo.Sql.Model(&Envs{})
	db.First(&value, id)

	return &value
}

func (repo *EnvsRepo) FindByKey(key string) *Envs {
	value := Envs{Key: key}
	db := repo.Sql.Model(&Envs{})
	db.Where("key = ?", key).FirstOrCreate(&value)

	return &value
}

func (repo *EnvsRepo) List() *[]Envs {
	envs := []Envs{}
	repo.Sql.Find(&envs)

	return &envs
}
