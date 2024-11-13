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
			lib.ErrorServer(w, "Error preparing query")
		}
		defer stmt.Close()

		var hashedPassword string
		err = stmt.QueryRow(email).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				data, err_getdata := lib.GetData(db, "null", "notlogged", "index")
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
			data, err_getdata := lib.GetData(db, "null", "notlogged", "index")
			if err_getdata != "OK" {
				lib.ErrorServer(w, err_getdata)
			}
			data = lib.ErrorMessage(w, data, "LoginPassword")
			lib.RenderTemplate(w, "layout/index", "page/index", data)
		}

		// Getting the UUID from the database
		var user_uuid string
		state := `SELECT UUID FROM User WHERE Email = ?`
		err_user := db.QueryRow(state, email).Scan(&user_uuid)
		if err_user != nil {
			lib.ErrorServer(w, "Error accessing User UUID in the database")
		}

		// Attribute a session to an User
		lib.CookieSession(user_uuid, w, r)

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
