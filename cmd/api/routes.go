// Filename: cmd/api/routes.go

package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)
func (app *application) routes () http.Handler{
	
	router := httprouter.New()
	
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/records/insert", app.requirePermission("dailyfitness:read", app.saveFitnessHandler))
	router.HandlerFunc(http.MethodGet, "/v1/records/show", app.requirePermission("dailyfitness:read", app.listFitnessHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}