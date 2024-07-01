package main

import (
	"net/http"
)

import "github.com/julienschmidt/httprouter"

func (app *application) routes() http.Handler {

	router := httprouter.New()

	return router
}
