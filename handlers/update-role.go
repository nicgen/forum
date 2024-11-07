package handlers

import (
	"fmt"
	"forum/cmd/lib"
	"net/http"
)

// Function to handle role updates for users
func UpdateUserToModerator(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		ErrorServer(w, "Invalid request method.")
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		ErrorServer(w, "User UUID is required.")
	}

	// Update the user's role to "Moderator"
	query := `UPDATE User SET Role = 'Moderator' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		ErrorServer(w, "Failed to update user role.")
		fmt.Println("Error updating user role:", err)
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func RemoveModerator(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		ErrorServer(w, "Invalid request method.")
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		ErrorServer(w, "User UUID is required.")
	}

	// Update the user's role to "Moderator"
	query := `UPDATE User SET Role = 'User' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		ErrorServer(w, "Failed to update user role.")
		fmt.Println("Error updating user role:", err)
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
