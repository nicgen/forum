package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"net/http"
)

func Report_Delete(w http.ResponseWriter, r *http.Request) {
	// recupérer ID post
	db := lib.GetDB()

	formValues := r.URL.Query()
	report_postID := formValues.Get("report_postID")
	report_ID := formValues.Get("report_ID")
	respons_text := formValues.Get("respons_text")
	name1 := formValues.Get("1")

	if len(name1) != 0 {
		state_reaction := `SELECT ID From Report WHERE ID = ? AND Post_Id = ?`
		err_reaction := db.QueryRow(state_reaction, report_ID, report_postID).Scan(&state_reaction)

		if err_reaction == sql.ErrNoRows {
			lib.ErrorServer(w, "Error Report exist")
			return
		}
		var state_delete string
		state_delete = `DELETE FROM Posts WHERE ID = ?`
		_, err_db := db.Exec(state_delete, report_postID)
		if err_db != nil {
			lib.ErrorServer(w, "Error delete_report")
		}

		state_report := `UPDATE Report SET Respons_Text = ? WHERE ID = ?`
		_, err_db = db.Exec(state_report, respons_text, report_ID)
		if err_db != nil {
			lib.ErrorServer(w, "Error updating respons_text")
		}
	} else {
		state_report := `UPDATE Report SET Respons_Text = ? WHERE ID = ?`
		_, err_db := db.Exec(state_report, respons_text, report_ID)
		if err_db != nil {
			lib.ErrorServer(w, "Error updating respons_text")
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
