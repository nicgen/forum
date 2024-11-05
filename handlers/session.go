package handlers

import (
	"forum/cmd/lib"
	"net/http"
	"time"
)

// ? Function to attribute a temporary cookie that will store User UUID in header
func CookieSession(user_uuid string, w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Setting the User UUID into the cookie
	cookie := &http.Cookie{
		Name:     "session_id", // Name of the cookie
		Value:    user_uuid,    // Using UUID as session token
		Path:     "/",          // Cookie is valid for all paths
		HttpOnly: true,         // Cannot be accessed by JavaScript
		Secure:   true,         // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,
		// Expires in 24 hours
		Expires: time.Now().Add(24 * time.Hour),
	}

	// Setting the User as Logged in the database
	state := `UPDATE User SET IsLogged = ? WHERE UUID = ?`
	_, err_db := db.Exec(state, true, user_uuid)
	if err_db != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	// Set the cookie in the response header
	http.SetCookie(w, cookie)
}
