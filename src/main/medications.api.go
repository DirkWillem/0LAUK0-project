package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"main/utils"
)

// HandleCreateMedication handles the creation of a medication
func HandleCreateMedication(w http.ResponseWriter, r *http.Request) {
	// Read the new medication from the request body
	var newMedication NewMedication

	err := utils.ReadJSONFromRequest(r, &newMedication)
	if err != nil {
		utils.WriteError(w, utils.BadRequestError(err))
		return
	}

	// Create and return the medication
	medication, err := CreateMedication(newMedication)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, medication)
}

// HandleListMedications handles a read of all medications
func HandleListMedications(w http.ResponseWriter, r *http.Request) {
	// Read and return medications
	medications, err := ListMedications()

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, medications)
}

// HandleReadMedication returns data on a single medication to the client
func HandleReadMedication(w http.ResponseWriter, r *http.Request) {
	// Read medication ID from URL
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["medicationId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'medicationId' isn't a valid integer.", vars["medicationId"])))
		return
	}

	// Read the medication from the database and write to the client
	medication, err := ReadMedication(id)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, medication)
}

// HandleUpdateMedication handles the update of a single medication
func HandleUpdateMedication(w http.ResponseWriter, r *http.Request) {
	// Read medication ID from URL
	vars := mux.Vars(r)

	medicationID, err := strconv.Atoi(vars["medicationId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'medicationId' isn't a valid integer.", vars["medicationId"])))
		return
	}

	// Read updated medication from the request body
	var updatedMedication UpdatedMedication

	err = utils.ReadJSONFromRequest(r, &updatedMedication)
	if err != nil {
		utils.WriteError(w, utils.BadRequestError(err))
		return
	}

	// Update the medication and respond
	medication, err := UpdateMedication(medicationID, updatedMedication)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, medication)
}

// HandleDeleteMedication handles the removal of a single medication
func HandleDeleteMedication(w http.ResponseWriter, r *http.Request) {
	// Read medication ID from URL
	vars := mux.Vars(r)

	medicationID, err := strconv.Atoi(vars["medicationId"])
	if err != nil {
		utils.WriteError(w, utils.BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'medicationId' isn't a valid integer.", vars["medicationId"])))
		return
	}

	// Delete the medication and respond
	err = DeleteMedication(medicationID)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}