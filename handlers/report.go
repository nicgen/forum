package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"net/http"
)

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	// recupérer ID post
	db := lib.GetDB()

	formValues := r.URL.Query()
	cookie_user_uuid, err_useruuid := r.Cookie("session_id")
	cookie_username, err_username := r.Cookie("username")

	uuid := cookie_user_uuid.Value
	username := cookie_username.Value
	report_post := formValues.Get("Report_post")
	report_title := formValues.Get("Report_title")

	var testID string
	state_check_report := `SELECT ID FROM Report WHERE User_UUID = ? AND Post_ID = ?`
	err1_db := db.QueryRow(state_check_report, uuid, report_post).Scan(&testID)

	if err_username != nil {
		lib.ErrorServer(w, "Error getting Username from the cookies")
	} else if err_useruuid != nil {
		lib.ErrorServer(w, "Error getting Creation Date from the cookies")
	}
	if err1_db == sql.ErrNoRows {
		var respons_report string

		state_report := `INSERT INTO Report (User_UUID, Username, Post_ID,Title, Respons_Text) VALUES (?, ?, ?, ?, ?)`
		_, err_db := db.Exec(state_report, uuid, username, report_post, report_title, respons_report)
		if err_db != nil {
			lib.ErrorServer(w, "Error report")
		}
	} else if len(testID) != 0 {
		lib.ErrorServer(w, "Report already completed")
	} else {
		lib.ErrorServer(w, "Error checking report")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
