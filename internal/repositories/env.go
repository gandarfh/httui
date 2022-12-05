package repositories

import (
	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
)

type Env struct {
	gorm.Model
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvsRepo struct {
	Sql *gorm.DB
}

func NewEnvs() (*EnvsRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Env{})

	if err != nil {
		return nil, err
	}

	return &EnvsRepo{
		Sql: db,
	}, nil
}

func (repo *EnvsRepo) Create(env *Env) error {
	db := repo.Sql.Create(&Env{Value: env.Value, Key: env.Key})

	return db.Error
}

func (repo *EnvsRepo) Update(resource *Env, value *Env) error {
	data := Env{
		Key:   value.Key,
		Value: value.Value,
	}

	db := repo.Sql.Model(resource).Session(&gorm.Session{FullSaveAssociations: true})
	db.Updates(data)

	return db.Error
}

func (repo *EnvsRepo) Delete(id uint) error {
	db := repo.Sql.Model(&Env{})
	db.Where("id IS ?", id).Delete(&Env{})

	return db.Error
}

func (repo *EnvsRepo) Find(id uint) (Env, error) {
	value := Env{}
	db := repo.Sql.Model(&Env{})
	db.First(&value, id)

	return value, db.Error
}

func (repo *EnvsRepo) FindByKey(key string) (Env, error) {
	value := Env{Key: key}
	db := repo.Sql.Model(&Env{})
	db.Where("key = ?", key).FirstOrCreate(&value)

	return value, db.Error
}

func (repo *EnvsRepo) List() ([]Env, error) {
	envs := []Env{}
	db := repo.Sql.Find(&envs)

	return envs, db.Error
}
