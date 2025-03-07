package repositories

import (
	"form-server/internals/models"

	"gorm.io/gorm"
)

type FormRepo interface {
	Create(form *models.FormResponse) error
	GetByID(id uint) (*models.FormResponse, error)
	Update(form *models.FormResponse) error
	Delete(id uint) error
}

type formRepo struct {
	db *gorm.DB
}

func NewFormRepo(db *gorm.DB) FormRepo {
	return &formRepo{db: db}
}

func (r *formRepo) Create(form *models.FormResponse) error {
	return r.db.Create(form).Error
}

func (r *formRepo) GetByID(id uint) (*models.FormResponse, error) {
	var form models.FormResponse
	if err := r.db.First(&form, id).Error; err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *formRepo) Update(form *models.FormResponse) error {
	return r.db.Save(form).Error
}

func (r *formRepo) Delete(id uint) error {
	return r.db.Delete(&models.FormResponse{}, id).Error
}
