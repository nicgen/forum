package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"time"
)

// ? Function to retrieve like/dislike infos and handler them
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	cookie, err := r.Cookie("session_id")
	if err != nil {
		//Erreur critique : Echec de la recup de cookie
		err := &models.CustomError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Session ID cookie not found",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Getting the form values
	id := r.FormValue("id")
	like_post := r.FormValue("like_post")
	dislike_post := r.FormValue("dislike_post")
	like_comment := r.FormValue("like_comment")
	dislike_comment := r.FormValue("dislike_comment")
	post_id := r.URL.Query().Get("post_id")

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
					//Erreur non critique : Echec de création de réaction
					lib.ErrorServer(w, "Error creating Post Reaction")
				}

				// Updating the post like value
				state_post_like := `UPDATE Posts SET Like = Like + 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					//Erreur non critique : Echec de mise à jour de réaction
					lib.ErrorServer(w, "Error updating Post Like value")
				}
			} else if dislike_post == "dislike_post" {
				// Creation of a new reaction
				state_create := `INSERT INTO Reaction (User_UUID, Post_ID, Status) VALUES (?, ?, ?)`
				_, err_creation := db.Exec(state_create, cookie.Value, id, "disliked")
				if err_creation != nil {
					//Erreur non critique: Echec de la création de réaction
					lib.ErrorServer(w, "Error creating Post Reaction")
				}

				// Updating the post like value
				state_post_like := `UPDATE Posts SET Dislike = Dislike + 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					//Erreur non critique : Echec de la mise a jour du nombre de dislike
					lib.ErrorServer(w, "Error updating Post Dislike value")
				}

				// Retrieve the last inserted comment ID
				reactionID, err := result.LastInsertId()
				if err != nil {
					lib.ErrorServer(w, "Error retrieving comment ID")
					return
				}

				// Insert a notification for the post author
				var postAuthorUUID string
				queryPostAuthor := `SELECT User_UUID FROM Posts WHERE ID = ?`
				err = db.QueryRow(queryPostAuthor, post_id).Scan(&postAuthorUUID)
				if err != nil {
					lib.ErrorServer(w, "Error finding post author")
					return
				}

				if postAuthorUUID != cookie.Value { // Avoid notifying the user if they commented on their own post
					state_notification := `INSERT INTO Notification (User_UUID, Comment_ID, Post_ID, CreatedAt, IsRead) VALUES (?, ?, ?, ?, ?)`
					_, errNotif := db.Exec(state_notification, postAuthorUUID, reactionID, post_id, time.Now(), false)
					if errNotif != nil {
						lib.ErrorServer(w, "Error inserting notification")
						return
					}
				}
				// // Notifier le propriétaire du post
				// var postOwnerUUID string
				// queryPostOwner := `SELECT User_UUID FROM Posts WHERE ID = ?`
				// err := db.QueryRow(queryPostOwner, id).Scan(&postOwnerUUID)
				// if err != nil {
				// 	lib.ErrorServer(w, "Error retrieving post owner")
				// 	return
				// }

				// if postOwnerUUID != cookie.Value { // Ne pas s'auto-notifier
				// 	reactionID := fmt.Sprintf("%s-%s", cookie.Value, id) // ID unique pour la réaction
				// 	errNotif := InsertNotification(db, postOwnerUUID, reactionID, sql.NullString{String: id, Valid: true}, sql.NullString{})
				// 	if errNotif != nil {
				// 		lib.ErrorServer(w, "Error creating notification for post owner")
				// 	}
				// }
			}

			// If the User already like the post
		} else if reaction_status == "liked" {
			if like_post == "like_post" {
				// Updating the post like value
				state_post_like := `UPDATE Posts SET Like = Like - 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					//Erreur non critique : Echec de la mise a jour du nombre de like
					lib.ErrorServer(w, "Error updating Post Like value")
				}

				// Deleting the reaction table
				state_delete := `DELETE FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
				_, err_delete := db.Exec(state_delete, cookie.Value, id)
				if err_delete != nil {
					//Erreur non critique : Echec de la suppression
					lib.ErrorServer(w, "Error deleting the reaction")
				}
			} else if dislike_post == "dislike_post" {

				// Updating the post like value
				state_post_like := `UPDATE Posts SET Like = Like - 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					//Erreur non critique : Echec de la mise a jour du nombre de like
					lib.ErrorServer(w, "Error updating Post Like value 1")
				}

				// Updating the post dislike value
				state_post_dislike := `UPDATE Posts SET Dislike = Dislike + 1 WHERE ID = ?`
				_, err_post_dislike := db.Exec(state_post_dislike, id)
				if err_post_dislike != nil {
					//Erreur non critique : Echec de la mise a jour de nombre de dislike
					lib.ErrorServer(w, "Error updating Post Dislike value 2")
				}

				// Deleting the reaction table
				state_delete := `UPDATE Reaction SET Status = ? WHERE User_UUID = ? AND Post_ID = ?`
				_, err_delete := db.Exec(state_delete, "disliked", cookie.Value, id)
				if err_delete != nil {
					// Erreur non critique : Erreur de la mise a jour de la réaction
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
					//Erreur non critique : Echec de la mise a jour du nbr de like
					lib.ErrorServer(w, "Error updating Post Like value 1")
				}

				// And also updating the dislike
				state_post_dislike := `UPDATE Posts SET Dislike = Dislike - 1 WHERE ID = ?`
				_, err_post_dislike := db.Exec(state_post_dislike, id)
				if err_post_dislike != nil {
					//Erreur non critique : Echec de la mise a jour du nbr de like
					lib.ErrorServer(w, "Error updating Post Dislike value 2")
				}

				// Changing the state into "liked"
				state_update_state := `UPDATE Reaction SET Status = ? WHERE User_UUID = ? AND Post_ID = ?`
				_, err_update_state := db.Exec(state_update_state, "liked", cookie.Value, id)
				if err_update_state != nil {
					//Erreur non critique: Echec de la mise a jour de la réaction
					lib.ErrorServer(w, "Error updating Post Like value 3")
				}
			} else if dislike_post == "dislike_post" {

				// Updating the post like value
				state_post_like := `UPDATE Posts SET Dislike = Dislike - 1 WHERE ID = ?`
				_, err_post_like := db.Exec(state_post_like, id)
				if err_post_like != nil {
					//Erreur non critique : Echec de la mise a jour du nbr de dislikes
					lib.ErrorServer(w, "Error updating Post Dislike value")
				}

				// Deleting the reaction table
				state_delete := `DELETE FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
				_, err_delete := db.Exec(state_delete, cookie.Value, id)
				if err_delete != nil {
					//Erreur non critique : Echec de la suppression de réaction
					lib.ErrorServer(w, "Error deleting the reaction")
				}
			}

		} else {
			//Erreur critique : Etat de réaction Inconnu
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error getting reaction status",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
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

func InsertNotification(db *sql.DB, recipientUUID string, reactionID string, postID sql.NullString, commentID sql.NullString) error {
	query := `INSERT INTO Notification (User_UUID, ReactionID, Post_ID, Comment_ID, CreatedAt, IsRead) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, recipientUUID, reactionID, postID, commentID, time.Now(), false)
	return err
}
