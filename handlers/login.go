package handlers

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Display the login form
		http.ServeFile(w, r, "./templates/login.html") // Path to your login form
		return
	}

	if r.Method == "POST" {
		// Retrieve form data
		username := r.FormValue("UsernameForm")
		password := r.FormValue("PasswordForm")

		// Validate input
		if username == "" || password == "" {
			http.Error(w, "Username and password fields are mandatory", http.StatusBadRequest)
			return
		}

		// Retrieve the hashed password from the database
		var hashedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			} else {
				http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
			}
			return
		}

		// Check if the provided password matches the stored hashed password
		if !CheckPassword(hashedPassword, password) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Successful login
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

// Function to check if the password matches the stored hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
