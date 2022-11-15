//Filename: cmd/api/users.go

package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"fitness.zioncastillo.net/internal/data"
	"fitness.zioncastillo.net/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	//Hold data from the request body
	var input struct {
		Name  	 string `json:"name"`
		Email 	 string `json:"email"`
		Password string `json:"password"`
	}

	//Parse the request body into the anonymous struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//copy the data to a new struct
	user:= &data.User{
		Name: input.Name,
		Email: input.Email,
		Activated: false,
	}

	//Generate a password hash
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//Perform validation on user input
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Insert the data in the database
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//Generate a token for the newly reated user
	token, err := app.models.Tokens.New(user.ID, 1*24*time.Hour, data.ScopeActivation)	
	if err != nil {
	app.serverErrorResponse(w, r, err)
	return
}
	app.background(func() {
		data := map[string]interface{} {
			"activationToken" : token.Plaintext,
			"userID" : user.ID,
		}
		//Send the email to the new user
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			//log the error
			app.logger.PrintError(err, nil)
		}
	})

	//write a 202 Accepted status
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Bcakground accepts a function as it parameter
func (app *application) background(fn func()) {
	//Increment the waitgroup counter
	app.wg.Add(1)
	//Launch/Create a goroutine which runs an anonymous
	//function that sends the welcome message
	go func() {
		defer app.wg.Done()
		//Recover from panics
		defer func () {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s",err), nil)
			}
		}()
		//execute the function
		fn()
	}()
}

func(app *application) activateUserHandler(w http.ResponseWriter, r *Request) {

	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//Perform validation
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	//Get the client details of the provided token or give the client feedback about an invalid token
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//Update the user status
	user.Activated = true
	//Save the updated user's record in our database
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//Delete the user's token that was used for activation
	err = app.models.Tokens.DeleteAllForUsers(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}