package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Read new user from request
	var newUser NewUser

	err := ReadJSONFromRequest(r, &newUser)
	if err != nil {
		WriteError(w, BadRequestError(err))
		return
	}

	// Create the user
	user, err := CreateUser(newUser)

	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, user)
}

// HandleListUsers returns a list of all users to the client
func HandleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ListUsers(map[string]string{
		"role": r.URL.Query().Get("role"),
	})

	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, users)
}

// HandleReadUser returns a single user to the client
func HandleReadUser(w http.ResponseWriter, r *http.Request) {
	// Read user ID from URL
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read the user from the database and write to the client
	user, err := ReadUser(userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, user)
}

// HandleUpdateUser handles the update of a user
func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Read user ID from URL
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read updated user from url
	var updatedUser UpdatedUser

	err = ReadJSONFromRequest(r, &updatedUser)
	if err != nil {
		WriteError(w, BadRequestError(err))
		return
	}

	// Update user
	user, err := UpdateUser(userID, updatedUser)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, user)
}

// HandleDeleteUser handles the removal of a user
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Read user ID from URL
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Delete user and respond
	err = DeleteUser(userID)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}