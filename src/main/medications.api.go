package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "strconv"
  "fmt"
)

// HandleCreateMedication handles the creation of a medication
func HandleCreateMedication(w http.ResponseWriter, r *http.Request) {
  // Read the new medication from the request body
  var newMedication NewMedication

  err := ReadJSONFromRequest(r, &newMedication)
  if err != nil {
    WriteError(w, BadRequestError(err))
    return
  }

  // Create and return the medication
  medication, err := CreateMedication(newMedication)

  if err != nil {
    WriteError(w, err)
    return
  }

  WriteJSON(w, medication)
}

// HandleListMedications handles a read of all medications
func HandleListMedications(w http.ResponseWriter, r *http.Request) {
  // Read and return medications
  medications, err := ListMedications()

  if err != nil {
    WriteError(w, err)
    return
  }

  WriteJSON(w, medications)
}

func HandleReadMedication(w http.ResponseWriter, r *http.Request) {
  // Read medication ID from URL
  vars := mux.Vars(r)

  id, err := strconv.Atoi(vars["medicationId"])
  if err != nil {
    WriteError(w, BadRequestErrorMessage(fmt.Sprintf("Value '%s' of URL parameter 'medicationId' isn't a valid integer.", vars["medicationId"])))
    return
  }

  // Read the medication from the database and write to the client
  medication, err := ReadMedication(id)
  if err != nil {
    WriteError(w, err)
    return
  }

  WriteJSON(w, medication)
}