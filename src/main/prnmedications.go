package main

import (
	"database/sql"
	"fmt"
	"main/utils"
)

type (
	// PRNMedicationSummary contains a summary of data on a PRN (Pro Re Nata) medication
	PRNMedicationSummary struct {
		ID          int               `json:"id"`
		Description string            `json:"description"`
		UserID      int               `json:"userId"`
		MaxDaily    int               `json:"maxDaily"`
		MinInterval int               `json:"minInterval"`
		Medication  MedicationSummary `json:"medication"`
	}

	// PRNMedicationDetails contains all data on a PRN (Pro Re Nata) medication
	PRNMedicationDetails struct {
		ID          int               `json:"id"`
		Description string            `json:"description"`
		UserID      int               `json:"userId"`
		MaxDaily    int               `json:"maxDaily"`
		MinInterval int               `json:"minInterval"`
		Medication  MedicationSummary `json:"medication"`
	}

	// NewPRNMedication contains data of a new PRN medication
	NewPRNMedication struct {
		Description  string `json:"description"`
		MaxDaily     int    `json:"maxDaily"`
		MinInterval  int    `json:"minInterval"`
		MedicationID int    `json:"medication"`
	}

	// UpdatedPRNMedication contains data of a to-be updated PRN medication
	UpdatedPRNMedication struct {
		Description  string `json:"description"`
		MaxDaily     int    `json:"maxDaily"`
		MinInterval  int    `json:"minInterval"`
		MedicationID int    `json:"medication"`
	}
)

// ToSummary transforms a PRNMedicationDetails into its summary counterpart
func (m PRNMedicationDetails) ToSummary() PRNMedicationSummary {
	return PRNMedicationSummary{
		ID:          m.ID,
		Description: m.Description,
		UserID:      m.UserID,
		MaxDaily:    m.MaxDaily,
		MinInterval: m.MinInterval,
		Medication:  m.Medication,
	}
}

// CreatePRNMedication creates a new PRN medication
func CreatePRNMedication(userID int, newMedication NewPRNMedication) (PRNMedicationDetails, error) {
	// Insert the medication into the database
	var medicationID int
	err := db.QueryRow(`INSERT INTO prnmedications (description, maxdaily, mininterval, medicationid)
	VALUES ($1, $2, $3, $4) RETURNING id`, newMedication.Description, newMedication.MaxDaily, newMedication.MinInterval, newMedication.MedicationID).Scan(&medicationID)

	if err != nil {
		return PRNMedicationDetails{}, utils.InternalServerError(err)
	}

	// Notify the dispatcher
	medication, err := ReadPRNMedication(userID, medicationID)

	if err != nil {
		return medication, utils.InternalServerError(err)
	}

	prnSubject.PRNMedicationAdded(userID, medication.ToSummary())

	return medication, nil
}

// ListPRNMedications returns a list of all PRN medications for a given user
func ListPRNMedications(userID int) ([]PRNMedicationSummary, error) {
	// Read the PRN medications from the database
	rows, err := db.Query(`SELECT p.id, p.description, p.userid, p.maxdaily, p.mininterval, m.id, m.title FROM prnmedications p
	LEFT JOIN medications m on p.medicationid = m.id
	WHERE p.userid = $1`, userID)

	if err != nil {
		return []PRNMedicationSummary{}, utils.InternalServerError(err)
	}

	// Iterate over all rows and store in a slice
	medications := []PRNMedicationSummary{}
	var m PRNMedicationSummary

	for rows.Next() {
		err = rows.Scan(&m.ID, &m.Description, &m.UserID, &m.MaxDaily, &m.MinInterval, &m.Medication.ID, &m.Medication.Title)
		if err != nil {
			return []PRNMedicationSummary{}, utils.InternalServerError(err)
		}

		medications = append(medications, m)
	}

	// Return list
	return medications, nil
}

// ReadPRNMedication returns a PRN medication for a user by its ID
func ReadPRNMedication(userID, prnMedicationID int) (PRNMedicationDetails, error) {
	// Read the PRN medication from the database
	var m PRNMedicationDetails

	err := db.QueryRow(`SELECT p.id, p.description, p.userid, p.maxdaily, p.mininterval, m.id, m.title FROM prnmedications p
	LEFT JOIN medications m on p.medicationid = m.id
	WHERE p.userid = $1`, userID, prnMedicationID).Scan(&m.ID, &m.Description, &m.UserID, &m.MaxDaily, &m.MinInterval, &m.Medication.ID, &m.Medication.Title)

	if err != nil {
		if err == sql.ErrNoRows {
			return PRNMedicationDetails{}, utils.NotFoundErrorMessage(fmt.Sprintf("No PRN medication with ID %d found", prnMedicationID))
		}
		return PRNMedicationDetails{}, utils.InternalServerError(err)
	}

	return m, nil
}

// UpdatePRNMedication updates an existing PRN medication
func UpdatePRNMedication(userID, prnMedicationID int, updatedMedication UpdatedPRNMedication) (PRNMedicationDetails, error) {
	// Update the medication
	_, err := db.Exec(`UPDATE prnmedications
	SET
		description = $1,
		maxdaily = $2,
		mininterval = $3,
		medicationid = $4
	WHERE id = $5 AND userid = $6`, updatedMedication.Description, updatedMedication.MaxDaily, updatedMedication.MinInterval, updatedMedication.MedicationID, prnMedicationID, userID)

	if err != nil {
		return PRNMedicationDetails{}, utils.InternalServerError(err)
	}

	// Notify the dispatcher
	medication, err := ReadPRNMedication(userID, prnMedicationID)
	if err != nil {
		return PRNMedicationDetails{}, utils.InternalServerError(err)
	}

	prnSubject.PRNMedicationUpdated(userID, medication.ToSummary())

	return medication, nil
}

// DeletePRNMedication deletes a PRN medication
func DeletePRNMedication(userID, prnMedicationID int) error {
	// Delete the entity in the database
	_, err := db.Exec(`DELETE FROM prnmedications WHERE id = $1 AND userid = $2`, prnMedicationID, userID)

	if err != nil {
		return utils.InternalServerError(err)
	}

	// Notify the dispatcher and return
	prnSubject.PRNMedicationDeleted(userID, prnMedicationID)
	return nil
}
