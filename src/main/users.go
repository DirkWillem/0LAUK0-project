package main

import (
	"database/sql"
	"fmt"
)

type (
	UserSummary struct {
		ID           int    `json:"id"`
		Username     string `json:"username"`
		FullName     string `json:"fullName"`
		Role         string `json:"role"`
		Email        string `json:"email"`
		EmailMD5     string `json:"emailMD5"`
	}

	UserDetails struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		FullName string `json:"fullName"`
		Role     string `json:"role"`
		Email    string `json:"email"`
		EmailMD5 string `json:"emailMD5"`

		Patients    []UserSummary `json:"patients,omitempty"`
		Customers   []UserSummary `json:"customers,omitempty"`
		Doctors     []UserSummary `json:"doctor,omitempty"`
		Pharmacists []UserSummary `json:"pharmacist,omitempty"`
	}
)

// ListUsers returns a list of all users
func ListUsers() ([]UserSummary, error) {
	// Read all users from the database
	rows, err := db.Query(`SELECT ID, Username, FullName, Role, Email FROM Users`)

	if err != nil {
		return []UserSummary{}, InternalServerError(err)
	}

	// Iterate over all rows and store in slice
	return readUsersFromRows(rows)
}

// ReadUser returns a user by its ID
func ReadUser(userID int) (UserDetails, error) {
	// Read user from the database
	var user UserDetails

	err := db.QueryRow(`SELECT ID, Username, FullName, Role, Email FROM Users
	WHERE ID = ?`, userID).Scan(&user.ID, &user.Username, &user.FullName, &user.Role, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return UserDetails{}, NotFoundErrorMessage(fmt.Sprintf("No user with ID %d found", userID))
		}
		return UserDetails{}, InternalServerError(err)
	}

	user.EmailMD5 = HashMD5(user.Email)

	// Retrieve relations
	switch user.Role {
	case PatientRole:
		user.Doctors, err = ListRelations(user.ID, DoctorRole)
		if err != nil {
			return user, err
		}
		user.Pharmacists, err = ListRelations(user.ID, PharmacistRole)
		if err != nil {
			return user, err
		}
	case DoctorRole:
		user.Patients, err = ListRelatedPatients(user.ID)
		if err != nil {
			return user, err
		}
	case PharmacistRole:
		user.Customers, err = ListRelatedPatients(user.ID)
		if err != nil {
			return user, err
		}
	}

	// Return
	return user, nil
}

// ListRelatedPatients returns a list of all patients related to a user
func ListRelatedPatients(userID int) ([]UserSummary, error) {
	// Read patients from database
	rows, err := db.Query(`SELECT U.ID, U.Username, U.FullName, U.Role, U.Email FROM PatientRelations PR
	LEFT JOIN Users U ON PR.PatientID = U.ID
	WHERE PR.RelationID = ?`, userID)

	if err != nil {
		return []UserSummary{}, InternalServerError(err)
	}

	return readUsersFromRows(rows)
}

// ListRelations returns a list of all relations of a patient
func ListRelations(userID int, role string) ([]UserSummary, error) {
	// Read patients from database
	rows, err := db.Query(`SELECT U.ID, U.Username, U.FullName, U.Role, U.Email FROM PatientRelations PR
	LEFT JOIN Users U ON PR.RelationID = U.ID
	WHERE PR.PatientID = ? AND U.Role = ?`, userID, role)

	if err != nil {
		return []UserSummary{}, InternalServerError(err)
	}

	return readUsersFromRows(rows)
}

// readUsersFromRows is a helper function to easily read a *sql.Rows of users into a slice
func readUsersFromRows(rows *sql.Rows) ([]UserSummary, error) {
	var users []UserSummary
	var user UserSummary

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Role, &user.Email)
		if err != nil {
			return users, InternalServerError(err)
		}

		user.EmailMD5 = HashMD5(user.Email)
		users = append(users, user)
	}

	return users, nil
}
