package models

import "time"

type Task struct {
	ID        uint      `gorm:"primaryKey"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	Result    string    `json:"result"`
	Retries   int       `json:"retries"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
