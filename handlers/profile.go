package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strings"
)

// ? Function to retrieve User data and redirect to his profile page
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	session_id := r.Cookies()

	// Setting up the variables of the User informations
	var user_username, user_email, user_creation, user_role string

	// Setting up the states for the database requests
	state_uuid := `SELECT Username FROM User WHERE UUID = ?`
	state_email := `SELECT Email FROM User WHERE UUID = ?`
	state_creation := `SELECT CreatedAt FROM User WHERE UUID = ?`
	state_role := `SELECT Role FROM User WHERE UUID = ?`

	// Users posts Request
	state_posts := `SELECT ID, Category_ID, Title, Text, CreatedAt FROM Posts WHERE User_UUID = ? ORDER BY CreatedAt DESC`
	var posts []*models.Post
	rows, err := db.Query(state_posts, session_id[0].Value)
	if err != nil {
		http.Error(w, "Error accessing user posts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.CreatedAt); err != nil {
			http.Error(w, "Error scanning user posts: "+err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over user posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Making requests to the database
	err_uuid := db.QueryRow(state_uuid, session_id[0].Value).Scan(&user_username)
	err_email := db.QueryRow(state_email, session_id[0].Value).Scan(&user_email)
	err_creation := db.QueryRow(state_creation, session_id[0].Value).Scan(&user_creation)
	err_role := db.QueryRow(state_role, session_id[0].Value).Scan(&user_role)

	// Checking for database requests errors
	if err_uuid != nil {
		http.Error(w, "Error accessing User UUID", http.StatusUnauthorized)
		return
	} else if err_email != nil {
		http.Error(w, "Error accessing User EMAIL", http.StatusUnauthorized)
		return
	} else if err_creation != nil {
		http.Error(w, "Error accessing User CREATION DATE", http.StatusUnauthorized)
		return
	} else if err_role != nil {
		http.Error(w, "Error accessing User ROLE", http.StatusUnauthorized)
		return
	}

	// Spliting the creation date into 2 different values
	creation := strings.Split(user_creation, "T")
	creation_Date := creation[0]
	creation_Hour := creation[1][:len(creation[1])-1]

	// Storing the data into a map that can be sent into the html
	data := map[string]interface{}{
		"Username":     user_username,
		"Email":        user_email,
		"CreationDate": creation_Date,
		"CreationHour": creation_Hour,
		"Role":         user_role,
		"Posts":        posts,
	}

	// Redirect User to the profile html page and sending the data to it
	renderTemplate(w, "layout/index", "page/profile", data)
}
