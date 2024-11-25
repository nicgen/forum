package handlers

import (
	"fmt"
	"forum/cmd/lib"
	"net/http"
	"time"
)

// ? Handler that will insert a new post into the database
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
	db := lib.GetDB()

    // Parse the form data (including query parameters and form body)
    err := r.ParseForm()
    if err != nil {
        fmt.Println("Error parsing form:", err)
        http.Error(w, "Unable to parse form data", http.StatusInternalServerError)
        return
    }
	
	// Getting the cookie (containing the UUID)
	cookie, _ := r.Cookie("session_id")
	
	// Storing the form values into variables
	title := r.FormValue("post_title")
	text := r.FormValue("post_text")
    selectedCategories := r.Form["categories"] // This gives you a slice of strings

	var categories string
	for i := 0; i < len(selectedCategories); i++ {
		categories += selectedCategories[i] + ","
	}

	var like_count, dislike_count int = 0, 0
	// Storing those values into the database with a database request
	state_post := `INSERT INTO Posts (User_UUID, Title, Category_ID, Text, Like, Dislike, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err_db := db.Exec(state_post, cookie.Value, title, categories[:len(categories)-1], text, like_count, dislike_count, time.Now())
	if err_db != nil {
		lib.ErrorServer(w, "Error inserting new Post")
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
