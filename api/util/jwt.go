package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = []byte("your-secret-key")

func CreateToken(username string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "micro-chat-app",                 // Issuer
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}
