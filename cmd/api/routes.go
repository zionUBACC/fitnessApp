// Filename: cmd/api/routes.go

package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)
func (app *application) routes () *httprouter.Router{
	
	router := httprouter.New()
	
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/insert", app.saveFitnessHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	return router
}