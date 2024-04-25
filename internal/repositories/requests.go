package repositories

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model                                             // `json:"-"`
	Type        string                                     `json:"type"` // group | request
	Name        string                                     `json:"name"`
	Description string                                     `json:"description"`
	Method      string                                     `json:"method"`
	Endpoint    string                                     `json:"endpoint"`
	QueryParams datatypes.JSONType[[]map[string]string]    `json:"query_params"`
	Headers     datatypes.JSONType[[]map[string]string]    `json:"headers"`
	Body        datatypes.JSONType[map[string]interface{}] `json:"body"`
	ParentID    *uint                                      `json:"parent_id"`
}

type RequestsRepo struct {
	Sql *gorm.DB
}

func NewRequest() *RequestsRepo {
	db := Database

	return &RequestsRepo{db}
}

func (repo *RequestsRepo) Create(value *Request) error {
	if err := repo.Sql.Create(value).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RequestsRepo) Update(value *Request) error {
	if err := repo.Sql.Model(&Request{}).
		Where("id = ?", value.ID).
		Updates(value).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RequestsRepo) FindOne(id uint) (*Request, error) {
	request := Request{}

	if err := repo.Sql.Model(&request).Where("id = ?", id).First(&request).Error; err != nil {
		return &request, err
	}

	return &request, nil
}

func (repo *RequestsRepo) List(parentId *uint, filter string) ([]Request, error) {
	requests := []Request{}
	query := repo.Sql.Model(&requests)

	if filter != "" {
		query.Where("name LIKE ?", "%"+filter+"%")
	}

	if filter == "" {
		if parentId == nil {
			query = query.Where("parent_id IS NULL")
		} else {
			query = query.Where("parent_id = ?", *parentId)
		}
	}

	if err := query.Find(&requests).Error; err != nil {
		return requests, err
	}

	return requests, nil
}

func (repo *RequestsRepo) Delete(id uint) error {
	if err := repo.Sql.Model(&Request{}).
		Where("id = ?", id).
		Delete(&Request{}).Error; err != nil {
		return err
	}
	return nil
}
