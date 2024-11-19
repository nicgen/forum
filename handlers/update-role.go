package handlers

import (
	"fmt"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"time"
)

// Function to handle role updates for users
func UpdateUserToModerator(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		// Erreur non critique : Méthode de requête invalide
		lib.ErrorServer(w, "Invalid request method.")
		return
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		// Erreur non critique : UUID de l'utilisateur requis
		lib.ErrorServer(w, "User  UUID is required.")
		return
	}

	// Update the user's role to "Moderator"
	query := `UPDATE User SET Role = 'Moderator' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		// Erreur critique : Échec de la mise à jour du rôle de l'utilisateur
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update user role to Moderator",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func RemoveModerator(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		// Erreur non critique : Méthode de requête invalide
		lib.ErrorServer(w, "Invalid request method.")
		return
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		// Erreur non critique : UUID de l'utilisateur requis
		lib.ErrorServer(w, "User  UUID is required.")
		return
	}

	query := `UPDATE User SET Role = 'User' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		// Erreur critique : Échec de la mise à jour du rôle de l'utilisateur
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update user role to User",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		// Erreur non critique : Méthode de requête invalide
		lib.ErrorServer(w, "Invalid request method.")
		return
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	fmt.Println("Received userUUID:", userUUID)
	if userUUID == "" {
		// Erreur non critique : UUID de l'utilisateur requis
		lib.ErrorServer(w, "User  UUID is required.")
		return
	}

	query := `UPDATE User SET Role = 'DeletUser' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		// Erreur critique : Échec de la mise à jour du rôle de l'utilisateur
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update user role to DeletUser ",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		// Erreur non critique : Méthode de requête invalide
		lib.ErrorServer(w, "Invalid request method.")
		return
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	fmt.Println("Received userUUID:", userUUID)
	if userUUID == "" {
		// Erreur non critique : UUID de l'utilisateur requis
		lib.ErrorServer(w, "User  UUID is required.")
		return
	}

	query := `UPDATE User SET Role = 'DeletUser' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		// Erreur critique : Échec de la mise à jour du rôle de l'utilisateur
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update user role to DeletUser ",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
	})

	// Redirect or send a success message
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
