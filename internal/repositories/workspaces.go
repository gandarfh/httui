package repositories

import (
	"github.com/gandarfh/maid-san/external/database"
	"gorm.io/gorm"
)

type Workspace struct {
	gorm.Model
	Name      string     `json:"name"`
	Uri       string     `json:"uri"`
	Resources []Resource `json:"resources"`
}

type WorkspacesRepo struct {
	Sql *gorm.DB
}

func NewWorkspace() (*WorkspacesRepo, error) {
	db, err := database.SqliteConnection()
	db.AutoMigrate(&Workspace{})

	return &WorkspacesRepo{db}, err
}

func (repo *WorkspacesRepo) Create(value *Workspace) error {
	db := repo.Sql.Create(value)
	return db.Error
}

func (repo *WorkspacesRepo) Update(workspace *Workspace, value *Workspace) error {
	db := repo.Sql.Model(workspace).Updates(value)
	return db.Error
}

func (repo *WorkspacesRepo) FindOne(id uint) (Workspace, error) {
	workspace := Workspace{}

	db := repo.Sql.Model(&workspace).Where("id IS ?", id).First(&workspace)
	return workspace, db.Error
}

func (repo *WorkspacesRepo) List() ([]Workspace, error) {
	workspaces := []Workspace{}

	db := repo.Sql.Model(&workspaces).Find(&workspaces)
	return workspaces, db.Error
}

func (repo *WorkspacesRepo) Delete(id uint) error {
	db := repo.Sql.Model(&Workspace{}).Where("id IS ?", id).Delete(&Workspace{})
	return db.Error
}
