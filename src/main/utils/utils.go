package utils

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
)

type (
	ErrorMessage struct {
		Message string `json:"message"`
	}

	// Contains the minimal information on an entity
	MinimalEntity struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
)

// HashMD5 returns the MD5 hash of a string
func HashMD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// ReadJSONFromRequest unmarshals a JSON request body
func ReadJSONFromRequest(r *http.Request, target interface{}) error {
	reader := json.NewDecoder(r.Body)
	return reader.Decode(target)
}

// WriteJSON writes a JSON value to a HTTP response
func WriteJSON(w http.ResponseWriter, response interface{}) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		WriteErrorMessage(w, http.StatusInternalServerError, "Failed to process request")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// WriteError writes a HttpError value to a HTTP response
func WriteError(w http.ResponseWriter, err error) {
	if httpError, ok := err.(*HttpError); ok {
		WriteErrorMessage(w, httpError.StatusCode, httpError.Message)
	} else {
		WriteErrorMessage(w, http.StatusInternalServerError, err.Error())
	}
}

// WriteErrorMessage writes an error message to a HTTP response
func WriteErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	LogErrorMessage(message)

	msg := ErrorMessage{message}
	data, err := json.Marshal(msg)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

// RollbackOrLog tries to rollback a SQL transaction, and logs if it fails
func RollbackOrLog(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		LogError(err)
	}
}
