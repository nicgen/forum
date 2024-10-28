package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	print("test")

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

	// Successful login, create a session token
	sessionUUID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Error generating session token", http.StatusInternalServerError)
		return
	}
	sessionToken := sessionUUID.String()

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Function to check if the password matches the stored hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
