package main

import "main/utils"

type (
	// NewDoseHistoryEntry represents a to-be inserted dose history entry
	NewDoseHistoryEntry struct {
		DoseID        int    `json:"doseId"`
		DispensedDay  string `json:"dispensedDay"`
		DispensedTime string `json:"dispensedTime"`
	}

	// DoseHistoryEntrySummary contains basic information on a dose history entry
	DoseHistoryEntrySummary struct {
		ID            int                 `json:"id"`
		DispensedDay  string              `json:"dispensedDay"`
		DispensedTime string              `json:"dispensedTime"`
		Dose          utils.MinimalEntity `json:"dose"`
	}

	// DoseHistoryEntryDetails contains basic information on a dose history entry
	DoseHistoryEntryDetails struct {
		ID            int                 `json:"id"`
		DispensedDay  string              `json:"dispensedDay"`
		DispensedTime string              `json:"dispensedTime"`
		Dose          utils.MinimalEntity `json:"dose"`
	}
)

var (
	doseHistorySearchMapping SearchMapping
)

func init() {
	// Initialize field mappings
	doseHistorySearchMapping = NewMapping()
	doseHistorySearchMapping.DefineFieldMapping("dose", FieldMapping{
		SearchType: SearchTypeEqual,
		DBField:    "DH.DoseID",
	})

	doseHistorySearchMapping.DefineFieldMapping("dispensedday", FieldMapping{
		SearchType: SearchTypeEqual,
		DBField:    "DH.DispensedDay",
	})
}

func CreateDoseHistoryEntry(userID int, newDoseHistoryEntry NewDoseHistoryEntry) (DoseHistoryEntryDetails, error) {
	// Insert the new dose history in the database
	result, err := db.Exec(`INSERT INTO DoseHistory (DoseID, DispensedDay, DispensedTime)
	VALUES (?, ?, ?)`, newDoseHistoryEntry.DoseID, newDoseHistoryEntry.DispensedDay, newDoseHistoryEntry.DispensedTime)

	if err != nil {
		return DoseHistoryEntryDetails{}, utils.InternalServerError(err)
	}

	// Notify the dispatcher that the dose summaries have been updated
	summaries, err := ListDoseSummaries(userID)
	if err != nil {
		return DoseHistoryEntryDetails{}, err
	}

	doseSummariesSubject.DoseSummariesUpdated(userID, summaries)

	// Read the updated dose history entry and return
	doseHistoryEntryID, err := result.LastInsertId()

	if err != nil {
		return DoseHistoryEntryDetails{}, utils.InternalServerError(err)
	}

	return ReadDoseHistoryEntry(userID, int(doseHistoryEntryID))
}

// ListDoseHistoryEntries returns a list of all dose history entries for a given user and search query
func ListDoseHistoryEntries(userID int, search map[string]string) ([]DoseHistoryEntrySummary, error) {
	// Create the query using the search mapping
	query, queryParams := doseHistorySearchMapping.CreateQuery(`SELECT DH.ID, DH.DispensedDay, DH.DispensedTime, D.ID, D.Title FROM DoseHistory DH
	LEFT JOIN Doses D ON DH.DoseID = D.ID
	WHERE D.UserID = ? AND %MAPPING_CONDITIONS%`, search, userID)

	// Read all matching dose history entries from the database
	rows, err := db.Query(query, queryParams...)

	if err != nil {
		return []DoseHistoryEntrySummary{}, utils.InternalServerError(err)
	}

	// Iterate over all results and store in a slice
	doseHistoryEntries := []DoseHistoryEntrySummary{}
	var dhe DoseHistoryEntrySummary

	for rows.Next() {
		err = rows.Scan(&dhe.ID, &dhe.DispensedDay, &dhe.DispensedTime, &dhe.Dose.ID, &dhe.Dose.Title)
		if err != nil {
			return []DoseHistoryEntrySummary{}, utils.InternalServerError(err)
		}

		doseHistoryEntries = append(doseHistoryEntries, dhe)
	}

	return doseHistoryEntries, nil
}

// ReadDoseHistoryEntry returns a dose history entry for the given user and dose history entry ID
func ReadDoseHistoryEntry(userID, doseHistoryEntryID int) (DoseHistoryEntryDetails, error) {
	// Read from the database and return
	var dhe DoseHistoryEntryDetails

	err := db.QueryRow(`SELECT DH.ID, DH.DispensedDay, DH.DispensedTime, D.ID, D.Title FROM DoseHistory DH
	LEFT JOIN Doses D ON DH.DoseID = D.ID
	WHERE D.UserID = ? AND DH.ID = ?`, userID, doseHistoryEntryID).Scan(&dhe.ID, &dhe.DispensedDay, &dhe.DispensedTime, &dhe.Dose.ID, &dhe.Dose.Title)

	if err != nil {
		return dhe, utils.InternalServerError(err)
	}

	return dhe, nil
}
