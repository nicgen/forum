package lib

import (
	"forum/models"
	"net/http"
)

func GetReport(w http.ResponseWriter, data map[string]interface{}, r *http.Request) map[string]interface{} {
	// db := lib.GetDB()

	state_report := `SELECT ID, User_UUID, Username,Post_ID,Title,Respons_Text FROM Report`
	// Users posts Request
	var reports []*models.Reports
	var err_report error
	rows_report, err_report := db.Query(state_report)

	if err_report != nil {
		ErrorServer(w, "Error accessing user reports")
	}

	defer rows_report.Close()

	for rows_report.Next() {
		var report models.Reports
		if err := rows_report.Scan(&report.ID, &report.User_UUID, &report.Username, &report.Post_ID, &report.Title, &report.Response_Text); err != nil {
			ErrorServer(w, "Error scanning report")
		}

		reports = append(reports, &report)
	}

	if err := rows_report.Err(); err != nil {
		ErrorServer(w, "Error iterating over user report")
	}

	data["Report"] = reports
	return data
}
