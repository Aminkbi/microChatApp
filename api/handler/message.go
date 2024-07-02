package handler

import (
	"context"
	"errors"
	"github.com/aminkbi/microChatApp/api/util"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func ListMessages(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RoomId primitive.ObjectID `json:"roomId"`
	}

	err := util.ReadJSON(w, r, &input)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	if input.RoomId.IsZero() {
		BadRequestResponse(w, r, errors.New("roomId must be provided"))
		return
	}

	coll := util.MongoDBClient.GetCollection("micro-chat", "message")
	filter := bson.M{"roomId": input.RoomId}
	findOptions := options.Find().SetSort(bson.D{{"timestamp", 1}})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()
	cur, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}
	defer cur.Close(context.TODO())

	var messages []data.Message
	for cur.Next(context.TODO()) {
		var message data.Message
		err = cur.Decode(&message)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}
		messages = append(messages, message)
	}

	if err = cur.Err(); err != nil {
		ServerErrorResponse(w, r, err)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, data.Envelope{"messages": messages}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}

}

func AddMessage(w http.ResponseWriter, r *http.Request) {

	var input data.MessageDTO

	err := util.ReadJSON(w, r, &input)
	if err != nil {
		BadRequestResponse(w, r, err)
	}

	v := validator.New()

	if data.ValidateMessageDTO(v, &input); !v.Valid() {
		FailedValidationResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "message")

	roomId, err := primitive.ObjectIDFromHex(input.RoomID)
	if err != nil {
		BadRequestResponse(w, r, err)
	}
	senderId, err := primitive.ObjectIDFromHex(input.SenderID)
	if err != nil {
		BadRequestResponse(w, r, err)
	}

	message := data.Message{
		Content:   input.Content,
		SenderID:  senderId,
		RoomID:    roomId,
		Timestamp: time.Now(),
	}

	_, err = coll.InsertOne(ctx, message)
	if err != nil {
		//if mongo.IsDuplicateKeyError(err) {
		//	BadRequestResponse(w, r, ErrDuplicateCredentials)
		//	return
		//}
		ServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
