package lib

import (
	"database/sql"
	"forum/models"
	"strings"
)

func GetData(db *sql.DB, uuid string) (any any, err string) {
	// Setting up the variables of the User informations
	var user_username, user_email, user_creation, user_role string

	// SQL queries to get user information
	state_uuid := `SELECT Username FROM User WHERE UUID = ?`
	state_email := `SELECT Email FROM User WHERE UUID = ?`
	state_creation := `SELECT CreatedAt FROM User WHERE UUID = ?`
	state_role := `SELECT Role FROM User WHERE UUID = ?`

	// Users posts Request
	state_posts := `SELECT ID, Category_ID, Title, Text, CreatedAt FROM Posts WHERE User_UUID = ? ORDER BY CreatedAt DESC`
	var posts []*models.Post
	rows, err_post := db.Query(state_posts, uuid)
	if err_post != nil {
		return nil, "Error accessing user posts"
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.CreatedAt); err != nil {
			return nil, "Error scanning user posts"
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, "Error iterating over user posts"
	}

	// Making requests to the database
	err_uuid := db.QueryRow(state_uuid, uuid).Scan(&user_username)
	err_email := db.QueryRow(state_email, uuid).Scan(&user_email)
	err_creation := db.QueryRow(state_creation, uuid).Scan(&user_creation)
	err_role := db.QueryRow(state_role, uuid).Scan(&user_role)

	// Checking for database requests errors
	if err_uuid != nil {
		return nil, "Error accessing User UUID"
	} else if err_email != nil {
		return nil, "Error accessing User EMAIL"
	} else if err_creation != nil {
		return nil, "Error accessing User CREATION DATE"
	} else if err_role != nil {
		return nil, "Error accessing User ROLE"
	}

	// Splitting the creation date into 2 different values
	creation := strings.Split(user_creation, "T")
	creation_Date := creation[0]
	creation_Hour := creation[1][:len(creation[1])-1]
	// Retrieve all users for admin view
	var allUsers []models.User
	if user_role == "Admin" {
		allUsersQuery := `SELECT UUID, Username, Email, Role FROM User`
		rows, err := db.Query(allUsersQuery)
		if err != nil {
			return nil, "Error accessing user list"
		}
		defer rows.Close()

		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.UUID, &user.Username, &user.Email, &user.Role); err != nil {
				return nil, "Error scanning users"
			}
			allUsers = append(allUsers, user)
		}

		if err := rows.Err(); err != nil {
			return nil, "Error iterating over users"
		}
	}

	// Storing the data into a map that can be sent into the html
	data := map[string]interface{}{
		"Username":     user_username,
		"Email":        user_email,
		"CreationDate": creation_Date,
		"CreationHour": creation_Hour,
		"Role":         user_role,
		"Posts":        posts,
		"AllUsers":     allUsers,
	}
	return data, "OK"
}
