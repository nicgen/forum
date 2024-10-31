package handlers

import (
	"forum/cmd/lib"
	"log"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// ? Handler to get form values, store them into database after checking them
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()

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
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				log.Printf("Error parsing form: %v", err1)
				return
			}

			// Getting form values
			username := r.FormValue("UsernameForm")
			password := r.FormValue("PasswordForm")
			confirmPassword := r.FormValue("ConfirmPasswordForm")
			email := r.FormValue("EmailForm")

			log.Printf("Received registration request: Username=%s, Email=%s", username, email)

			// if !lib.IsValidPassword(password) {
			// 	http.Error(w, "Invalid Password Format, Your password must be at least 8 characters long and include at least one number and one special character.", http.StatusBadRequest)
			// 	return
			// }
			// if !lib.IsValidEmail(email) {
			// 	http.Error(w, "Invalid Email Format, Example of a valid email address: example@domain.com", http.StatusBadRequest)
			// 	return
			// }
			// Check if passwords match
			if password != confirmPassword {
				http.Error(w, "Password doesn't match", http.StatusBadRequest)
				return
			}

			// Generate UUID
			userUUID, errUUID := uuid.NewV4()
			if errUUID != nil {
				http.Error(w, "Error generating user UUID", http.StatusInternalServerError)
				log.Printf("Error generating UUID: %v", errUUID)
				return
			}

			var usernameExists bool
			var emailExists bool

			err_user := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE Username=?)", username).Scan(&usernameExists)
			err_email := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE Email=?)", email).Scan(&emailExists)

			if err_user != nil {
				http.Error(w, "Error checking username", http.StatusInternalServerError)
				log.Printf("Error checking username: %v", err_user)
				return
			}
			if err_email != nil {
				http.Error(w, "Error checking email", http.StatusInternalServerError)
				log.Printf("Error checking email: %v", err_email)
				return
			}

			log.Printf("Username exists: %v, Email exists: %v", usernameExists, emailExists)

			// Hash password
			hashedPassword, err_password := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err_password != nil {
				http.Error(w, "Error hashing the password", http.StatusInternalServerError)
				log.Printf("Error hashing password: %v", err_password)
				return
			}

			// Check if the user already exists
			if usernameExists {
				http.Error(w, "Username already taken", http.StatusConflict)
				log.Printf("Username already taken: %s", username)
				return
			}

			// Check if the email is already taken
			if emailExists {
				http.Error(w, "Email already used", http.StatusConflict)
				log.Printf("Email already used: %s", email)
				return
			}

			// Insert the new user into the database as a User
			_, err_db := db.Exec("INSERT INTO User (UUID, Username, Password, Email, Role) VALUES (?, ?, ?, ?, ?)", userUUID, username, hashedPassword, email, "User")
			if err_db != nil {
				log.Printf("Error adding user to the database: %v", err_db)
				http.Error(w, "Error adding user to the database", http.StatusInternalServerError)
				return
			}

			// Getting the UUID from the database
			var user_uuid string
			state := `SELECT UUID FROM User WHERE Email = ?`
			err_user = db.QueryRow(state, email).Scan(&user_uuid)
			if err_user != nil {
				http.Error(w, "Error accessing User UUID", http.StatusUnauthorized)
				return
			}

			// Attribute a session to an User
			CookieSession(user_uuid, w, r)

			// Notify server new User as been added
			log.Printf("User %s added successfully", username)

			// Redirect to a success page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	} else {
		// If a logged User tries to register without logging out
		http.Error(w, "You must log-out before loggin in again", http.StatusUnauthorized)
		return
	}

}
