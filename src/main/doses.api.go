package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "strconv"
  "fmt"
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
