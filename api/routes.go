package main

import (
	"github.com/aminkbi/microChatApp/api/handler"
	"net/http"
)

import "github.com/julienschmidt/httprouter"

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc("POST", "/v1/register", handler.Register)
	router.HandlerFunc("POST", "/v1/login", handler.Login)

	return router
}
