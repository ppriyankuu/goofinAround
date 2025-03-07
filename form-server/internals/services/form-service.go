package services

import (
	"form-server/internals/models"
	"form-server/internals/repositories"
)

type FormService interface {
	SubmitForm(formData string) (*models.FormResponse, error)
	GetFormByID(id uint) (*models.FormResponse, error)
	UpdateForm(form *models.FormResponse) error
	DeleteForm(id uint) error
}

type formService struct {
	repo repositories.FormRepo
}

func NewFormService(repo repositories.FormRepo) FormService {
	return &formService{repo: repo}
}

func (s *formService) SubmitForm(formData string) (*models.FormResponse, error) {
	form := &models.FormResponse{FormData: formData}
	if err := s.repo.Create(form); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *formService) GetFormByID(id uint) (*models.FormResponse, error) {
	return s.repo.GetByID(id)
}

func (s *formService) UpdateForm(form *models.FormResponse) error {
	return s.repo.Update(form)
}

func (s *formService) DeleteForm(id uint) error {
	return s.repo.Delete(id)
}
