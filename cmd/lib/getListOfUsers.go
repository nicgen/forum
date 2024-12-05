package lib

import (
	"forum/models"
	"net/http"
)

// ? Function to store the list of Users into the data map if the role is Admin
func GetListOfUsers(w http.ResponseWriter, role string, data map[string]interface{}) map[string]interface{} {
	// Retrieve all users for admin view
	var allUsers []models.User
	if role == "Admin" {
		allUsersQuery := `SELECT UUID, Username, Email, Role, IsRequest FROM User`
		rows, err := db.Query(allUsersQuery)
		if err != nil {
			ErrorServer(w, "Error accessing user list")
		}
		defer rows.Close()

		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.UUID, &user.Username, &user.Email, &user.Role, &user.IsRequest); err != nil {
				ErrorServer(w, "Error scanning users")
			}
			allUsers = append(allUsers, user)
		}

		if err := rows.Err(); err != nil {
			ErrorServer(w, "Error iterating over users")
		}
	}
	data["AllUsers"] = allUsers
	return data
}
