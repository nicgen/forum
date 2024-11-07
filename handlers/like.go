package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"net/http"
	"strings"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	cookie, _ := r.Cookie("session_id")

	// Getting the form values
	id := strings.TrimSpace(r.URL.Query().Get("id"))

	var reaction_status string
	state_reaction := `SELECT Status FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
	err_reaction := db.QueryRow(state_reaction, cookie.Value, id).Scan(&reaction_status)

	// If the User didn't reacted to that post yet
	if err_reaction == sql.ErrNoRows {
		// Creation of a new reaction
		state_create := `INSERT INTO Reaction (User_UUID, Post_ID, Status) VALUES (?, ?, ?)`
		_, err_creation := db.Exec(state_create, cookie.Value, id, "liked")
		if err_creation != nil {
			ErrorServer(w, "Error creating Post Reaction")
		}

		// Getting the Post like value
		var post_like int
		state_post_like := `SELECT Like FROM Posts WHERE User_UUID = ? AND ID = ?`
		err_post_like := db.QueryRow(state_post_like, cookie.Value, id).Scan(&post_like)
		if err_post_like != nil {
			ErrorServer(w, "Error getting Post Like value")
		}

		// Sending back the modified like value
		state_post_like = `UPDATE Posts SET Like = Like + 1 WHERE ID = ?`
		_, err_post_like = db.Exec(state_post_like, id)
		if err_post_like != nil {
			ErrorServer(w, "Error updating Post Like value")
		}

		// If the User already like the post
	} else if reaction_status == "liked" {

		// Getting the Post like value
		var post_like int
		state_post_like := `SELECT Like FROM Posts WHERE User_UUID = ? AND ID = ?`
		err_post_like := db.QueryRow(state_post_like, cookie.Value, id).Scan(&post_like)
		if err_post_like != nil {
			ErrorServer(w, "Error getting Post Like value")
		}

		// Sending back the modified like value
		state_post_like = `UPDATE Posts SET Like = Like - 1 WHERE ID = ?`
		_, err_post_like = db.Exec(state_post_like, id)
		if err_post_like != nil {
			ErrorServer(w, "Error updating Post Like value")
		}

		// Deleting the reaction table
		state_delete := `DELETE FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
		_, err_delete := db.Exec(state_delete, cookie.Value, id)
		if err_delete != nil {
			ErrorServer(w, "Error deleting the reaction")
		}

		// If the post is disliked
	} else if reaction_status == "disliked" {
		// Getting the Post like value
		var post_like int
		state_post_like := `SELECT Like FROM Posts WHERE User_UUID = ? AND ID = ?`
		err_post_like := db.QueryRow(state_post_like, cookie.Value, id).Scan(&post_like)
		if err_post_like != nil {
			ErrorServer(w, "Error getting Post Like value")
		}

		// Sending back the modified like value
		state_post_like = `UPDATE Posts SET Like = Like + 2 WHERE ID = ?`
		_, err_post_like = db.Exec(state_post_like, id)
		if err_post_like != nil {
			ErrorServer(w, "Error updating Post Like value")
		}
	} else {
		ErrorServer(w, "Error getting reaction status")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
