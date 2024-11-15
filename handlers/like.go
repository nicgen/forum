package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	cookie, _ := r.Cookie("session_id")

	// Getting the form values
	id := r.FormValue("id")
	like_post := r.FormValue("like_post")
	dislike_post := r.FormValue("dislike_post")
	like_comment := r.FormValue("like_comment")
	dislike_comment := r.FormValue("dislike_comment")

	if like_post != "" || dislike_post != "" {
		var reaction_status string
		state_reaction := `SELECT Status FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
		err_reaction := db.QueryRow(state_reaction, cookie.Value, id).Scan(&reaction_status)

		// If the User didn't reacted to that post yet
		if err_reaction == sql.ErrNoRows {
			if like_post == "like_post" {
				// Creation of a new reaction
				state_create := `INSERT INTO Reaction (User_UUID, Post_ID, Status) VALUES (?, ?, ?)`
				_, err_creation := db.Exec(state_create, cookie.Value, id, "liked")
				if err_creation != nil {
					lib.ErrorServer(w, "Error creating Post Reaction")
				}

				// Updating the post like value
				state_post_like := `UPDATE Posts SET Like = Like + 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					lib.ErrorServer(w, "Error updating Post Like value")
				}
			} else if dislike_post == "dislike_post" {
				// Creation of a new reaction
				state_create := `INSERT INTO Reaction (User_UUID, Post_ID, Status) VALUES (?, ?, ?)`
				_, err_creation := db.Exec(state_create, cookie.Value, id, "disliked")
				if err_creation != nil {
					lib.ErrorServer(w, "Error creating Post Reaction")
				}

				// Updating the post like value
				state_post_like := `UPDATE Posts SET Dislike = Dislike + 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					lib.ErrorServer(w, "Error updating Post Dislike value")
				}
			}

			// If the User already like the post
		} else if reaction_status == "liked" {

			if like_post == "like_post" {
				// Updating the post like value
				state_post_like := `UPDATE Posts SET Like = Like - 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					lib.ErrorServer(w, "Error updating Post Like value")
				}

				// Deleting the reaction table
				state_delete := `DELETE FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
				_, err_delete := db.Exec(state_delete, cookie.Value, id)
				if err_delete != nil {
					lib.ErrorServer(w, "Error deleting the reaction")
				}
			} else if dislike_post == "dislike_post" {

				// Updating the post like value
				state_post_like := `UPDATE Posts SET Like = Like - 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					lib.ErrorServer(w, "Error updating Post Like value 1")
				}

				// Updating the post dislike value
				state_post_dislike := `UPDATE Posts SET Dislike = Dislike + 1 WHERE ID = ?`
				_, err_post_dislike := db.Exec(state_post_dislike, id)
				if err_post_dislike != nil {
					lib.ErrorServer(w, "Error updating Post Dislike value 2")
				}

				// Deleting the reaction table
				state_delete := `UPDATE Reaction SET Status = ? WHERE User_UUID = ? AND Post_ID = ?`
				_, err_delete := db.Exec(state_delete, "disliked", cookie.Value, id)
				if err_delete != nil {
					lib.ErrorServer(w, "Error updating the reaction from like to dislike")
				}
			}

			// If the post is disliked
		} else if reaction_status == "disliked" {

			if like_post == "like_post" {
				// Updating the post like value
				state_post_like := `UPDATE Posts SET Like = Like + 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					lib.ErrorServer(w, "Error updating Post Like value 1")
				}

				// And also updating the dislike
				state_post_dislike := `UPDATE Posts SET Dislike = Dislike - 1 WHERE ID = ?`
				_, err_post_dislike := db.Exec(state_post_dislike, id)
				if err_post_dislike != nil {
					lib.ErrorServer(w, "Error updating Post Like value 2")
				}

				// Changing the state into "liked"
				state_update_state := `UPDATE Reaction SET Status = ? WHERE User_UUID = ? AND Post_ID = ?`
				_, err_update_state := db.Exec(state_update_state, "liked", cookie.Value, id)
				if err_update_state != nil {
					lib.ErrorServer(w, "Error updating Post Like value 3")
				}
			} else if dislike_post == "dislike_post" {

				// Updating the post like value
				state_post_like := `UPDATE Posts SET Dislike = Dislike - 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					lib.ErrorServer(w, "Error updating Post Dislike value")
				}

				// Deleting the reaction table
				state_delete := `DELETE FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
				_, err_delete := db.Exec(state_delete, cookie.Value, id)
				if err_delete != nil {
					lib.ErrorServer(w, "Error deleting the reaction")
				}
			}

		} else {
			lib.ErrorServer(w, "Error getting reaction status")
		}
	} else {
		var reaction_status string
		state_reaction := `SELECT Status FROM Reaction WHERE User_UUID = ? AND Comment_ID = ?`
		err_reaction := db.QueryRow(state_reaction, cookie.Value, id).Scan(&reaction_status)

		// If the User didn't reacted to that post yet
		if err_reaction == sql.ErrNoRows {
			if like_comment == "like_comment" {
				// Creation of a new reaction
				state_create := `INSERT INTO Reaction (User_UUID, Comment_ID, Status) VALUES (?, ?, ?)`
				_, err_creation := db.Exec(state_create, cookie.Value, id, "liked")
				if err_creation != nil {
					lib.ErrorServer(w, "Error creating Comment Reaction")
				}

				// Updating the post like value
				state_comment_like := `UPDATE Comments SET Like = Like + 1 WHERE ID = ?`
				_, err_comment_like := db.Exec(state_comment_like, id)
				if err_comment_like != nil {
					lib.ErrorServer(w, "Error updating Comment Like value")
				}
			} else if dislike_post == "dislike_comment" {
				// Creation of a new reaction
				state_create := `INSERT INTO Reaction (User_UUID, Comment_ID, Status) VALUES (?, ?, ?)`
				_, err_creation := db.Exec(state_create, cookie.Value, id, "disliked")
				if err_creation != nil {
					lib.ErrorServer(w, "Error creating Comment Reaction")
				}

				// Updating the post like value
				state_comment_like := `UPDATE Comments SET Dislike = Dislike + 1 WHERE ID = ?`
				_, err_comment_like := db.Exec(state_comment_like, id)
				if err_comment_like != nil {
					lib.ErrorServer(w, "Error updating Comment Dislike value")
				}
			}

			// If the User already like the post
		} else if reaction_status == "liked" {

			if like_comment == "like_comment" {
				// Updating the post like value
				state_comment_like := `UPDATE Comments SET Like = Like - 1 WHERE ID = ?`
				_, err_comment_like := db.Exec(state_comment_like, id)
				if err_comment_like != nil {
					lib.ErrorServer(w, "Error updating Comment Like value")
				}

				// Deleting the reaction table
				state_delete := `DELETE FROM Reaction WHERE User_UUID = ? AND Comment_ID = ?`
				_, err_delete := db.Exec(state_delete, cookie.Value, id)
				if err_delete != nil {
					lib.ErrorServer(w, "Error deleting the reaction")
				}
			} else if dislike_comment == "dislike_comment" {

				// Updating the post like value
				state_comment_like := `UPDATE Comments SET Like = Like - 1 WHERE ID = ?`
				_, err_comment_like := db.Exec(state_comment_like, id)
				if err_comment_like != nil {
					lib.ErrorServer(w, "Error updating Comment Like value 1")
				}

				// Updating the post dislike value
				state_comment_dislike := `UPDATE Comments SET Dislike = Dislike + 1 WHERE ID = ?`
				_, err_comment_dislike := db.Exec(state_comment_dislike, id)
				if err_comment_dislike != nil {
					lib.ErrorServer(w, "Error updating Comment Dislike value 2")
				}

				// Deleting the reaction table
				state_delete := `UPDATE Reaction SET Status = ? WHERE User_UUID = ? AND Comment_ID = ?`
				_, err_delete := db.Exec(state_delete, "disliked", cookie.Value, id)
				if err_delete != nil {
					lib.ErrorServer(w, "Error updating the reaction from like to dislike")
				}
			}

			// If the post is disliked
		} else if reaction_status == "disliked" {

			if like_comment == "like_comment" {
				// Updating the post like value
				state_comment_like := `UPDATE Comments SET Like = Like + 1 WHERE ID = ?`
				_, err_comment_like := db.Exec(state_comment_like, id)
				if err_comment_like != nil {
					lib.ErrorServer(w, "Error updating Comment Like value 1")
				}

				// And also updating the dislike
				state_post_dislike := `UPDATE Comments SET Dislike = Dislike - 1 WHERE ID = ?`
				_, err_comment_dislike := db.Exec(state_post_dislike, id)
				if err_comment_dislike != nil {
					lib.ErrorServer(w, "Error updating Comment Like value 2")
				}

				// Changing the state into "liked"
				state_update_state := `UPDATE Reaction SET Status = ? WHERE User_UUID = ? AND Comment_ID = ?`
				_, err_update_state := db.Exec(state_update_state, "liked", cookie.Value, id)
				if err_update_state != nil {
					lib.ErrorServer(w, "Error updating Comment Like value 3")
				}
			} else if dislike_comment == "dislike_comment" {

				// Updating the post like value
				state_comment_like := `UPDATE Posts SET Dislike = Dislike - 1 WHERE ID = ?`
				_, err_comment_like := db.Exec(state_comment_like, id)
				if err_comment_like != nil {
					lib.ErrorServer(w, "Error updating Comment Dislike value")
				}

				// Deleting the reaction table
				state_delete := `DELETE FROM Reaction WHERE User_UUID = ? AND Comment_ID = ?`
				_, err_delete := db.Exec(state_delete, cookie.Value, id)
				if err_delete != nil {
					lib.ErrorServer(w, "Error deleting the reaction")
				}
			}

		} else {
			lib.ErrorServer(w, "Error getting reaction status")
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
