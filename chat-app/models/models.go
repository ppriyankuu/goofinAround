package models

import "gorm.io/gorm"

type ChatRoom struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

// type Message struct {
// 	gorm.Model
// 	RoomID  uint   `gorm:"not null"`
// 	Content string `gorm:"not null"`
// 	Sender  string `gorm:"not null"`
// }
