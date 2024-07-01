package util

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type MongoClient struct {
	Client *mongo.Client
}

func ConnectMongoDB() (error, *MongoClient) {
	if err := godotenv.Load(); err != nil {
		return err, nil
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("set your 'MONGODB_URI' environment variable"), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err, nil
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return err, nil
	}

	Logger.Println("Connected To Mongo")

	return nil, &MongoClient{Client: client}
}
func (m MongoClient) GetCollection(database, collection string) *mongo.Collection {
	return m.Client.Database(database).Collection(collection)
}
