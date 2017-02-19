package main

import "net/http"

// HandleAuthenticate handles an authentication request
func HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	// Read credentials from request
	var credentials Credentials

	err := ReadJSONFromRequest(r, &credentials)

	if err != nil {
		WriteError(w, BadRequestError(err))
		return
	}

	// Try to authenticate the user
	token, err := Authenticate(credentials)

	if err != nil {
		WriteError(w, err)
		return
	}

	// Return session token to user
	WriteJSON(w, token)
}
