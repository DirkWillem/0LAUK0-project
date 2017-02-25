package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
	"main/utils"
)

// HandleCreateDoseHistoryEntry handles the creation of a new dose history entry
func HandleCreateDoseHistoryEntry(w http.ResponseWriter, r *http.Request) {
	// Read user ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read new dose history entry from the request body
	var newDoseHistoryEntry NewDoseHistoryEntry
	err = utils.ReadJSONFromRequest(r, &newDoseHistoryEntry)

	if err != nil {
		utils.WriteError(w, utils.BadRequestError(err))
		return
	}

	// Create the new dose history entry and respond
	doseHistoryEntry, err := CreateDoseHistoryEntry(userID, newDoseHistoryEntry)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, doseHistoryEntry)
}

// HandleListDoseHistoryEntries returns a list of all dose history entries for a user to the client
func HandleListDoseHistoryEntries(w http.ResponseWriter, r *http.Request) {
	// Read user ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read dose history entries from database
	doses, err := ListDoseHistoryEntries(userID, map[string]string{
		"dose": r.URL.Query().Get("dose"),
		"dispensedday": r.URL.Query().Get("dispensedday"),
	})

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, doses)
}

// HandleReadDose returns a single dose history entry for a given user and dose history entry ID
func HandleReadDoseHistoryEntry(w http.ResponseWriter, r *http.Request) {
	// Read user and dose history entry ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	doseHistoryEntryID, err := strconv.Atoi(vars["doseHistoryEntryId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'doseHistoryEntryId' isn't a valid integer.", vars["doseHistoryEntryId"])))
		return
	}

	// Read dose history entry from database
	doseHistoryEntry, err := ReadDoseHistoryEntry(userID, doseHistoryEntryID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, doseHistoryEntry)
}