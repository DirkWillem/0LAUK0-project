package main

import (
	"database/sql"
	"fmt"
)

type (
	// DoseSummary contains basic information on a dose
	DoseSummary struct {
		ID             int    `json:"id"`
		Title          string `json:"title"`
		DispenseAfter  string `json:"dispenseAfter"`
		DispenseBefore string `json:"dispenseBefore"`
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
)

// CreateDose creates a new dose
func CreateDose(userID int, newDose NewDose) (DoseDetails, error) {
  // Begin a database transaction
  tx, err := db.Begin()
  if err != nil {
    return DoseDetails{}, InternalServerError(err)
  }

  // Insert the dose into the Doses table
  result, err := tx.Exec(`INSERT INTO Doses (Title, Description, UserID, DispenseAfter, DispenseBefore)
  VALUES (?, ?, ?, ?, ?)`, newDose.Title, newDose.Description, userID, newDose.DispenseAfter, newDose.DispenseBefore)

  if err != nil {
    tx.Rollback()
    return DoseDetails{}, InternalServerError(err)
  }

  doseID, err := result.LastInsertId()

  if err != nil {
    tx.Rollback()
    return DoseDetails{}, InternalServerError(err)
  }

  // Insert the dose medications
  for _, medication := range newDose.Medications {
    _, err = tx.Exec(`INSERT INTO DoseMedications (DoseID, MedicationID, Amount)
    VALUES (?, ?, ?)`, doseID, medication.MedicationID, medication.Amount)

    if err != nil {
      tx.Rollback()
      return DoseDetails{}, InternalServerError(err)
    }
  }

  // Commit SQL transaction
  err = tx.Commit()

  if err != nil {
    tx.Rollback()
    return DoseDetails{}, InternalServerError(err)
  }

  return ReadDose(userID, int(doseID))
}

// ListDoses returns a list of all doses for a user
func ListDoses(userID int) ([]DoseSummary, error) {
	// Read doses from the database
	rows, err := db.Query(`SELECT ID, Title, DispenseAfter, DispenseBefore
  FROM Doses
  WHERE UserID = ?`, userID)

	if err != nil {
		return []DoseSummary{}, InternalServerError(err)
	}

	// Read doses into a slice
	var doses []DoseSummary
	var dose DoseSummary

	for rows.Next() {
		err := rows.Scan(&dose.ID, &dose.Title, &dose.DispenseAfter, &dose.DispenseBefore)
		if err != nil {
			return doses, InternalServerError(err)
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
			return dose, NotFoundErrorMessage(fmt.Sprintf("No dose with ID '%d' for user with ID '%d' found.", doseID, userID))
		}
		return dose, InternalServerError(err)
	}

	// Read the dose medications from the database
	rows, err := db.Query(`SELECT DM.Amount, M.ID, M.Title, M.Description FROM DoseMedications DM
  LEFT JOIN Medications M ON DM.MedicationID = M.ID
  WHERE DoseID = ?`, doseID)

	if err != nil {
		return dose, InternalServerError(err)
	}

	// Read does medications into a slice
	dose.Medications = make([]DoseMedication, 0)
	var dm DoseMedication

	for rows.Next() {
		err = rows.Scan(&dm.Amount, &dm.Medication.ID, &dm.Medication.Title, &dm.Medication.Description)
		if err != nil {
			return dose, InternalServerError(err)
		}

		dose.Medications = append(dose.Medications, dm)
	}

	return dose, nil
}
