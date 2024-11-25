package handlers

import (
	"fmt"
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

	fmt.Println(report_ID)
	fmt.Println(report_postID)
	fmt.Println(respons_text)

	var state_delete string
	state_delete = `DELETE FROM Posts WHERE ID = ?`

	_, err_db := db.Exec(state_delete, report_postID)
	if err_db != nil {
		lib.ErrorServer(w, "Error delete_report")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
