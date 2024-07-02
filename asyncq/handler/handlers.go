package handler

import (
	"context"
	"github.com/aminkbi/microChatApp/asyncq/util"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func HandleArchiveMessagesTask(ctx context.Context, t *asynq.Task) error {
	util.Logger.Printf("Handling archive messages task %v", string(t.Payload()))

	coll := util.MongoDBClient.GetCollection("micro-chat", "message")
	archiveColl := util.MongoDBClient.GetCollection("micro-chat", "archiveMessage")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Calculate the date 1 minute ago
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)

	// Create a filter to find documents with a timestamp older than 1 minute
	filter := bson.M{"timestamp": bson.M{"$lt": oneMinuteAgo}}

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())

	var messages []interface{}
	for cur.Next(context.TODO()) {
		var message data.Message
		err = cur.Decode(&message)
		if err != nil {
			return err
		}
		messages = append(messages, message)
	}

	if len(messages) == 0 {
		util.Logger.Println("No messages to archive")
		return nil
	}

	// Insert all fetched messages into archiveColl
	_, err = archiveColl.InsertMany(ctx, messages)
	if err != nil {
		return err
	}

	// Delete the inserted messages from the original collection
	_, err = coll.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	util.Logger.Println("Archived messages:", messages)

	if err = cur.Err(); err != nil {
		return err
	}

	return nil
}

func HandleReportMessagesTask(ctx context.Context, t *asynq.Task) error {
	util.Logger.Printf("Handling report messages task %v", string(t.Payload()))
	userColl := util.MongoDBClient.GetCollection("micro-chat", "user")
	messageColl := util.MongoDBClient.GetCollection("micro-chat", "message")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Fetch all users
	cur, err := userColl.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())

	var users []data.User
	for cur.Next(context.TODO()) {
		var user data.User
		err = cur.Decode(&user)
		if err != nil {
			return err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		util.Logger.Println("No users to send reports")
		return nil
	}

	// Iterate over users and count messages
	for _, user := range users {
		messageCount, err := messageColl.CountDocuments(ctx, bson.M{"userId": user.ID})
		if err != nil {
			return err
		}
		util.Logger.Printf("User ID: %v has sent %d messages\n", user.ID, messageCount)
	}

	return nil
}
