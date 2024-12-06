package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
)

// ? Function to check the User session before accessing informations
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Storing database into a variable
		db := lib.GetDB()

		// Checking if the User have a cookie or not
		cookie, err_cookie := r.Cookie("session_id")
		cookie_role, err_cookie_role := r.Cookie("role")
		if err_cookie == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else if err_cookie != nil {
			// Erreur critique: See Other
			err := &models.CustomError{
				StatusCode: http.StatusSeeOther,
				Message:    "See Other",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}

		// Verify if user is logged in the database
		var userUUID, role string
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

		// Checking if the cookie "role" exist
		if err_cookie_role == http.ErrNoCookie {
		} else if err_cookie_role != nil {
			err := &models.CustomError{
				StatusCode: http.StatusSeeOther,
				Message:    "Error getting role (AuthMiddleWare)",
			}
			// Appel de HandleError
			HandleError(w, err.StatusCode, err.Message)
			return
		} else {

			// Checking if the User have the correct role
			state_role := `SELECT Role FROM User WHERE UUID = ?`
			err_db_role := db.QueryRow(state_role, userUUID).Scan(&role)
			if err_db_role != nil && err_db_role != sql.ErrNoRows {
				err := &models.CustomError{
					StatusCode: http.StatusSeeOther,
					Message:    "See Other",
				}
				HandleError(w, err.StatusCode, err.Message)
				return
			}
			if cookie_role.Value != role {
				lib.LogoutHandler(w, r)
			}
		}

		// Session is valid, proceed to handler
		next.ServeHTTP(w, r)
	}
}
