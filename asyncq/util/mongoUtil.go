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

type mongoClient struct {
	Client *mongo.Client
}

var MongoDBClient mongoClient

func ConnectMongoDB() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

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

	//wc := writeconcern.New(writeconcern.WMajority())
	//userCollection = client.Database("micro-chat", options.Database().SetWriteConcern(wc)).Collection("user")
	//
	//// Create unique index for email field
	//EnsureIndexes()

	return nil
}

//func EnsureIndexes() {
//	indexModel := mongo.IndexModel{
//		Keys: bson.M{
//			"email": 1, // create index on `email` field
//		},
//		Options: options.Index().SetUnique(true),
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	_, err := userCollection.Indexes().CreateOne(ctx, indexModel)
//	if err != nil {
//		log.Fatalf("Could not create index: %v", err)
//	}
//}

func (m mongoClient) GetCollection(database, collection string) *mongo.Collection {
	return m.Client.Database(database).Collection(collection)
}
