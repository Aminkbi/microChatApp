package middleware

import (
	"context"
	"errors"
	"github.com/aminkbi/microChatApp/api/handler"
	"github.com/aminkbi/microChatApp/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handler.InvalidCredentialsResponse(w, r)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return util.SecretKey, nil
		})
		if err != nil || !token.Valid {
			handler.BadRequestResponse(w, r, errors.New("token is expired. please login again"))
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			handler.BadRequestResponse(w, r, err)
		}

		ctx := context.WithValue(r.Context(), "user", sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
