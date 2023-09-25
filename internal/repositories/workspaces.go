package repositories

import (
	"github.com/gandarfh/httui/external/database"
	"gorm.io/gorm"
)

type Workspace struct {
	gorm.Model
	Name string `gorm:"unique" json:"name"`
	Envs []Env  `json:"envs" gorm:"foreignKey:WorkspaceId;constraint:Onupdate:CASCADE;"`
}

type WorkspacesRepo struct {
	Sql *gorm.DB
}

func NewWorkspace() *WorkspacesRepo {
	db := database.Client
	db.AutoMigrate(&Workspace{})

	return &WorkspacesRepo{db}
}

func (repo *WorkspacesRepo) Create(value *Workspace) error {
	db := repo.Sql.Create(value)
	return db.Error
}

func (repo *WorkspacesRepo) Update(value *Workspace) error {
	if err := repo.Sql.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Where("id = ?", value.ID).
		Updates(value).Error; err != nil {
		return err
	}

	return nil

}

func (repo *WorkspacesRepo) FindOne(id uint) (Workspace, error) {
	workspace := Workspace{}

	if err := repo.Sql.Model(&workspace).
		Where("id IS ?", id).First(&workspace).Error; err != nil {
		return workspace, err
	}

	return workspace, nil
}

func (repo *WorkspacesRepo) List() ([]Workspace, error) {
	workspaces := []Workspace{}

	if err := repo.Sql.Model(&workspaces).
		Find(&workspaces).Error; err != nil {
		return workspaces, err
	}

	return workspaces, nil
}

func (repo *WorkspacesRepo) Delete(id uint) error {
	db := repo.Sql.Model(&Workspace{}).Where("id IS ?", id).Delete(&Workspace{})
	return db.Error
}
