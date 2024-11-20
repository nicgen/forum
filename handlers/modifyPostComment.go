package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

func ModifyPostComment(w http.ResponseWriter, r *http.Request) {
	// Getting the database data
	db := lib.GetDB()

	var state_modify string
	status := r.URL.Query().Get("delete")
	id := r.URL.Query().Get("id")
	content := r.URL.Query().Get("content")

	if status == "post" {
		state_modify = `UPDATE Posts SET Title = ? WHERE ID = ?`
	} else {
		state_modify = `UPDATE Comments SET Text = ? WHERE ID = ?`
	}
	_, err_delete := db.Exec(state_modify, content, id)
	if err_delete != nil && status == "post" {
		lib.ErrorServer(w, "Error modifying User's post")
	} else if err_delete != nil && status == "comment" {
		lib.ErrorServer(w, "Error modifying User's comment")
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
