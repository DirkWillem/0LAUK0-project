package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

// CheckJWT checks whether a valid JSON web token is present in the request headers
func CheckJWT(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check whether the token string is present
		tokStr := r.Header.Get("X-JWT")
		if len(tokStr) == 0 {
			WriteError(w, UnauthorizedErrorMessage("No X-JWT header was present"))
			return
		}

		// Check whether the token string is valid
		_, err := jwt.Parse(tokStr, JWTKeyFunc)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				WriteError(w, UnauthorizedError(err))
			} else {
				WriteError(w, InternalServerError(err))
			}
			return
		}

		next(w, r)
	}
}

// CheckRole checks whether the current user has the appropriate role to acces the resource
func CheckRole(roles string, next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read token from request
		session, err := ReadJWTSession(r)
		if err != nil {
			WriteError(w, err)
		}

		for _, role := range strings.Split(roles, ",") {
			if session.Role == role {
				next(w, r)
				return
			}
		}

		WriteError(w, UnauthorizedErrorMessage(fmt.Sprintf("Your role (%s) is not authorized to access this resource.", session.Role)))
	}
}
