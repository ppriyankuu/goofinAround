package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	BookID uint   `json:"book_id" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
	Status string `json:"status" binding:"required"` // e.g., "Pending", "Completed", "Cancelled"
}
