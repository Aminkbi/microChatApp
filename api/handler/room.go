package handler

import (
	"context"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/util"
	"github.com/aminkbi/microChatApp/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func ListRooms(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "room")

	filter := bson.M{}

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}
	defer cur.Close(ctx)

	var rooms []data.Room

	for cur.Next(ctx) {
		var room data.Room
		err = cur.Decode(&room)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}
		rooms = append(rooms, room)
	}

	if err = cur.Err(); err != nil {
		ServerErrorResponse(w, r, err)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, data.Envelope{"rooms": rooms}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}

}

func AddRoom(w http.ResponseWriter, r *http.Request) {

	var input data.Room

	err := util.ReadJSON(w, r, &input)
	if err != nil {
		BadRequestResponse(w, r, err)
	}

	v := validator.New()

	if data.ValidateRoom(v, &input); !v.Valid() {
		FailedValidationResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "room")

	room := data.Room{
		Name:      input.Name,
		CreatedAt: time.Now(),
	}

	_, err = coll.InsertOne(ctx, room)
	if err != nil {
		//if mongo.IsDuplicateKeyError(err) {
		//	BadRequestResponse(w, r, err)
		//	return
		//}
		ServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
