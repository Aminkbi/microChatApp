package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aminkbi/microChatApp/api/handler"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListRooms(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "room")

	_, err := coll.DeleteMany(ctx, bson.M{"name": "Test Room"})
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	// Insert a sample room
	room := data.Room{Name: "Test Room", CreatedAt: time.Now()}
	_, err = coll.InsertOne(ctx, room)
	if err != nil {
		t.Fatalf("Failed to insert sample room: %v", err)
	}

	// Generate a token
	token, err := util.CreateToken("6686d1e67b975c3fcbc6a69c")
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ListRooms(w, req)

	// Check the response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, res.StatusCode)
	}

	// Check the response body for the list of rooms
	var envelope struct {
		Rooms []data.Room `json:"rooms"`
	}
	err = json.NewDecoder(res.Body).Decode(&envelope)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(envelope.Rooms) == 0 {
		t.Errorf("Expected rooms in the response, got %v", envelope.Rooms)
	}
}

func TestAddRoom(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "room")

	// Clean up the database before testing
	_, err := coll.DeleteMany(ctx, bson.M{"name": "New Test Room"})
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	// Prepare input
	input := data.Room{Name: "New Test Room"}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	// Generate a token
	token, err := util.CreateToken("sampleUserID")
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/rooms", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.AddRoom(w, req)

	// Check the response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, res.StatusCode)
	}

	// Check if the room was created in the database
	var room data.Room
	err = coll.FindOne(ctx, bson.M{"name": "New Test Room"}).Decode(&room)
	if err != nil {
		t.Fatalf("Failed to find room in database: %v", err)
	}

	if room.Name != "New Test Room" {
		t.Errorf("Room data does not match: got %+v", room)
	}
}
