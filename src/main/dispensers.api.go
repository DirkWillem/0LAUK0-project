package main

import (
	"net/http"
	"main/utils"
)

//HandleAuthenticateDispenser handles an authentication request from a dispenser
func HandleAuthenticateDispenser(w http.ResponseWriter, r *http.Request) {
	// Read auth token from request
	var auth DispenserAuth

	err := utils.ReadJSONFromRequest(r, &auth)

	if err != nil {
		utils.WriteError(w, utils.BadRequestError(err))
		return
	}

	// Try to authenticate the dispenser
	token, err := AuthenticateDispenser(auth)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	// Return session token to user
	utils.WriteJSON(w, token)
}