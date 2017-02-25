package main

import "main/utils"

type (
	// DoseSummarySummary contains summary information on a dose summary
	DoseSummarySummary struct {
		Date           string `json:"date"`
		DispensedCount int    `json:"dispensedCount"`
		PendingCount   int    `json:"pendingCount"`
		TotalCount     int    `json:"totalCount"`
	}

	// DoseStatus contains status information on a dose for a given day
	DoseStatus struct {
		DispensedTime  string              `json:"dispensedTime"`
		Dispensed      bool                `json:"dispensed"`
		Pending        bool                `json:"pending"`
		BeingDispensed bool                `json:"beingDispensed"`
		Dose           utils.MinimalEntity `json:"dose"`
	}
)

// ListDoseSummaries returns a list of dose summaries for a given user ID
func ListDoseSummaries(userID int) ([]DoseSummarySummary, error) {
	// Read dose summaries from the database
	rows, err := db.Query(`SELECT H.DispensedDay AS HistoryDay, COUNT(H.DoseID) AS DispensedDayCount,
		(SELECT COUNT(*) FROM Doses D
			LEFT JOIN DoseHistory DH ON DH.DoseID = D.ID AND DH.DispensedDay = CURRENT_DATE()
			WHERE ISNULL(DH.ID) AND
				H.DispensedDay = CURRENT_DATE() AND
				((D.DispenseAfter < D.DispenseBefore AND CURRENT_TIME() < D.DispenseBefore) OR
					D.DispenseAfter > D.DispenseBefore)) AS PendingDayCount,
		(SELECT COUNT(*) FROM Doses
			WHERE UserID = ? AND DATE(CreatedOn) <= H.DispensedDay) AS TotalDayCount
	FROM (SELECT
		CASE
			WHEN (D.DispenseAfter > D.DispenseBefore AND DH.DispensedTime < D.DispenseAfter)
				THEN DATE_ADD(DH.DispensedDay, INTERVAL -1 DAY)
			ELSE DH.DispensedDay
		END AS DispensedDay,
		D.ID AS DoseID
	FROM Doses D
		RIGHT JOIN DoseHistory DH ON DH.DoseID = D.ID
	WHERE D.UserID = ?) H GROUP BY H.DispensedDay
	ORDER BY H.DispensedDay DESC`, userID, userID)

	if err != nil {
		return []DoseSummarySummary{}, utils.InternalServerError(err)
	}

	// Iterate over rows and fill slice of dose summaries
	summaries := []DoseSummarySummary{}
	var summary DoseSummarySummary

	for rows.Next() {
		err = rows.Scan(&summary.Date, &summary.DispensedCount, &summary.PendingCount, &summary.TotalCount)
		if err != nil {
			return summaries, utils.InternalServerError(err)
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// ReadDoseSummary reads the dose summary details for a given user ID and date
func ReadDoseSummary(userID int, date string) ([]DoseStatus, error) {
	// Read dose statuses from the database
	rows, err := db.Query(`SELECT D.ID AS DoseID, D.Title AS DoseTitle, IFNULL(DH.DispensedTime, '') AS DispensedTime,
  (NOT ISNULL(DH.ID)) AS Dispensed,
  (ISNULL(DH.ID) AND (
    (D.DispenseAfter < D.DispenseBefore AND CURRENT_TIME() < D.DispenseBefore) OR D.DispenseAfter > D.DispenseBefore
  ) AND ? = CURRENT_DATE()) AS Pending,
  (ISNULL(DH.ID) AND (
    (D.DispenseAfter < D.DispenseBefore AND CURRENT_TIME() BETWEEN D.DispenseBefore AND D.DispenseAfter) OR
    (D.DispenseAfter > D.DispenseBefore AND CURRENT_TIME() >= D.DispenseAfter)
  ) AND ? = CURRENT_DATE()) AS BeingDispensed
  FROM Doses D
  LEFT JOIN DoseHistory DH
    ON DH.DoseID = D.ID AND DH.DispensedDay = ? AND
       (D.DispenseAfter < D.DispenseBefore OR DH.DispensedTime >= D.DispenseAfter)
  WHERE
    DATE(D.CreatedOn) <= ? AND D.UserID = ?
	ORDER BY D.DispenseAfter`, date, date, date, date, userID)

	if err != nil {
		return []DoseStatus{}, utils.InternalServerError(err)
	}

	// Iterate over rows and store in slice
	statuses := []DoseStatus{}
	var status DoseStatus
	var dispensed, pending, beingDispensed byte

	for rows.Next() {
		err = rows.Scan(&status.Dose.ID, &status.Dose.Title, &status.DispensedTime, &dispensed, &pending, &beingDispensed)
		if err != nil {
			return statuses, utils.InternalServerError(err)
		}

		status.Dispensed = dispensed == 1
		status.Pending = pending == 1
		status.BeingDispensed = beingDispensed == 1

		statuses = append(statuses, status)
	}

	return statuses, nil
}
