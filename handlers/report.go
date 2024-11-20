package handlers

import (
	"fmt"
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

	if err_username != nil {
		lib.ErrorServer(w, "Error getting Username from the cookies")
	} else if err_useruuid != nil {
		lib.ErrorServer(w, "Error getting Creation Date from the cookies")
	}

	fmt.Println(report_post)
	fmt.Println(uuid)
	fmt.Println(username)

	var respons_report string

	state_report := `INSERT INTO Report (User_UUID, Username, Post_ID, Respons_Text) VALUES (?, ?, ?, ?)`
	_, err_db := db.Exec(state_report, uuid, username, report_post, respons_report)
	if err_db != nil {
		lib.ErrorServer(w, "Error report")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
