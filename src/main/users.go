package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/hlandau/passlib.v1"
	"main/utils"
)

type (
	NewUser struct {
		Username  string `json:"username"`
		FullName  string `json:"fullName"`
		Role      string `json:"role"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		Birthdate string `json:"birthdate"`
		Gender    string `json:"gender"`
		Phone     string `json:"phone"`

		PatientIDs    []int `json:"patientIds"`
		CustomerIDs   []int `json:"customerIds"`
		DoctorIDs     []int `json:"doctorIds"`
		PharmacistIDs []int `json:"pharmacistIds"`
	}

	// UserSummary contains basic information on a user
	UserSummary struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		FullName string `json:"fullName"`
		Role     string `json:"role"`
		Email    string `json:"email"`
		EmailMD5 string `json:"emailMD5"`
		Phone     string `json:"phone"`
	}

	// UserDetails contains all information on a user
	UserDetails struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		FullName string `json:"fullName"`
		Role     string `json:"role"`
		Email    string `json:"email"`
		EmailMD5 string `json:"emailMD5"`
		Birthdate string `json:"birthdate"`
		Gender    string `json:"gender"`
		Phone     string `json:"phone"`

		Patients    []UserSummary `json:"patients,omitempty"`
		Customers   []UserSummary `json:"customers,omitempty"`
		Doctors     []UserSummary `json:"doctors,omitempty"`
		Pharmacists []UserSummary `json:"pharmacists,omitempty"`
	}

	// UpdatedUser represents an updated user
	UpdatedUser struct {
		Username string `json:"username"`
		FullName string `json:"fullName"`
		Email    string `json:"email"`

		Patients []struct {
			ID int `json:"id"`
		} `json:"patients"`

		Customers []struct {
			ID int `json:"id"`
		} `json:"customers"`

		Doctors []struct {
			ID int `json:"id"`
		} `json:"doctors"`

		Pharmacists []struct {
			ID int `json:"id"`
		} `json:"pharmacists"`
	}
)

var (
	usersSearchMapping SearchMapping
)

func init() {
	usersSearchMapping = NewMapping()
	usersSearchMapping.DefineFieldMapping("role", FieldMapping{
		SearchType: SearchTypeEqual,
		DBField:    "Role",
	})
}

// CreateUser creates a new user
func CreateUser(newUser NewUser) (UserDetails, error) {
	// Begin SQL transaction
	tx, err := db.Begin()

	if err != nil {
		utils.RollbackOrLog(tx)
		return UserDetails{}, utils.InternalServerError(err)
	}

	// Create the user record
	passHash, err := passlib.Hash(newUser.Password)

	if err != nil {
		utils.RollbackOrLog(tx)
		return UserDetails{}, utils.InternalServerError(err)
	}

	result, err := tx.Exec(`INSERT INTO Users (Username, FullName, PasswordHash, Role, Email)
	VALUES ($1, $2, $3, $4, $5)`, newUser.Username, newUser.FullName, passHash, newUser.Role, newUser.Email)

	if err != nil {
		utils.RollbackOrLog(tx)
		return UserDetails{}, utils.InternalServerError(err)
	}

	userID, err := result.LastInsertId()

	if err != nil {
		utils.RollbackOrLog(tx)
		return UserDetails{}, utils.InternalServerError(err)
	}

	// Determine to-be inserted relation IDs
	insertedRelationIDs := []int{}
	insertedPatientIDs := []int{}

	switch newUser.Role {
	case PatientRole:
		insertedRelationIDs = append(newUser.DoctorIDs, newUser.PharmacistIDs...)
	case DoctorRole:
		insertedPatientIDs = newUser.PatientIDs
	case PharmacistRole:
		insertedPatientIDs = newUser.CustomerIDs
	}

	// Insert relations into the database
	for _, insertedRelationID := range insertedRelationIDs {
		_, err = db.Exec(`INSERT INTO PatientRelations (PatientID, RelationID) VALUES ($1, $2)`, userID, insertedRelationID)
		if err != nil {
			utils.RollbackOrLog(tx)
			return UserDetails{}, utils.InternalServerError(err)
		}
	}

	for _, insertedPatientID := range insertedPatientIDs {
		_, err = db.Exec(`INSERT INTO PatientRelations (PatientID, RelationID) VALUES ($1, $2)`, insertedPatientID, userID)
		if err != nil {
			utils.RollbackOrLog(tx)
			return UserDetails{}, utils.InternalServerError(err)
		}
	}

	// Commit SQL transaction
	err = tx.Commit()

	if err != nil {
		utils.RollbackOrLog(tx)
		return UserDetails{}, utils.InternalServerError(err)
	}

	return ReadUser(int(userID))
}

// ListUsers returns a list of all users
func ListUsers(search map[string]string) ([]UserSummary, error) {
	query, queryParams := usersSearchMapping.CreateQuery(`SELECT id, username, fullname, role, email, phone FROM users
	WHERE %MAPPING_CONDITIONS%`, search)

	// Read all users from the database
	rows, err := db.Query(query, queryParams...)

	if err != nil {
		return []UserSummary{}, utils.InternalServerError(err)
	}

	// Iterate over all rows and store in slice
	return readUsersFromRows(rows)
}

// ReadUser returns a user by its ID
func ReadUser(userID int) (UserDetails, error) {
	// Read user from the database
	var user UserDetails

	err := db.QueryRow(`SELECT ID, Username, FullName, Role, Email, birthdate, gender, phone FROM Users
	WHERE ID = $1`, userID).Scan(&user.ID, &user.Username, &user.FullName, &user.Role, &user.Email, &user.Birthdate, &user.Gender, &user.Phone)

	if err != nil {
		if err == sql.ErrNoRows {
			return UserDetails{}, utils.NotFoundErrorMessage(fmt.Sprintf("No user with ID %d found", userID))
		}
		return UserDetails{}, utils.InternalServerError(err)
	}

	user.EmailMD5 = utils.HashMD5(user.Email)

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

// UpdateUser updates a user
func UpdateUser(userID int, updatedUser UpdatedUser) (UserDetails, error) {
	user, err := ReadUser(userID)
	if err != nil {
		return UserDetails{}, err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		utils.RollbackOrLog(tx)
		return UserDetails{}, err
	}

	// Update user record
	_, err = db.Exec(`UPDATE Users
	SET
		Username = $1,
		FullName = $2,
		Email = $3
	WHERE ID = $4`, updatedUser.Username, updatedUser.FullName, updatedUser.Email, userID)

	if err != nil {
		utils.RollbackOrLog(tx)
		return UserDetails{}, err
	}

	// Update relations
	if user.Role == PatientRole {
		processedIds := []int{}

		// Insert all new and keep track of which relations didn't change
		for _, updatedRelation := range append(updatedUser.Doctors, updatedUser.Pharmacists...) {
			isNew := true

			for _, relation := range append(user.Doctors, user.Pharmacists...) {
				if relation.ID == updatedRelation.ID {
					processedIds = append(processedIds, updatedRelation.ID)
					isNew = false
				}
			}

			if isNew {
				_, err := tx.Exec(`INSERT INTO PatientRelations (PatientID, RelationID) VALUES ($1, $2)`, userID, updatedRelation.ID)
				if err != nil {
					utils.RollbackOrLog(tx)
					return UserDetails{}, err
				}
			}
		}

		// Remove all old relations
		for _, relation := range append(user.Doctors, user.Pharmacists...) {
			processed := false

			for _, processedID := range processedIds {
				if processedID == relation.ID {
					processed = true
				}
			}

			if !processed {
				_, err := tx.Exec(`DELETE FROM PatientRelations WHERE PatientID = $1 AND RelationID = $2`, userID, relation.ID)
				if err != nil {
					utils.RollbackOrLog(tx)
					return UserDetails{}, err
				}
			}
		}
	} else if user.Role == DoctorRole || user.Role == PharmacistRole {
		processedIds := []int{}

		// Insert all new and keep track of which relations didn't change
		for _, updatedPatient := range append(updatedUser.Patients, updatedUser.Customers...) {
			isNew := true

			for _, patient := range append(user.Patients, user.Customers...) {
				if patient.ID == updatedPatient.ID {
					processedIds = append(processedIds, updatedPatient.ID)
					isNew = false
				}
			}

			if isNew {
				_, err := tx.Exec(`INSERT INTO PatientRelations (PatientID, RelationID) VALUES ($1, $2)`, updatedPatient.ID, userID)
				if err != nil {
					utils.RollbackOrLog(tx)
					return UserDetails{}, err
				}
			}
		}

		// Remove all old relations
		for _, patient := range append(user.Patients, user.Customers...) {
			processed := false

			for _, processedID := range processedIds {
				if processedID == patient.ID {
					processed = true
				}
			}

			if !processed {
				_, err := tx.Exec(`DELETE FROM PatientRelations WHERE PatientID = $1 AND RelationID = $2`, patient.ID, userID)
				if err != nil {
					utils.RollbackOrLog(tx)
					return UserDetails{}, err
				}
			}
		}
	}

	// Commit the transaction and return
	err = tx.Commit()

	if err != nil {
		utils.RollbackOrLog(tx)
		return UserDetails{}, utils.InternalServerError(err)
	}

	return ReadUser(userID)
}

// DeleteUser deletes a user
func DeleteUser(userID int) error {
	_, err := db.Exec(`DELETE FROM Users WHERE ID = $1`, userID)

	if err != nil {
		return utils.InternalServerError(err)
	}

	return nil
}

// ListRelatedPatients returns a list of all patients related to a user
func ListRelatedPatients(userID int) ([]UserSummary, error) {
	// Read patients from database
	rows, err := db.Query(`SELECT U.ID, U.Username, U.FullName, U.Role, U.Email, U.phone FROM PatientRelations PR
	LEFT JOIN Users U ON PR.PatientID = U.ID
	WHERE PR.RelationID = $1`, userID)

	if err != nil {
		return []UserSummary{}, utils.InternalServerError(err)
	}

	return readUsersFromRows(rows)
}

// ListRelations returns a list of all relations of a patient
func ListRelations(userID int, role string) ([]UserSummary, error) {
	// Read patients from database
	rows, err := db.Query(`SELECT U.ID, U.Username, U.FullName, U.Role, U.Email, U.phone FROM PatientRelations PR
	LEFT JOIN Users U ON PR.RelationID = U.ID
	WHERE PR.PatientID = $1 AND U.Role = $2`, userID, role)

	if err != nil {
		return []UserSummary{}, utils.InternalServerError(err)
	}

	return readUsersFromRows(rows)
}

// readUsersFromRows is a helper function to easily read a *sql.Rows of users into a slice
func readUsersFromRows(rows *sql.Rows) ([]UserSummary, error) {
	users := []UserSummary{}
	var user UserSummary

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Role, &user.Email, &user.Phone)
		if err != nil {
			return users, utils.InternalServerError(err)
		}

		user.EmailMD5 = utils.HashMD5(user.Email)
		users = append(users, user)
	}

	return users, nil
}
