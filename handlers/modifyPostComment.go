package handlers

import (
	"fmt"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
)

func ModifyPostComment(w http.ResponseWriter, r *http.Request) {
	// Getting the database data
	db := lib.GetDB()

	var state_modify string
	status := r.URL.Query().Get("delete")
	modify_type := r.URL.Query().Get("type")
	id := r.URL.Query().Get("id")
	content := r.URL.Query().Get("content")
	fmt.Println(content)

	if status == "post" && modify_type == "title" {
		state_modify = `UPDATE Posts SET Title = ? WHERE ID = ?`
	} else if status == "post" && modify_type == "content" {
		state_modify = `UPDATE Posts SET Text = ? WHERE ID = ?`
	} else {
		state_modify = `UPDATE Comments SET Text = ? WHERE ID = ?`
	}
	_, err_delete := db.Exec(state_modify, content, id)
	if err_delete != nil && status == "post" && len(content) != 0 {

		//Erreur critique : échec de la modification du commentaire
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error modifying User's comment. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
	} else if err_delete != nil && status == "comment" && len(content) != 0 {

		//Erreur critique : echec de la modification du commentaire
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error modifying user's comment. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
	} else {
		data := lib.DataTest(w, r)
		data = lib.ErrorMessage(w, data, "ContentEmpty")

		lib.RenderTemplate(w, "layout/index", "page/index", data)
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
