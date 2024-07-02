package main

import (
	"context"
	"github.com/aminkbi/microChatApp/api/util"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
}

func main() {
	util.InitLogger()

	conf := getConfig()
	err := util.ConnectMongoDB()
	if err != nil {
		util.Logger.Fatal("Can't connect to mongo: ", err)
	}

	app := &application{
		config: conf,
	}

	defer func() {
		if err = util.MongoDBClient.Client.Disconnect(context.TODO()); err != nil {
			util.Logger.Fatal(err)
		}
	}()

	util.Logger.Println("Server started on :4000")
	err = app.serve()
	if err != nil {
		util.Logger.Fatal(err)
	}

}

func getConfig() config {

	if err := godotenv.Load(); err != nil {
		util.Logger.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		util.Logger.Fatal("Set your 'JWT_SECRET' environment variable. ")
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		util.Logger.Fatal("Set your 'ENVIRONMENT' environment variable. ")
	}

	port := os.Getenv("PORT")
	if port == "" {
		util.Logger.Fatal("Set your 'PORT' environment variable. ")
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		util.Logger.Fatal(err)
	}

	return config{
		port: intPort,
		env:  env,
	}
}
