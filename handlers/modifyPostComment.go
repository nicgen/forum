package handlers

import (
	"forum/cmd/lib"
	"forum/models"
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

		//Erreur critique : échec de la modification du commentaire
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error modifying User's comment. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
	} else if err_delete != nil && status == "comment" {

		//Erreur critique : echec de la modification du commentaire
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error modifying user's comment. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
