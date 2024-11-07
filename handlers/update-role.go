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
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		http.Error(w, "User UUID is required.", http.StatusBadRequest)
		return
	}

	// Update the user's role to "Moderator"
	query := `UPDATE User SET Role = 'Moderator' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		http.Error(w, "Failed to update user role.", http.StatusInternalServerError)
		fmt.Println("Error updating user role:", err)
		return
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func RemoveModerator(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		http.Error(w, "User UUID is required.", http.StatusBadRequest)
		return
	}

	// Update the user's role to "Moderator"
	query := `UPDATE User SET Role = 'User' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		http.Error(w, "Failed to update user role.", http.StatusInternalServerError)
		fmt.Println("Error updating user role:", err)
		return
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
