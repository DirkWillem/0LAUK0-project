package main

import (
	"database/sql"
	"fmt"
)

type (
	// MedicationSummary contains basic information on a single medication
	MedicationSummary struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	// MedicationDetails contains all information on a medication
	MedicationDetails struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	NewMedication struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
)

// CreateMedication creates a new medication
func CreateMedication(newMedication NewMedication) (MedicationDetails, error) {
  // Insert the medication into the database
  result, err := db.Exec(`INSERT INTO Medications (Title, Description)
  VALUES (?, ?)`, newMedication.Title, newMedication.Description)

  if err != nil {
    return MedicationDetails{}, InternalServerError(err)
  }

  // Read and return the created medication
  medicationID, err := result.LastInsertId()

  if err != nil {
    return MedicationDetails{}, InternalServerError(err)
  }

  return ReadMedication(int(medicationID))
}

// ListMedications returns a list of all medications
func ListMedications() ([]MedicationSummary, error) {
	// Read medications from the database
	rows, err := db.Query(`SELECT ID, Title, Description FROM Medications`)

	if err != nil {
		return []MedicationSummary{}, InternalServerError(err)
	}

	// Iterate over all rows and store in slice
	var medications []MedicationSummary
	var medication MedicationSummary

	for rows.Next() {
		err = rows.Scan(&medication.ID, &medication.Title, &medication.Description)
		if err != nil {
			return []MedicationSummary{}, InternalServerError(err)
		}

		medications = append(medications, medication)
	}

	// Return list
	return medications, nil
}

// ReadMedication returns a single medication
func ReadMedication(id int) (MedicationDetails, error) {
	// Read medication from the database and return
	var medication MedicationDetails

	err := db.QueryRow(`SELECT ID, Title, Description FROM Medications
  WHERE ID = ?`, id).Scan(&medication.ID, &medication.Title, &medication.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return MedicationDetails{}, NotFoundErrorMessage(fmt.Sprintf("No medication with ID %d found", id))
		}
		return MedicationDetails{}, InternalServerError(err)
	}

	return medication, err
}
