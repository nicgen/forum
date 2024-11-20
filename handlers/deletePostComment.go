package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

func DeletePostComment(w http.ResponseWriter, r *http.Request) {
	// Getting the database data
	db := lib.GetDB()

	var state_delete string
	status := r.URL.Query().Get("delete")
	id := r.URL.Query().Get("id")

	if status == "post" {
		state_delete = `DELETE FROM Posts WHERE ID = ?`
	} else {
		state_delete = `DELETE FROM Comments WHERE ID = ?`
	}
	_, err_delete := db.Exec(state_delete, id)
	if err_delete != nil && status == "post" {
		lib.ErrorServer(w, "Error deleting User's post")
	} else if err_delete != nil && status == "comment" {
		lib.ErrorServer(w, "Error deleting User's comment")
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
