package handlers

import (
	"forum/cmd/lib"
	"net/http"
	"time"
)

// ? Handler that will insert a new post into the database
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
	db := lib.GetDB()

	// Getting the cookie (containing the UUID)
	cookie, _ := r.Cookie("session_id")

	// Storing the form values into variables
	title := r.FormValue("post_title")
	text := r.FormValue("post_text")
	category := r.FormValue("post_category")

	var like_count, dislike_count int = 0, 0
	// Storing those values into the database with a database request
	state_post := `INSERT INTO Posts (User_UUID, Title, Category_ID, Text, Like, Dislike, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err_db := db.Exec(state_post, cookie.Value, title, category, text, like_count, dislike_count, time.Now())
	if err_db != nil {
		lib.ErrorServer(w, "Error inserting new Post")
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
