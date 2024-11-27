package handlers

import (
	"forum/cmd/lib"
	"forum/models"
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
		//Erreur critique : Erreur deleting user post
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting user's post. Please try again.",
		}
		HandleError(w, err.StatusCode, err.Message)
	} else if err_delete != nil && status == "comment" {
		//Erreur critique : Erreur deleting user comment
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting user's comment",
		}
		HandleError(w, err.StatusCode, err.Message)
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
