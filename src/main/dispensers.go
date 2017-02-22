package main

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/hlandau/passlib.v1"
)

type (
	// DispenserAuth contains the authentication credentials for a dispenser
	DispenserAuth struct {
		ID        int    `json:"id"`
		AuthToken string `json:"authToken"`
	}
)

// UpdateDispenserAuthToken updates the auth token of a dispenser
func UpdateDispenserAuthToken(id int, newToken string) error {
	_, err := db.Exec(`UPDATE Dispensers
	SET AuthToken = ?
	WHERE ID = ?`, newToken, id)

	if err != nil {
		LogErrorMessage(err.Error())
	}

	return err
}

// CreateDispenserJWT
func CreateDispenserJWT(id int) (SessionToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   id,
		"username": "",
		"fullName": "",
		"role":     "dispenser",
		"email":    "",
	})
	tokenString, err := token.SignedString([]byte(config.JWT.Secret))

	return SessionToken{tokenString}, err
}

// AuthenticateDispenser creates a JSON web token for a
func AuthenticateDispenser(auth DispenserAuth) (SessionToken, error) {
	var token SessionToken

	// Read the dispenser from the database
	var tokenHash string

	err := db.QueryRow(`SELECT AuthToken FROM Dispensers WHERE ID = ?`, auth.ID).Scan(&tokenHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return token, NotFoundErrorMessage(fmt.Sprintf("No dispenser with ID %d", auth.ID))
		} else {
			return token, InternalServerError(err)
		}
	}

	// Check whether the password hashes match
	newTokenHash, err := passlib.Verify(auth.AuthToken, tokenHash)
	if err != nil {
		return token, InternalServerError(err)
	}

	// Update the auth token if necessary
	if newTokenHash != "" {
		err := UpdateDispenserAuthToken(auth.ID, newTokenHash)
		if err != nil {
			return token, InternalServerError(err)
		}
	}

	return CreateDispenserJWT(auth.ID)
}
