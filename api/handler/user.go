package handler

import (
	"context"
	"errors"
	"github.com/aminkbi/microChatApp/api/util"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var ErrDuplicateCredentials = errors.New("email or username already exists")

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "user")
	_, err = coll.InsertOne(ctx, newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			badRequestResponse(w, r, ErrDuplicateCredentials)
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
	coll := util.MongoDBClient.GetCollection("micro-chat", "user")
	if input.Username != "" {
		err = coll.FindOne(context.Background(), bson.M{"username": input.Username}).Decode(&user)
	} else if input.Email != "" {
		err = coll.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
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

	headers := make(http.Header)
	headers.Set("Authorization", "Bearer "+token)
	err = util.WriteJSON(w, http.StatusOK, data.Envelope{"token": token}, headers)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
