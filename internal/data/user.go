package data

import (
	"github.com/aminkbi/microChatApp/api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserModel struct {
	mongo *utils.MongoClient
}

// User contains personal information, Username, Email address and PasswordHash.
// It also contains a CreatedAt to trace the time of user being created.
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string             `bson:"username" json:"username" validate:"required,unique"`
	Email        string             `bson:"email" json:"email" validate:"required,email,unique"`
	PasswordHash string             `bson:"password_hash" json:"password_hash" validate:"required"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
