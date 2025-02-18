package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupID   string             `bson:"group_id" json:"group_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Content   string             `bson:"content" json:"content"`
	Timestamp int64              `bson:"timestamp" json:"timestamp"`
}
