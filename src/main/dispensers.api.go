package main

import "net/http"

//HandleAuthenticateDispenser handles an authentication request from a dispenser
func HandleAuthenticateDispenser(w http.ResponseWriter, r *http.Request) {
	// Read auth token from request
	var auth DispenserAuth

	err := ReadJSONFromRequest(r, &auth)

	if err != nil {
		WriteError(w, BadRequestError(err))
		return
	}

	// Try to authenticate the dispenser
	token, err := AuthenticateDispenser(auth)

	if err != nil {
		WriteError(w, err)
		return
	}

	// Return session token to user
	WriteJSON(w, token)
}