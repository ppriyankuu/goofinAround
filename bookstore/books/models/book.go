package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string  `json:"title" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Description string  `json:"description"`
}
