package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aminkbi/microChatApp/api/handler"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestRegister(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "user")

	// Clean up the database before testing
	_, err := coll.DeleteOne(ctx, bson.M{"username": "testUser"})
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	// Prepare input
	input := data.UserDTO{
		Username: "testUser",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	// Check the response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, res.StatusCode)
	}

	// Check if the user was created in the database
	var user data.User
	err = coll.FindOne(ctx, bson.M{"username": "testUser"}).Decode(&user)
	if err != nil {
		t.Fatalf("Failed to find user in database: %v", err)
	}

	if user.Username != "testUser" || user.Email != "test@example.com" {
		t.Errorf("User data does not match: got %+v", user)
	}
}

func TestLogin(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	coll := util.MongoDBClient.GetCollection("micro-chat", "user")

	// Clean up the database before testing
	_, err := coll.DeleteMany(ctx, bson.M{"username": "testUser"})
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	// Create a user for login testing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	newUser := data.User{
		Username:     "testUser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	_, err = coll.InsertOne(ctx, newUser)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Prepare input for login
	input := data.UserDTO{
		Username: "testUser",
		Password: "password123",
	}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the Login function
	handler.Login(w, req)

	// Check the response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, res.StatusCode)
	}

	// Check the response body for the token
	var envelope data.Envelope
	err = json.NewDecoder(res.Body).Decode(&envelope)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	token, ok := envelope["token"].(string)
	if !ok || token == "" {
		t.Errorf("Expected a token in the response, got %v", envelope)
	}
}
