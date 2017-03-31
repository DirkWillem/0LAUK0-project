package main

import "main/utils"

type (
	// NewPRNHistoryEntry represents a to-be inserted PRN history entry
	NewPRNHistoryEntry struct {
		PRNMedicationID int    `json:"prnMedicationId"`
		DispensedDay    string `json:"dispensedDay"`
		DispensedTime   string `json:"dispensedTime"`
	}
)

// CreatePRNHistoryEntry creates a new PORN history entry
func CreatePRNHistoryEntry(userID int, newPRNHistoryEntry NewPRNHistoryEntry) error {
	// Insert the new PRN history entry
	var prnHistoryEntryID int
	err := db.QueryRow(`INSERT INTO prnhistory (prnmedicationid, dispensedday, dispensedtime)
	values ($1, $2, $3) RETURNING id`, newPRNHistoryEntry.PRNMedicationID, newPRNHistoryEntry.DispensedDay, newPRNHistoryEntry.DispensedTime).Scan(&prnHistoryEntryID)

	if err != nil {
		utils.InternalServerError(err)
	}

	return nil
}
