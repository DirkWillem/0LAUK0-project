package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// HandleCreateDose handles the creation of a new dose
func HandleCreateDose(w http.ResponseWriter, r *http.Request) {
	// Read user ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read new dose from the request body
	var newDose NewDose
	err = ReadJSONFromRequest(r, &newDose)

	if err != nil {
		WriteError(w, BadRequestError(err))
		return
	}

	// Create the new dose and respond
	dose, err := CreateDose(userID, newDose)

	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, dose)
}

// HandleListDoses returns a list of all doses for a user to the client
func HandleListDoses(w http.ResponseWriter, r *http.Request) {
	// Read user ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read doses from database
	doses, err := ListDoses(userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, doses)
}

// HandleReadDose returns a single dose for a given user and dose ID
func HandleReadDose(w http.ResponseWriter, r *http.Request) {
	// Read user and dose ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	doseID, err := strconv.Atoi(vars["doseId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'doseId' isn't a valid integer.", vars["doseId"])))
		return
	}

	// Read doses from database
	dose, err := ReadDose(userID, doseID)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, dose)
}

// HandleUpdateDose handles a dose update
func HandleUpdateDose(w http.ResponseWriter, r *http.Request) {
	// Read user and dose ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	doseID, err := strconv.Atoi(vars["doseId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'doseId' isn't a valid integer.", vars["doseId"])))
		return
	}

	// Read updated dose from request
	var updatedDose UpdatedDose

	err = ReadJSONFromRequest(r, &updatedDose)

	if err != nil {
		WriteError(w, BadRequestError(err))
		return
	}

	// Update and return the updated dose
	dose, err := UpdateDose(userID, doseID, updatedDose)

	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, dose)
}

// HandleDeleteDose handles the removal of a doseda
func HandleDeleteDose(w http.ResponseWriter, r *http.Request) {
	// Read user and dose ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	doseID, err := strconv.Atoi(vars["doseId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'doseId' isn't a valid integer.", vars["doseId"])))
		return
	}

	// Delete the dose and respond
	err = DeleteDose(userID, doseID)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}