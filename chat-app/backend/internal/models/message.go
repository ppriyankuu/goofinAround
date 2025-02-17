package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    string             `bson:"room_id" json:"room_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Content   string             `bson:"content" json:"content"`
	Timestamp primitive.DateTime `bson:"timestamp" json:"timestamp"`
}
