package main

import (
	"database/sql"
	"fmt"
	"main/utils"
)

type (
	// DoseSummary contains basic information on a dose
	DoseSummary struct {
		ID             int    `json:"id"`
		Title          string `json:"title"`
		DispenseAfter  string `json:"dispenseAfter"`
		DispenseBefore string `json:"dispenseBefore"`
		Description    string `json:"description"`
	}

	// DoseMedication contains information on a medication in a dose
	DoseMedication struct {
		Amount     int               `json:"amount"`
		Medication MedicationSummary `json:"medication"`
	}

	// DoseDetails contains all information on a dose
	DoseDetails struct {
		ID             int              `json:"id"`
		Title          string           `json:"title"`
		DispenseAfter  string           `json:"dispenseAfter"`
		DispenseBefore string           `json:"dispenseBefore"`
		Description    string           `json:"description"`
		Medications    []DoseMedication `json:"medications"`
	}

	// NewDose contains all information on a to-be inserted dose
	NewDose struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		DispenseAfter  string `json:"dispenseAfter"`
		DispenseBefore string `json:"dispenseBefore"`
		Medications    []struct {
			MedicationID int `json:"medicationId"`
			Amount       int `json:"amount"`
		}
	}

	// UpdatedDose contains all information on a to-be updated dose
	UpdatedDose struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		DispenseAfter  string `json:"dispenseAfter"`
		DispenseBefore string `json:"dispenseBefore"`
		Medications    []struct {
			Amount     int `json:"amount"`
			Medication struct {
				ID int `json:"id"`
			} `json:"medication"`
		}
	}
)

// ToSummary transforms a DoseDetails into its DoseSummary counterpart
func (dd DoseDetails) ToSummary() DoseSummary {
	return DoseSummary{
		ID:             dd.ID,
		Title:          dd.Title,
		DispenseAfter:  dd.DispenseAfter,
		DispenseBefore: dd.DispenseBefore,
		Description:    dd.Description,
	}
}

// CreateDose creates a new dose
func CreateDose(userID int, newDose NewDose) (DoseDetails, error) {
	// Begin a SQL transaction
	tx, err := db.Begin()
	if err != nil {
		return DoseDetails{}, utils.InternalServerError(err)
	}

	// Insert the dose into the Doses table
	result, err := tx.Exec(`INSERT INTO Doses (Title, Description, UserID, DispenseAfter, DispenseBefore)
  VALUES (?, ?, ?, ?, ?)`, newDose.Title, newDose.Description, userID, newDose.DispenseAfter, newDose.DispenseBefore)

	if err != nil {
		utils.RollbackOrLog(tx)
		return DoseDetails{}, utils.InternalServerError(err)
	}

	doseID, err := result.LastInsertId()

	if err != nil {
		utils.RollbackOrLog(tx)
		return DoseDetails{}, utils.InternalServerError(err)
	}

	// Insert the dose medications
	for _, medication := range newDose.Medications {
		_, err = tx.Exec(`INSERT INTO DoseMedications (DoseID, MedicationID, Amount)
    VALUES (?, ?, ?)`, doseID, medication.MedicationID, medication.Amount)

		if err != nil {
			utils.RollbackOrLog(tx)
			return DoseDetails{}, utils.InternalServerError(err)
		}
	}

	// Commit SQL transaction
	err = tx.Commit()

	if err != nil {
		utils.RollbackOrLog(tx)
		return DoseDetails{}, utils.InternalServerError(err)
	}

	// Notify the dispatcher and return
	dose, err := ReadDose(userID, int(doseID))

	if err != nil {
		return dose, err
	}

	dosesSubject.DoseAdded(userID, dose.ToSummary())

	return dose, err
}

// ListDoses returns a list of all doses for a user
func ListDoses(userID int) ([]DoseSummary, error) {
	// Read doses from the database
	rows, err := db.Query(`SELECT ID, Title, DispenseAfter, DispenseBefore, Description
  FROM Doses
  WHERE UserID = ?
  ORDER BY DispenseAfter`, userID)

	if err != nil {
		return []DoseSummary{}, utils.InternalServerError(err)
	}

	// Read doses into a slice
	doses := []DoseSummary{}
	var dose DoseSummary

	for rows.Next() {
		err := rows.Scan(&dose.ID, &dose.Title, &dose.DispenseAfter, &dose.DispenseBefore, &dose.Description)
		if err != nil {
			return doses, utils.InternalServerError(err)
		}

		doses = append(doses, dose)
	}

	return doses, nil
}

