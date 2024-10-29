package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()

	email := r.FormValue("EmailForm")
	password := r.FormValue("PasswordForm")

	// if email == "" || password == "" {
	// 	http.Error(w, "Email and password fields are mandatory", http.StatusBadRequest)
	// 	return
	// }

	// Prepared request to avoid SQL injection
	stmt, err := db.Prepare("SELECT password FROM User WHERE email = ?")
	if err != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var hashedPassword string
	err = stmt.QueryRow(email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
		}
		return
	}

	// Verify password
	if !CheckPassword(hashedPassword, password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	var user_uuid string
	state := `SELECT UUID FROM User`
	err_user := db.QueryRow(state).Scan(&user_uuid)
	if err_user != nil {
		http.Error(w, "Error accessing User UUID", http.StatusUnauthorized)
		return
	}

	println("User UUID: ", user_uuid)

	// Attribute a session to an User
	CookieSession(user_uuid, w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Function to check if the password matches the stored hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
