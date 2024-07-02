package util

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var SecretKey = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(username string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "micro-chat-app",                 // Issuer
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	// Print information about the created token
	Logger.Println("Token claims added: %+v\n", claims)
	return tokenString, nil
}