// ReadDose returns a dose for a given user and dose ID
func ReadDose(userID, doseID int) (DoseDetails, error) {
	// Read dose from the database
	var dose DoseDetails

	err := db.QueryRow(`SELECT ID, Title, DispenseAfter, DispenseBefore, Description
  FROM Doses
  WHERE ID = ? AND UserID = ?`, doseID, userID).Scan(&dose.ID, &dose.Title, &dose.DispenseAfter, &dose.DispenseBefore, &dose.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return dose, utils.NotFoundErrorMessage(fmt.Sprintf("No dose with ID '%d' for user with ID '%d' found.", doseID, userID))
		}
		return dose, utils.InternalServerError(err)
	}

	// Read the dose medications from the database
	rows, err := db.Query(`SELECT DM.Amount, M.ID, M.Title, M.Description FROM DoseMedications DM
  LEFT JOIN Medications M ON DM.MedicationID = M.ID
  WHERE DoseID = ?`, doseID)

	if err != nil {
		return dose, utils.InternalServerError(err)
	}

	// Read does medications into a slice
	dose.Medications = make([]DoseMedication, 0)
	var dm DoseMedication

	for rows.Next() {
		err = rows.Scan(&dm.Amount, &dm.Medication.ID, &dm.Medication.Title, &dm.Medication.Description)
		if err != nil {
			return dose, utils.InternalServerError(err)
		}

		dose.Medications = append(dose.Medications, dm)
	}

	return dose, nil
}

// UpdateDose updates a dose for a given user and dose ID
func UpdateDose(userID, doseID int, updatedDose UpdatedDose) (DoseDetails, error) {
	// Begin a SQL transaction
	tx, err := db.Begin()
	if err != nil {
		return DoseDetails{}, utils.InternalServerError(err)
	}

	// Get the current state of the dose
	dose, err := ReadDose(userID, doseID)
	if err != nil {
		utils.RollbackOrLog(tx)
		return DoseDetails{}, utils.InternalServerError(err)
	}

	// Update the dose
	_, err = tx.Exec(`UPDATE Doses
	SET
		Title = ?,
		Description = ?,
		DispenseAfter = ?,
		DispenseBefore = ?
	WHERE UserID = ? AND ID = ?`, updatedDose.Title, updatedDose.Description, updatedDose.DispenseAfter, updatedDose.DispenseBefore, userID, doseID)

	if err != nil {
		utils.RollbackOrLog(tx)
		return DoseDetails{}, utils.InternalServerError(err)
	}

	// Update or create dose medications
	processedDoseMedicationIDs := []int{}

	for _, updatedDoseMedication := range updatedDose.Medications {
		// Check whether the dose medication is new
		isNew := true

		// Iterate over all existing medications
		for _, doseMedication := range dose.Medications {
			if updatedDoseMedication.Medication.ID == doseMedication.Medication.ID {
				processedDoseMedicationIDs = append(processedDoseMedicationIDs, updatedDoseMedication.Medication.ID)
				isNew = false

				// If the amount changed, update it in the database
				if updatedDoseMedication.Amount != doseMedication.Amount {
					_, err = tx.Exec(`UPDATE DoseMedications
					SET
						Amount = ?
					WHERE DoseID = ? AND MedicationID = ?`, updatedDoseMedication.Amount, doseID, updatedDoseMedication.Medication.ID)

					if err != nil {
						utils.RollbackOrLog(tx)
						return DoseDetails{}, utils.InternalServerError(err)
					}
				}
			}
		}

		// If the dose medication was new, insert it in the database
		if isNew {
			_, err = tx.Exec(`INSERT INTO DoseMedications (DoseID, MedicationID, Amount)
			VALUES (?, ?, ?)`, doseID, updatedDoseMedication.Medication.ID, updatedDoseMedication.Amount)

			if err != nil {
				utils.RollbackOrLog(tx)
				return DoseDetails{}, utils.InternalServerError(err)
			}
		}
	}

	// Remove all dose medications that werent in the updated dose
	for _, doseMedication := range dose.Medications {
		processed := false

		for _, medicationID := range processedDoseMedicationIDs {
			if medicationID == doseMedication.Medication.ID {
				processed = true
				break
			}
		}

		if !processed {
			_, err = tx.Exec(`DELETE FROM DoseMedications
			WHERE DoseID = ? AND MedicationID = ?`, doseID, doseMedication.Medication.ID)

			if err != nil {
				utils.RollbackOrLog(tx)
				return DoseDetails{}, utils.InternalServerError(err)
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()

	if err != nil {
		utils.RollbackOrLog(tx)
		return DoseDetails{}, utils.InternalServerError(err)
	}

	// Notify the dispatcher and return
	dose, err = ReadDose(userID, int(doseID))

	if err != nil {
		return dose, err
	}

	dosesSubject.DoseUpdated(userID, dose.ToSummary())

	return dose, err
}

// DeleteDose deletes a dose for a given user and dose ID
func DeleteDose(userID, doseID int) error {
	// Delete the dose in the database
	_, err := db.Exec(`DELETE FROM Doses WHERE UserID = ? AND ID = ?`, userID, doseID)

	if err != nil {
		return utils.InternalServerError(err)
	}

	// Notify the dispatcher and return
	dosesSubject.DoseDeleted(userID, doseID)

	return nil
}
