package handlers

import (
	"fmt"
	"forum/cmd/lib"
	"net/http"
	"time"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database data into a variable
	db := lib.GetDB()

	// Getting the cookie (containing the UUID)
	cookie, _ := r.Cookie("session_id")

	// Storing the form values into variables
	text := r.FormValue("comment_text")
	post_id := r.URL.Query().Get("post_id")

	fmt.Println("post_id: ", post_id)

	var like_count, dislike_count int = 0, 0
	// Insert the new comment into the Comments table
	state_post := `INSERT INTO Comments (User_UUID, Post_ID, Text, Like, Dislike, CreatedAt) VALUES (?, ?, ?, ?, ?, ?)`
	result, err_db := db.Exec(state_post, cookie.Value, post_id, text, like_count, dislike_count, time.Now())
	if err_db != nil {
		lib.ErrorServer(w, "Error inserting new Comment")
		return
	}

	// Retrieve the last inserted comment ID
	commentID, err := result.LastInsertId()
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
		_, errNotif := db.Exec(state_notification, postAuthorUUID, commentID, post_id, time.Now(), false)
		if errNotif != nil {
			lib.ErrorServer(w, "Error inserting notification")
			return
		}
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
