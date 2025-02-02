package services

import (
	"api-testing/interfaces"
	"api-testing/models"
)

type UserService struct {
	db interfaces.DBInterface
}

func NewUserService(db interfaces.DBInterface) *UserService {
	return &UserService{db: db}
}

func (us *UserService) CreateUser(user *models.User) error {
	return us.db.Create(user).Error
}

func (us *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := us.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
