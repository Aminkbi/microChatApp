package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/aminkbi/microChatApp/internal/util"
	"github.com/aminkbi/microChatApp/websocket/handler"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// Client represents a single chat user
type Client struct {
	ID     string
	Conn   *websocket.Conn
	RoomID string
}

// Message represents a message in a chat room
type Message struct {
	data.MessageDTO
	Token string `json:"token"`
}

// ChatRoom manages the clients and broadcast messages
type ChatRoom struct {
	ID         string
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

// ChatServer manages multiple chat rooms
type ChatServer struct {
	Rooms map[string]*ChatRoom
}

// NewChatRoom creates a new ChatRoom
func NewChatRoom(id string) *ChatRoom {
	return &ChatRoom{
		ID:         id,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the chat room to handle register/unregister clients and broadcasting messages
func (room *ChatRoom) Run() {
	for {
		select {
		case client := <-room.Register:
			room.Clients[client] = true
		case client := <-room.Unregister:
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				client.Conn.Close()
			}
		case message := <-room.Broadcast:

			for client := range room.Clients {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					util.Logger.Printf("error: %v", err)
					client.Conn.Close()
					delete(room.Clients, client)
				}
			}
		}

	}
}

// ServeWs handles websocket requests from the peer
func (server *ChatServer) ServeWs(w http.ResponseWriter, r *http.Request) {
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
			util.Logger.Fatal(err)
		}
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.Logger.Println(err)
		return
	}

	client := &Client{ID: r.RemoteAddr, Conn: conn, RoomID: roomID}

	room, ok := server.Rooms[roomID]
	if !ok {
		room = NewChatRoom(roomID)
		server.Rooms[roomID] = room
		go room.Run()
	}

	room.Register <- client

	go server.handleMessages(client, room)
}

// handleMessages reads messages from the WebSocket connection and forwards them to the chat room
func (server *ChatServer) handleMessages(client *Client, room *ChatRoom) {
	defer func() {
		room.Unregister <- client
		client.Conn.Close()
	}()

	for {
		var msg Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			util.Logger.Printf("error: %v", err)
			break
		}

		err = sendMessageToAPI(msg)
		if err != nil {
			util.Logger.Printf("error: %v", err)
			client.Conn.WriteJSON(err)
			break
		}

		room.Broadcast <- msg
	}
}

func main() {

	port := os.Getenv("WEBSOCKET_PORT")
	if port == "" {
		util.Logger.Fatal("Please provide WEBSOCKET_PORT in env")
	}

	util.InitLogger()

	err := util.ConnectMongoDB()
	if err != nil {
		util.Logger.Fatal("Can't connect to mongo: ", err)
	}

	r := mux.NewRouter()
	server := &ChatServer{Rooms: make(map[string]*ChatRoom)}

	r.HandleFunc("/ws/{roomID}", server.ServeWs).Methods("GET")
	r.HandleFunc("/ws/{roomID}", handler.RoomCheck).Methods("POST")

	http.Handle("/", r)
	util.Logger.Printf("Server started on :%s", port)

	util.Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func sendMessageToAPI(message Message) error {
	url := fmt.Sprintf("http://go-app:%s/v1/messages", os.Getenv("APP_PORT"))

	tobeSent := data.MessageDTO{
		Content:  message.Content,
		SenderID: message.SenderID,
		RoomID:   message.RoomID,
	}
	jsonData, err := json.Marshal(tobeSent)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+message.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		util.Logger.Println(resp.StatusCode)
		return fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	return nil
}
