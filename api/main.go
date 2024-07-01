package main

import (
	"context"
	"github.com/aminkbi/microChatApp/api/utils"
	"github.com/aminkbi/microChatApp/internal/data"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	Models data.Models
}

func main() {

	conf := getConfig()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	err, mongoClient := utils.ConnectMongoDB()
	if err != nil {
		log.Fatal("Can't connect to mongo: ", err)
	}

	models := data.NewModels(mongoClient)

	app := &application{
		config: conf,
		logger: logger,
		Models: models,
	}

	defer func() {
		if err = mongoClient.Client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Server started on :4000")
	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}

}

func getConfig() config {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		log.Fatal("Set your 'ENVIRONMENT' environment variable. ")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Set your 'PORT' environment variable. ")
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	return config{
		port: intPort,
		env:  env,
	}
}
