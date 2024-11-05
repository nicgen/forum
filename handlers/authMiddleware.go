package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

// ? Function to check the User session before accessing informations
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Storing database into a variable
		db := lib.GetDB()

		// Get session cookie
		cookie, err_cookie := r.Cookie("session_id")
		if err_cookie != nil {
			// Redirect User to login page if the cookie doesn't exist
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Verify if user is logged in the database
		var userUUID string
		state := "SELECT IsLogged FROM User WHERE UUID = ?"
		err_db := db.QueryRow(state, cookie.Value).Scan(&userUUID)
		if err_db != nil {
			// Session not found or expired
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Session is valid, proceed to handler
		next.ServeHTTP(w, r)
	}
}
