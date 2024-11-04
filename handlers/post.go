package handlers

import (
	"forum/cmd/lib"
	"net/http"
	"time"
)

// ? Handler that will insert a new post into the database
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
	db := lib.GetDB()

	// Getting the cookie (containing the UUID)
	session_id := r.Cookies()

	// Storing the form values into variables
	title := r.FormValue("post_title")
	text := r.FormValue("post_text")
	category := r.FormValue("post_category")

	// Storing those values into the database with a database request
	state_post := `INSERT INTO Posts (User_UUID, Title, Category_ID, Text, CreatedAt) VALUES (?, ?, ?, ?, ?)`
	_, err_db := db.Exec(state_post, session_id[0].Value, title, category, text, time.Now())
	if err_db != nil {
		http.Error(w, "Error inserting new Post", http.StatusUnauthorized)
		return
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}