package handlers

import (
	"forum/cmd/lib"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// ? Handler to get form values, store them into database after checking them
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
	db := lib.GetDB()

	// Setting up the variables we are going to set into cookies
	var username, creation_date, creation_hour string
	actual_time := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	creation_date = actual_time[0]
	creation_hour = actual_time[1]

	if r.Method == "GET" {
		http.ServeFile(w, r, "./templates/index.html")
		return
	}
	if r.Method == "POST" {
		// Checking if the User is already logged or not
		_, err_cookie := r.Cookie("session_id")
		if err_cookie == http.ErrNoCookie {
			// Parsing form values
			err1 := r.ParseForm()
			if err1 != nil {
				lib.ErrorServer(w, "Error parsing form data")
			}

			// Getting form values
			username = r.FormValue("UsernameForm")
			password := r.FormValue("PasswordForm")
			confirmPassword := r.FormValue("ConfirmPasswordForm")
			email := r.FormValue("EmailForm")

			log.Printf("Received registration request: Username=%s, Email=%s", username, email)

			// if !lib.IsValidPassword(password) {
			// 	data, err_getdata := lib.GetData(db, "null", "notlogged", "index")
			// 	if err_getdata != "OK" {
			// 		lib.ErrorServer(w, err_getdata)
			// 	}
			// 	data = lib.ErrorMessage(w, data, "RegisterPassword")
			// data["NavRegister"] = "show"
			// 	lib.RenderTemplate(w, "layout/index", "page/index", data)
			// }
			// if !lib.IsValidEmail(email) {
			// 	data, err_getdata := lib.GetData(db, "null", "notlogged", "index")
			// 	if err_getdata != "OK" {
			// 		lib.ErrorServer(w, err_getdata)
			// 	}
			// 	data = lib.ErrorMessage(w, data, "EmailFormat")
			// data["NavRegister"] = "show"
			// 	lib.RenderTemplate(w, "layout/index", "page/index", data)
			// }
			// Check if passwords match
			if password != confirmPassword {
				data, err_getdata := lib.GetData(db, "null", "notlogged", "index", r)
				if err_getdata != "OK" {
					lib.ErrorServer(w, err_getdata)
				}
				data = lib.ErrorMessage(w, data, "PasswordMatch")
				data["NavRegister"] = "show"
				lib.RenderTemplate(w, "layout/index", "page/index", data)
			}

			// Generate UUID
			userUUID, errUUID := uuid.NewV4()
			if errUUID != nil {
				lib.ErrorServer(w, "Error generating user UUID on register")
			}

			var usernameExists bool
			var emailExists bool

			err_user := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE Username=?)", username).Scan(&usernameExists)
			err_email := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE Email=?)", email).Scan(&emailExists)

			if err_user != nil {
				lib.ErrorServer(w, "Error checking username in database")
			} else if err_email != nil {
				lib.ErrorServer(w, "Error checking email in database")
			}

			log.Printf("Username exists: %v, Email exists: %v", usernameExists, emailExists)

			// Hash password
			hashedPassword, err_password := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err_password != nil {
				lib.ErrorServer(w, "Error hashing the password")
			}

			// Check if the user already exists
			if usernameExists {
				data, err_getdata := lib.GetData(db, "null", "notlogged", "index", r)
				if err_getdata != "OK" {
					lib.ErrorServer(w, err_getdata)
				}
				data = lib.ErrorMessage(w, data, "RegisterUsername")
				data["NavRegister"] = "show"
				lib.RenderTemplate(w, "layout/index", "page/index", data)
			}

			// Check if the email is already taken
			if emailExists {
				data, err_getdata := lib.GetData(db, "null", "notlogged", "index", r)
				if err_getdata != "OK" {
					lib.ErrorServer(w, err_getdata)
				}
				data = lib.ErrorMessage(w, data, "RegisterEmail")
				data["NavRegister"] = "show"
				lib.RenderTemplate(w, "layout/index", "page/index", data)
			}

			// Insert the new user into the database as a User
			_, err_db := db.Exec("INSERT INTO User (UUID, Username, Password, Email, Role) VALUES (?, ?, ?, ?, ?)", userUUID, username, hashedPassword, email, "User")
			if err_db != nil {
				lib.ErrorServer(w, "Error adding user to the database")
			}

			// Getting the UUID from the database
			var user_uuid string
			state := `SELECT UUID FROM User WHERE Email = ?`
			err_user = db.QueryRow(state, email).Scan(&user_uuid)
			if err_user != nil {
				lib.ErrorServer(w, "Error accessing User UUID")
			}

			// Attribute a session to an User
			lib.CookieSession(user_uuid, username, creation_date, creation_hour, email, "User", w, r)

			// Notify server new User as been added
			log.Printf("User %s added successfully", username)

			// Redirect to a success page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			// If a logged User tries to register without logging out
			lib.ErrorServer(w, "You must log-out before loggin in again")
		}
	} else {
		lib.ErrorServer(w, "Unsupported method")
	}
}
