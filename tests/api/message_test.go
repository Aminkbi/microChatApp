package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aminkbi/microChatApp/api/handler"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "message")

	roomID := primitive.NewObjectID()
	senderID := primitive.NewObjectID()

	// Clean up the database before testing
	_, err := coll.DeleteMany(ctx, bson.M{"senderId": senderID, "roomId": roomID})
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	message := data.Message{
		Content:  "Hello, world!",
		SenderID: senderID,
		RoomID:   roomID,
	}
	_, err = coll.InsertOne(ctx, message)
	if err != nil {
		t.Fatalf("Failed to insert sample message: %v", err)
	}

	// Prepare input
	input := struct {
		RoomId primitive.ObjectID `json:"roomId"`
	}{RoomId: roomID}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/messages", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ListMessages(w, req)

	// Check the response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, res.StatusCode)
	}

	// Check the response body for the list of messages
	var envelope struct {
		Messages []data.Message `json:"messages"`
	}
	err = json.NewDecoder(res.Body).Decode(&envelope)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(envelope.Messages) == 0 {
		t.Errorf("Expected messages in the response, got %v", envelope.Messages)
	}
}

func TestAddMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "message")

	roomID := primitive.NewObjectID()
	senderID := primitive.NewObjectID()

	// Clean up the database before testing
	_, err := coll.DeleteMany(ctx, bson.M{"roomId": roomID, "senderID": senderID})
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	// Prepare input
	input := data.MessageDTO{
		Content:  "Hello, world!",
		SenderID: senderID.Hex(),
		RoomID:   roomID.Hex(),
	}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.AddMessage(w, req)

	// Check the response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, res.StatusCode)
	}

	// Check if the message was created in the database
	var message data.Message
	err = coll.FindOne(ctx, bson.M{"content": "Hello, world!"}).Decode(&message)
	if err != nil {
		t.Fatalf("Failed to find message in database: %v", err)
	}

	if message.Content != "Hello, world!" {
		t.Errorf("Message data does not match: got %+v", message)
	}
}
