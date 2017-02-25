package main

import (
	"net/http"
	"main/utils"
)

// HandleAuthenticate handles an authentication request
func HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	// Read credentials from request
	var credentials Credentials

	err := utils.ReadJSONFromRequest(r, &credentials)

	if err != nil {
		utils.WriteError(w, utils.BadRequestError(err))
		return
	}

	// Try to authenticate the user
	token, err := Authenticate(credentials)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	// Return session token to user
	utils.WriteJSON(w, token)
}
