package main

import (
	"github.com/aminkbi/microChatApp/api/handler"
	"github.com/aminkbi/microChatApp/api/middleware"
	"net/http"
)

import "github.com/julienschmidt/httprouter"

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc("POST", "/v1/register", handler.Register)
	router.HandlerFunc("POST", "/v1/login", handler.Login)

	router.HandlerFunc("GET", "/v1/messages", middleware.AuthMiddleware(handler.ListMessages))
	router.HandlerFunc("POST", "/v1/messages", middleware.AuthMiddleware(handler.AddMessage))

	router.HandlerFunc("GET", "/v1/rooms", middleware.AuthMiddleware(handler.ListRooms))
	router.HandlerFunc("POST", "/v1/rooms", middleware.AuthMiddleware(handler.AddRoom))

	return router
}
