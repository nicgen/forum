package handlers

import (
	"forum/cmd/lib"
	"forum/models"
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
			// Erreur critique: See Other
			err := &models.CustomError{
				StatusCode: http.StatusSeeOther,
				Message:    "See Other",
			}

			HandleError(w, err.StatusCode, err.Message)
			return
		}

		// Verify if user is logged in the database
		var userUUID string
		state := "SELECT IsLogged FROM User WHERE UUID = ?"
		err_db := db.QueryRow(state, cookie.Value).Scan(&userUUID)
		if err_db != nil {
			//Erreur critique: Session not found
			err := &models.CustomError{
				StatusCode: http.StatusSeeOther,
				Message:    "Session not found or expired",
			}

			// Appel de HandleError
			HandleError(w, err.StatusCode, err.Message)
			return

		}

		// Session is valid, proceed to handler
		next.ServeHTTP(w, r)
	}
}
