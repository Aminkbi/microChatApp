package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Message contains a Content of type string, related to sender and room through SenderID and RoomID
// It also contains a Timestamp to trace the time of message being created.
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content   string             `bson:"content" json:"content" `
	SenderID  primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	RoomID    primitive.ObjectID `bson:"room_id" json:"room_id" `
	Timestamp time.Time          `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}
