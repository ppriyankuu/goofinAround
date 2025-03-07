package models

import "time"

type FormResponse struct {
	ID        uint `gorm:"primay_key"`
	FormData  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
