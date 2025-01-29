package services

import (
	"basic-api/models"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

var users = map[string]models.User{
	"1": {ID: "1", Email: "user@example.com", Password: "topsecretpassword"},
}

func AuthenticateUser(email, password string) (*models.User, error) {
	for _, user := range users {
		if user.Email == email && user.Password == password {
			return &user, nil
		}
	}
	return nil, ErrInvalidCredentials
}

func GetUserByID(id string) (*models.User, error) {
	user, exists := users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return &user, nil
}
