package handlers

import (
	"database/sql"
	"forum/cmd/lib"
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
			lib.ErrorServer(w, "Error preparing query")
		}
		defer stmt.Close()

		var hashedPassword string
		err = stmt.QueryRow(email).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				data, err_getdata := lib.GetData(db, "null", "notlogged", "index", r)
				if err_getdata != "OK" {
					lib.ErrorServer(w, err_getdata)
				}
				data = lib.ErrorMessage(w, data, "LoginMail")
				lib.RenderTemplate(w, "layout/index", "page/index", data)
			} else {
				lib.ErrorServer(w, "Error retrieving user data")
			}
			return
		}

		// Verify password
		if !CheckPassword(hashedPassword, password) {
			data, err_getdata := lib.GetData(db, "null", "notlogged", "index", r)
			if err_getdata != "OK" {
				lib.ErrorServer(w, err_getdata)
			}
			data = lib.ErrorMessage(w, data, "LoginPassword")
			lib.RenderTemplate(w, "layout/index", "page/index", data)
		}

		var uuid, username, creation_date, creation_hour string
		createdAt := time.Now()

		// Making the query for infos that will be stored into the cookie
		state_cookie := `SELECT UUID, Username, CreatedAt FROM User WHERE Email = ?`
		err_cookie := db.QueryRow(state_cookie, email).Scan(&uuid, &username, &createdAt)
		if err_cookie != nil {
			lib.ErrorServer(w, "Error getting User informations")
		}
		creation_date = createdAt.Format("2006-01-02") // YYYY-MM-DD
		creation_hour = createdAt.Format("15:04:05")   // HH:MM:SS

		// Attribute a session to an User
		lib.CookieSession(uuid, username, creation_date, creation_hour, w, r)

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
