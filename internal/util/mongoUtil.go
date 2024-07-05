package util

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type mongoClient struct {
	Client *mongo.Client
}

var MongoDBClient mongoClient

func ConnectMongoDB() error {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("set your 'MONGODB_URI' environment variable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	Logger.Println("Connected To Mongo")

	MongoDBClient = mongoClient{
		Client: client,
	}

	return nil
}

func (m mongoClient) GetCollection(database, collection string) *mongo.Collection {
	return m.Client.Database(database).Collection(collection)
}
