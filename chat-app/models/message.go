package models

import (
	"time"

	"gorm.io/gorm"
)

// Message represents a single chat message
type Message struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text"`
	Sender    string `gorm:"type:varchar(100)"`
	Room      string `gorm:"type:varchar(100)"`
	RoomID    uint   `gorm:"not null"` // Ensure this is linked to the room ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SaveMessage saves a new message to the database
func SaveMessage(db *gorm.DB, content, sender, room string, roomID uint) error {
	message := Message{
		Content: content,
		Sender:  sender,
		Room:    room,
		RoomID:  roomID, // Set room ID
	}
	return db.Create(&message).Error
}

// GetMessages fetches the message history for a given room
func GetMessages(db *gorm.DB, room string) ([]Message, error) {
	var messages []Message
	err := db.Where("room = ?", room).Order("created_at asc").Find(&messages).Error
	return messages, err
}
