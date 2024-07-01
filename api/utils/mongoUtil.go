package utils

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
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
		err := errors.New("set your 'MONGODB_URI' environment variable")
		if err != nil {
			return err, nil
		}
	}
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		return err, nil
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return err, nil
	}
	log.Println("Connected To Mongo")

	return nil, &MongoClient{Client: client}

}

func (m MongoClient) GetCollection(database, collection string) *mongo.Collection {
	return m.Client.Database(database).Collection(collection)
}
