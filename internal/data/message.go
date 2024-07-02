package data

import (
	"github.com/aminkbi/microChatApp/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Message contains a Content of type string, related to sender and room through SenderID and RoomID
// It also contains a Timestamp to trace the time of message being created.
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content   string             `bson:"content" json:"content" `
	SenderID  primitive.ObjectID `bson:"senderId" json:"senderId"`
	RoomID    primitive.ObjectID `bson:"roomId" json:"roomId" `
	Timestamp time.Time          `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

type MessageDTO struct {
	Content  string `bson:"content" json:"content" `
	SenderID string `bson:"senderId" json:"senderId"`
	RoomID   string `bson:"roomId" json:"roomId" `
}

func ValidateMessageDTO(v *validator.Validator, message *MessageDTO) {
	v.Check(message.Content != "", "content", "must be provided")
	v.Check(message.RoomID != "", "roomId", "must be provided")
	v.Check(message.SenderID != "", "senderId", "must be provided")
}
