package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"net/http"
	"strconv"
	"strings"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	cookie, _ := r.Cookie("session_id")

	// Getting the form values
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	post_id, _ := strconv.Atoi(id)

	println(post_id)

	// Checking if the User already liked the post
	var is_liked bool = true
	state_isliked := `SELECT IsLiked FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
	err_isliked := db.QueryRow(state_isliked, cookie.Value, post_id).Scan(&is_liked)

	println("is liked: ", is_liked)

	if err_isliked == sql.ErrNoRows {
		state_isliked := `INSERT INTO Reaction (Post_ID, User_UUID, IsLiked) VALUES (?, ?, ?)`
		_, err_add_react := db.Exec(state_isliked, post_id, cookie.Value, is_liked)
		if err_add_react != nil {
			ErrorServer(w, "Error creating User reaction")
		}
	} else if err_isliked != nil {
		ErrorServer(w, "Error checking User like")
	}

	// If the post is already liked by the User, simply redirect to home page
	println("is liked: ", is_liked)
	if !is_liked {
		// Redirect User to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {

		// Getting the number of likes of the post in the database
		var like_number int
		state_like := `SELECT Like FROM Posts WHERE User_UUID = ?`
		err_like := db.QueryRow(state_like, cookie.Value).Scan(&like_number)
		if err_like != nil {
			ErrorServer(w, "Error accessing User ID")
		}

		// Creating a reaction to keep track of liked posts
		state_reaction := `INSERT INTO Reaction (Post_ID, User_UUID, IsLiked) VALUES (?, ?, ?)`
		_, err_db := db.Exec(state_reaction, post_id, cookie.Value, false)
		if err_db != nil {
			ErrorServer(w, "Error liking the post")
		}

		// Updating the number of likes
		state_add_like := `UPDATE Posts SET Like = ? WHERE ID = ?`
		_, err_addlike := db.Exec(state_add_like, like_number+1, post_id)
		if err_addlike != nil {
			ErrorServer(w, "Error liking the post")
		}

		// Redirect User to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
