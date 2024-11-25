package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"time"

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
			//Erreur critique : échec de la préparation de la rêquete
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error preparing query. Please try again later.",
			}
			HandleError(w, err.StatusCode, err.Message)
		}
		defer stmt.Close()

		var hashedPassword string
		err = stmt.QueryRow(email).Scan(&hashedPassword)
		if err != nil {
			//Erreur non critique : Email
			if err == sql.ErrNoRows {
				data := lib.GetData(db, "null", "notlogged", "index", w, r)
				data = lib.ErrorMessage(w, data, "LoginMail")
				data["NavLogin"] = "show"
				lib.RenderTemplate(w, "layout/index", "page/index", data)
			} else {
				//Erreur critique : echec de la recuperation des données utilisateur
				err := &models.CustomError{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error retrieving user data. Please try again later.",
				}
				HandleError(w, err.StatusCode, err.Message)
			}
			return
		}

		// Verify password
		if !CheckPassword(hashedPassword, password) {
			//Erreur non critique : MDP
			data := lib.GetData(db, "null", "notlogged", "index", w, r)
			data = lib.ErrorMessage(w, data, "LoginPassword")
			data["NavLogin"] = "show"
			lib.RenderTemplate(w, "layout/index", "page/index", data)
		}

		var uuid, username, creation_date, creation_hour, role string
		createdAt := time.Now()

		// Making the query for infos that will be stored into the cookie
		state_cookie := `SELECT UUID, Username, CreatedAt, Role FROM User WHERE Email = ?`
		err_cookie := db.QueryRow(state_cookie, email).Scan(&uuid, &username, &createdAt, &role)
		if err_cookie != nil {
			//Erreur critique: echec de la recuperation des informations utilisateur.
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error getting User information. Please try again later.",
			}
			HandleError(w, err.StatusCode, err.Message)
		}
		creation_date = createdAt.Format("2006-01-02") // YYYY-MM-DD
		creation_hour = createdAt.Format("15:04:05")   // HH:MM:SS

		// Attribute a session to an User
		lib.CookieSession(uuid, username, creation_date, creation_hour, email, role, w, r)

	} else {
		// If the User is already logged and tries to log-in
		lib.ErrorServer(w, "You must log-out before loggin in again")
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
