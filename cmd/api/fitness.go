//Filename: cmd/api/fitness.go

package main

import (
	"fmt"
	"time"
	"net/http"

	"fitness.zioncastillo.net/internal/data"
	"fitness.zioncastillo.net/internal/validator"
)

func (app* application) saveFitnessHandler(w http.ResponseWriter, r *http.Request) {

	var input struct{
		UserId  int     `json:"user_id"`
		Steps   int     `json:"steps"`
		Cups    int     `json:"cups"`	
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//Copy the values from the input struct to a new fitness struct
	fitness := &data.Fitness{
		User_id: input.UserId,
		Steps: input.Steps,
		Cups: input.Cups,
	}

	//Initialize a new Vaalidator instance
	v := validator.New()

	//Check the map to determine if there were any validation errors
	if data.ValidateItem(v, fitness); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Fitness.Insert(fitness)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	//Create a Location header for the newly create resource/
	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/fitness/%d", fitness.ID))
	//Write the JSON response with 201 - Created status code with the body
	//Being the fitness data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"fitness": fitness}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listFitnessRecord(w http.ResponseWriter, r *http.Request) {

	//Create input struct to hold our query Parameters
	var input struct {
			ID	    int
			UserId  int
			Steps   int
			Cups    int
			Date 	time.Time
			data.Filters
	}

	//Initialize a validator
	v := validator.New()

	qf := r.URL.Query()

	//Use the helper methods to extract the values
	input.ID = app.readInt(qf, "id", 0, v)
	input.UserId = app.readInt(qf, "user_id", 0, v)
	input.Steps = app.readInt(qf, "steps", 0, v)
	input.Cups = app.readInt(qf, "cups", 0, v)

	//Get the page information
	input.Filters.Page = app.readInt(qf, "page", 1, v)
	input.Filters.PageSize = app.readInt(qf, "page_size", 20, v)

	//Get the page information
	input.Filters.Sort = app.readString(qf, "sort", "id")

	//Specify the allowed sort values
	input.Filters.SortList = []string{"id", "user_id", "steps", "cups", "date", "-id", "-user_id", "-steps", "-cups", "-date"}

	//Check for validation errors
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Get a listing of all fitness records
	lists, metadata, err := app.models.Fitness.GetAll(input.ID, input.UserId, input.Steps, input.Cups, input.Date, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//Send a JSON response containing all the fitnessRecord
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": lists, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}