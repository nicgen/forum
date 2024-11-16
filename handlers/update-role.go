package handlers

import (
	"database/sql"
	"fmt"
	"forum/cmd/lib"
	"net/http"
	"time"
)

// Function to handle role updates for users
func UpdateUserToModerator(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		lib.ErrorServer(w, "Invalid request method.")
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		lib.ErrorServer(w, "User UUID is required.")
	}

	// Update the user's role to "Moderator"
	query := `UPDATE User SET Role = 'Moderator' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		lib.ErrorServer(w, "Failed to update user role.")
		fmt.Println("Error updating user role:", err)
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func RemoveModerator(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		lib.ErrorServer(w, "Invalid request method.")
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		lib.ErrorServer(w, "User UUID is required.")
	}

	query := `UPDATE User SET Role = 'User' WHERE UUID = ?`
	_, err := db.Exec(query, userUUID)
	if err != nil {
		lib.ErrorServer(w, "Failed to update user role.")
		fmt.Println("Error updating user role:", err)
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		lib.ErrorServer(w, "Invalid request method.")
	}

	// Retrieve the user's UUID from the form
	userUUID := r.FormValue("userUUID")
	fmt.Println("Received userUUID:", userUUID)
	if userUUID == "" {
		lib.ErrorServer(w, "User UUID is required.")
	}

	// Génère un nom d'utilisateur anonyme unique
	newUsername, err := getNextAnonymousUsername(db)
	if err != nil {
		lib.ErrorServer(w, "Failed to generate anonymous username.")
		fmt.Println("Error generating anonymous username:", err)
		return
	}

	query := `UPDATE User SET Role = 'DeleteUser', username = ? WHERE UUID = ?`
	_, err = db.Exec(query, newUsername, userUUID)
	if err != nil {
		lib.ErrorServer(w, "Failed to update user role.")
		fmt.Println("Error updating user role:", err)
	}

	// Redirect or send a success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func getNextAnonymousUsername(db *sql.DB) (string, error) {
	var lastNumber int
	query := `SELECT COALESCE(MAX(CAST(SUBSTRING(username, 8) AS UNSIGNED)), 0) FROM User WHERE username LIKE 'Anonyme%'`
	err := db.QueryRow(query).Scan(&lastNumber)
	if err != nil {
		return "", err
	}
	nextNumber := lastNumber + 1
	return fmt.Sprintf("Anonyme%d", nextNumber), nil
}

// Fonction pour supprimer un utilisateur en le rendant anonyme
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Vérifie que la méthode de requête est POST
	if r.Method != http.MethodPost {
		lib.ErrorServer(w, "Invalid request method.")
		return
	}

	// Récupère l'UUID de l'utilisateur depuis le formulaire
	userUUID := r.FormValue("userUUID")
	fmt.Println("Received userUUID:", userUUID)
	if userUUID == "" {
		lib.ErrorServer(w, "User UUID is required.")
		return
	}

	// Génère un nom d'utilisateur anonyme unique
	newUsername, err := getNextAnonymousUsername(db)
	if err != nil {
		lib.ErrorServer(w, "Failed to generate anonymous username.")
		fmt.Println("Error generating anonymous username:", err)
		return
	}

	// Met à jour le rôle et le nom de l'utilisateur avec un numéro incrémenté
	query := `UPDATE User SET Role = 'DeleteUser', username = ? WHERE UUID = ?`
	_, err = db.Exec(query, newUsername, userUUID)
	if err != nil {
		lib.ErrorServer(w, "Failed to update user role.")
		fmt.Println("Error updating user role:", err)
		return
	}

	// Supprime le cookie de session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immédiatement
	})

	// Redirige ou envoie un message de succès
	lib.LogoutHandler(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
