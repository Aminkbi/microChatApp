package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/aminkbi/microChatApp/api/util"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var ErrDuplicateEmail = errors.New("email already exists")

func Register(w http.ResponseWriter, r *http.Request) {
	var input data.UserDTO

	err := util.ReadJSON(w, r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUserDTO(v, &input); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	newUser := data.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	_, err = util.MongoDBClient.GetCollection("micro-chat", "user").InsertOne(context.Background(), newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			badRequestResponse(w, r, errors.New("username or email already exists"))
			return
		}
		serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input data.UserDTO

	err := util.ReadJSON(w, r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	var user data.User
	if input.Username != "" {
		err = util.MongoDBClient.GetCollection("micro-chat", "user").FindOne(context.Background(), bson.M{"username": input.Username}).Decode(&user)
	} else if input.Email != "" {
		err = util.MongoDBClient.GetCollection("micro-chat", "user").FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	}
	if err != nil {
		InvalidCredentialsResponse(w, r)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		InvalidCredentialsResponse(w, r)
		return
	}

	token, err := util.CreateToken(user.Username)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	w.Write([]byte(fmt.Sprintf("Hello, %s!", user)))
}
