package models

type Message struct {
	ID        string `bson:"_id,omitempty"`
	UserID    string `bson:"user_id"`
	RoomID    string `bson:"room_id"`
	Content   string `bson:"content"`
	CreatedAt int64  `bson:"created_at"`
}
