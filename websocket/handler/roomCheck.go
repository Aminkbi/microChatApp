package handler

import (
	"context"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func RoomCheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]
	// Convert roomID to ObjectID
	oid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	roomColl := util.MongoDBClient.GetCollection("micro-chat", "room")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	roomFilter := bson.M{"_id": oid}

	var roomCheck data.Room
	rm := roomColl.FindOne(ctx, roomFilter)
	err = rm.Decode(&roomCheck)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Room not found", http.StatusNotFound)
		} else {
			log.Fatal(err)
		}
		return
	}

}
