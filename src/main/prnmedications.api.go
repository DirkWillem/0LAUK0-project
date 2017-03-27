package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"main/utils"
	"net/http"
	"strconv"
)

// HandleCreatePRNMedication handles the creation of a new PRN medication
func HandleCreatePRNMedication(w http.ResponseWriter, r *http.Request) {
	// Read user ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read new medication from the request body
	var newMedication NewPRNMedication
	err = utils.ReadJSONFromRequest(r, &newMedication)

	if err != nil {
		utils.WriteError(w, utils.BadRequestError(err))
		return
	}

	// Create the new medication and respond
	medication, err := CreatePRNMedication(userID, newMedication)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, medication)
}

// HandleListPRNMedications returns a list of all PRN medications for a user to the client
func HandleListPRNMedications(w http.ResponseWriter, r *http.Request) {
	// Read user ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	// Read medications from database
	medications, err := ListPRNMedications(userID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, medications)
}

// HandleReadPRNMedication returns a single PRN medication for a given user and PRN medication ID
func HandleReadPRNMedication(w http.ResponseWriter, r *http.Request) {
	// Read user and PRN medication ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	medicationId, err := strconv.Atoi(vars["prnMedicationId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'prnMedicationId' isn't a valid integer.", vars["prnMedicationId"])))
		return
	}

	// Read medications from database
	medication, err := ReadPRNMedication(userID, medicationId)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, medication)
}

// HandleUpdatePRNMedication handles a PRN medication update
func HandleUpdatePRNMedication(w http.ResponseWriter, r *http.Request) {
	// Read user and medication ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	medicationId, err := strconv.Atoi(vars["prnMedicationId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'prnMedicationId' isn't a valid integer.", vars["prnMedicationId"])))
		return
	}

	// Read updated medication from request
	var updatedMedication UpdatedPRNMedication

	err = utils.ReadJSONFromRequest(r, &updatedMedication)

	if err != nil {
		utils.WriteError(w, utils.BadRequestError(err))
		return
	}

	// Update and return the updated medication
	medication, err := UpdatePRNMedication(userID, medicationId, updatedMedication)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, medication)
}

// HandleDeletePRNMedication handles the removal of a PRN medication
func HandleDeletePRNMedication(w http.ResponseWriter, r *http.Request) {
	// Read user and dose ID from the URL parameters
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'userId' isn't a valid integer.", vars["userId"])))
		return
	}

	medicationId, err := strconv.Atoi(vars["prnMedicationId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'prnMedicationId' isn't a valid integer.", vars["prnMedicationId"])))
		return
	}

	// Delete the dose and respond
	err = DeletePRNMedication(userID, medicationId)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
