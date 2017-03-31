package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"main/utils"
	"net/http"
	"strconv"
)

// HandleListDoseSummaries returns all dose summaries for a given user to a client
func HandleListDoseSummaries(w http.ResponseWriter, r *http.Request) {
	// Read user and dose history entry ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// List the dose summaries and return to the client
	summaries, err := ListDoseSummaries(userID)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, summaries)
}

// HandleReadDoseSummary returns a dose summary to the client
func HandleReadDoseSummary(w http.ResponseWriter, r *http.Request) {
	// Read user and dose history entry ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read the dose summary and return to the client
	statuses, err := ReadDoseSummary(userID, vars["date"])

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, statuses)
}

// HandleReadPRNSummary returns a PRN summary to the client
func HandleReadPRNSummary(w http.ResponseWriter, r *http.Request) {
	// Read user and dose history entry ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read the dose summary and return to the client
	statuses, err := ReadPRNStatuses(userID, vars["date"])

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, statuses)
}
