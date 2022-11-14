package main

import (
	"fmt"
	"net/http"
)
func (app *application) logError(r *http.Request, err error){
	app.logger.Println(err)
}
//we want to send json-formatted error message
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	//create json response
	env := envelope{"error":message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
//server error response
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//we will log the error
	app.logError(r, err)
	//prepare a message with the error
	message := "the server encounter a problem and could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
//the not found response
func (app *application) notFoundResponse (w http.ResponseWriter, r *http.Request) {
	//create our message
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}
//a method not allowed response
func (app *application) methodNotAllowedResponse (w http.ResponseWriter, r *http.Request) {
	//create our message
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
//user provided a bad request
func (app *application) badRequestResponse (w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
//validation error
func (app * application) failedValidationResponse (w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}