package main

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/hlandau/passlib.v1"
	"net/http"
	"main/utils"
)

type (
	// Credentials contains the credentials to authenticate a user
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// SessionToken containins a JWT token in the response
	SessionToken struct {
		Token string `json:"token"`
	}

	// Session contains all data the client needs on the current session
	Session struct {
		UserID   int
		Username string
		FullName string
		Role     string
		Email    string
	}
)

// UpdatePasswordHash updates the password hash if necessary
func UpdatePasswordHash(username, newPasswordHash string) error {
	_, err := db.Exec(`UPDATE Users
    SET PasswordHash = $1
    WHERE Username = $2`, newPasswordHash, username)

	if err != nil {
		utils.LogErrorMessage(err.Error())
	}

	return err
}

// JWTKeyFunc handles the decoding of a clients JWT
func JWTKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(config.JWT.Secret), nil
}

// ReadJWT reads the JWT session claimsfrom the request headers
func ReadJWTSession(r *http.Request) (Session, error) {
	// Read token string from the header
	tokStr := r.Header.Get("X-JWT")
	if len(tokStr) == 0 {
		return Session{}, utils.UnauthorizedErrorMessage("No X-JWT header was present")
	}

	// Parse token
	token, err := jwt.Parse(tokStr, JWTKeyFunc)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return Session{}, utils.UnauthorizedError(err)
		}
		return Session{}, utils.InternalServerError(err)
	}

	claims := token.Claims.(jwt.MapClaims)

	// Fill session record with the token claims
	return Session{
		UserID:   int(claims["userId"].(float64)),
		Username: claims["username"].(string),
		FullName: claims["fullName"].(string),
		Role:     claims["role"].(string),
		Email:    claims["email"].(string),
	}, nil
}

// CreateJWT creates a session token struct given a session struct
func CreateJWT(session Session) (SessionToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   session.UserID,
		"username": session.Username,
		"fullName": session.FullName,
		"role":     session.Role,
		"email":    session.Email,
	})

	tokenString, err := token.SignedString([]byte(config.JWT.Secret))

	return SessionToken{tokenString}, err
}

// Authenticate authenticates a user and returns a session token
func Authenticate(credentials Credentials) (SessionToken, error) {
	var token SessionToken

	// Select the user from the database
	var passwordHash string
	var session Session

	err := db.QueryRow(`SELECT ID, Username, FullName, PasswordHash, Role, Email
    FROM Users
    WHERE Username = $1`, credentials.Username).Scan(&session.UserID, &session.Username, &session.FullName, &passwordHash, &session.Role, &session.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return token, utils.NotFoundErrorMessage(fmt.Sprintf("Username %s not found", credentials.Username))
		} else {
			return token, utils.InternalServerError(err)
		}
	}

	// Check whether the password hashes match
	newPassHash, err := passlib.Verify(credentials.Password, passwordHash)
	if err != nil {
		return token, utils.InternalServerError(err)
	}

	// Update the password hash if necessary
	if newPassHash != "" {
		err := UpdatePasswordHash(credentials.Username, newPassHash)
		if err != nil {
			return token, utils.InternalServerError(err)
		}
	}

	// Return a signed JSON web token
	return CreateJWT(session)
}
