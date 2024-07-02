package handler

import (
	"github.com/aminkbi/microChatApp/api/util"
	"github.com/aminkbi/microChatApp/internal/data"
	"net/http"
)

func logError(r *http.Request, err error) {

	util.Logger.Println(err)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := data.Envelope{"error": message}
	err := util.WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(500)
	}
}

func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	errorResponse(w, r, http.StatusUnauthorized, message)
}
