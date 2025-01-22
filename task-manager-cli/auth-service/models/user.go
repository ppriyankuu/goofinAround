package models

import "task-manager/auth-service/database"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(user *User) error {
	db := database.GetDB()
	return db.Create(user).Error
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	db := database.GetDB()
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
