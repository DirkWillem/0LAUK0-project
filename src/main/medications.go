package main

import (
	"database/sql"
	"fmt"
	"main/utils"
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

	// NewMedication contains all information on a to-be inserted medication
	NewMedication struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	// UpdatedMedication contains all information on a to-be updated medication
	UpdatedMedication struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
)

// ToSummary transforms a MedicationDetails into its MedicationSummary counterpart
func (md MedicationDetails) ToSummary() MedicationSummary {
	return MedicationSummary{
		ID:          md.ID,
		Title:       md.Title,
		Description: md.Description,
	}
}

// CreateMedication creates a new medication
func CreateMedication(newMedication NewMedication) (MedicationDetails, error) {
	// Insert the medication into the database
	result, err := db.Exec(`INSERT INTO Medications (Title, Description)
  VALUES ($1, $2)`, newMedication.Title, newMedication.Description)

	if err != nil {
		return MedicationDetails{}, utils.InternalServerError(err)
	}

	// Read and return the created medication
	medicationID, err := result.LastInsertId()

	if err != nil {
		return MedicationDetails{}, utils.InternalServerError(err)
	}

	// Notify the dispatcher a new entry has been inserted and return
	medication, err := ReadMedication(int(medicationID))

	if err != nil {
		return medication, err
	}

	medicationsSubject.EntityAdded(int(medicationID), medication)

	return medication, err
}

// ListMedications returns a list of all medications
func ListMedications() ([]MedicationSummary, error) {
	// Read medications from the database
	rows, err := db.Query(`SELECT ID, Title, Description FROM Medications`)

	if err != nil {
		return []MedicationSummary{}, utils.InternalServerError(err)
	}

	// Iterate over all rows and store in slice
	medications := []MedicationSummary{}
	var medication MedicationSummary

	for rows.Next() {
		err = rows.Scan(&medication.ID, &medication.Title, &medication.Description)
		if err != nil {
			return []MedicationSummary{}, utils.InternalServerError(err)
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
  WHERE ID = $1`, id).Scan(&medication.ID, &medication.Title, &medication.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return MedicationDetails{}, utils.NotFoundErrorMessage(fmt.Sprintf("No medication with ID %d found", id))
		}
		return MedicationDetails{}, utils.InternalServerError(err)
	}

	return medication, err
}

// UpdateMedication updates a medication with a given ID
func UpdateMedication(id int, updatedMedication UpdatedMedication) (MedicationDetails, error) {
	// Update the medication in the database
	_, err := db.Exec(`UPDATE Medications
	SET
		Title = $1,
		Description = $2
	WHERE ID = $3`, updatedMedication.Title, updatedMedication.Description, id)

	if err != nil {
		return MedicationDetails{}, utils.InternalServerError(err)
	}

	// Notify the dispatcher and return
	medication, err := ReadMedication(id)

	if err != nil {
		return medication, err
	}

	medicationsSubject.EntityUpdated(id, medication.ToSummary())

	return medication, err
}

// DeleteMedication deletes a medication with a given ID
func DeleteMedication(id int) error {
	// Delete the entity in the database
	_, err := db.Exec(`DELETE FROM Medications WHERE ID = $1`, id)

	if err != nil {
		return utils.InternalServerError(err)
	}

	// Notify the dispatcher and return
	medicationsSubject.EntityDeleted(id)
	return nil
}
