package data

import (
	"github.com/aminkbi/microChatApp/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User contains personal information, Username, Email address and PasswordHash.
// It also contains a CreatedAt to trace the time of user being created.
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"password_hash"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

type UserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateUserDTO(v *validator.Validator, user *UserDTO) {
	// Ensure either username or email is provided
	v.Check(user.Username != "" || user.Email != "", "username or email", "must be provided")

	// If email is provided, it must be valid
	if user.Email != "" {
		v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be valid")
	}

	// If username is provided, it must be at least 4 characters long
	if user.Username != "" {
		v.Check(len(user.Username) >= 4, "username", "must not be less than 4 characters")
	}

	// Validate password
	v.Check(user.Password != "", "password", "must be provided")
	v.Check(len(user.Password) >= 4, "password", "must not be less than 4 characters")
}
