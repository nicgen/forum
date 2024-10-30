package handlers

import (
	"forum/cmd/lib"
	"net/http"
	"time"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	session_id := r.Cookies()

	title := r.FormValue("post_title")
	text := r.FormValue("post_text")
	category := r.FormValue("post_category")

	state_post := `INSERT INTO Posts (User_UUID, Title, Category_ID, Text, CreatedAt) VALUES (?, ?, ?, ?, ?)`
	_, err_db := db.Exec(state_post, session_id[0].Value, title, category, text, time.Now())
	if err_db != nil {
		http.Error(w, "Error inserting new Post", http.StatusUnauthorized)
		return
	}

	// Printing Post values
	println(title)
	println(text)
	println(category)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
