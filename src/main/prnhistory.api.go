package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"main/utils"
	"net/http"
	"strconv"
)

// HandleCreatePRNHistoryEntry handles the creation of a new PRN history entry
func HandleCreatePRNHistoryEntry(w http.ResponseWriter, r *http.Request) {
	// Read user ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read new dose history entry from the request body
	var newPRNHistoryEntry NewPRNHistoryEntry
	err = utils.ReadJSONFromRequest(r, &newPRNHistoryEntry)

	if err != nil {
		utils.WriteError(w, utils.BadRequestError(err))
		return
	}

	// Create the new dose history entry and respond
	err = CreatePRNHistoryEntry(userID, newPRNHistoryEntry)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
