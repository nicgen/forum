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
	// Storing those values into the database with a database request
	state_post := `INSERT INTO Comments (User_UUID, Post_ID, Text, Like, Dislike, CreatedAt) VALUES (?, ?, ?, ?, ?, ?)`
	_, err_db := db.Exec(state_post, cookie.Value, post_id, text, like_count, dislike_count, time.Now())
	if err_db != nil {
		lib.ErrorServer(w, "Error inserting new Comment")
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
