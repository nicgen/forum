package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// ? Function that will verify the form values for the login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
	db := lib.GetDB()

	// Checking if the User is already logged or not
	_, err_cookie := r.Cookie("session_id")
	if err_cookie == http.ErrNoCookie {
		// Storing form values into variables
		email := r.FormValue("EmailForm")
		password := r.FormValue("PasswordForm")

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

		// Getting the UUID from the database
		var user_uuid string
		state := `SELECT UUID FROM User WHERE Email = ?`
		err_user := db.QueryRow(state, email).Scan(&user_uuid)
		if err_user != nil {
			http.Error(w, "Error accessing User UUID", http.StatusUnauthorized)
			return
		}

		// Attribute a session to an User
		CookieSession(user_uuid, w, r)

	} else {
		// If the User is already logged and tries to log-in
		http.Error(w, "You must log-out before loggin in again", http.StatusUnauthorized)
		return
	}

	// Redirecting to the home page after successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Function to check if the password matches the stored hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
