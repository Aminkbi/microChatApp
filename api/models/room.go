package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Room contains Name of the room.
// It also contains a CreatedAt to trace the time of user being created.
type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required,unique"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
