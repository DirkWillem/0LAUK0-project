package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
)

// HandleListDoseSummaries returns all dose summaries for a given user to a client
func HandleListDoseSummaries(w http.ResponseWriter, r *http.Request) {
	// Read user and dose history entry ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// List the dose summaries and return to the client
	summaries, err := ListDoseSummaries(userID)

	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, summaries)
}

// HandleReadDoseSummary returns a dose summary to the client
func HandleReadDoseSummary(w http.ResponseWriter, r *http.Request) {
	// Read user and dose history entry ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read the dose summary and return to the client
	statuses, err := ReadDoseSummary(userID, vars["date"])

	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, statuses)
}