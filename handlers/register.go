package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"log"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// ? Handler to get form values, store them into database after checking them
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
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
				// Erreur non critique : Problème lors de l'analyse des données du formulaire
				lib.ErrorServer(w, "Error parsing form data")
				return
			}

			// Getting form values
			username := r.FormValue("UsernameForm")
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
			// 	lib.RenderTemplate(w, "layout/index", "page/index", data)
			// }
			// if !lib.IsValidEmail(email) {
			// 	data, err_getdata := lib.GetData(db, "null", "notlogged", "index")
			// 	if err_getdata != "OK" {
			// 		lib.ErrorServer(w, err_getdata)
			// 	}
			// 	data = lib.ErrorMessage(w, data, "EmailFormat")
			// 	lib.RenderTemplate(w, "layout/index", "page/index", data)
			// }
			// Check if passwords match
			if password != confirmPassword {
				data, err_getdata := lib.GetData(db, "null", "notlogged", "index")
				if err_getdata != "OK" {
					lib.ErrorServer(w, err_getdata)
				}
				data = lib.ErrorMessage(w, data, "PasswordMatch")
				lib.RenderTemplate(w, "layout/index", "page/index", data)
			}

			// Generate UUID
			userUUID, errUUID := uuid.NewV4()
			if errUUID != nil {
				// Erreur critique : Problème de génération de l'UUID de l'utilisateur
				err := &models.CustomError{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error generating user UUID on register",
				}
				HandleError(w, err.StatusCode, err.Message)
				return
			}

			var usernameExists bool
			var emailExists bool

			err_user := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE Username=?)", username).Scan(&usernameExists)
			err_email := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE Email=?)", email).Scan(&emailExists)

			if err_user != nil {
				// Erreur critique : Problème lors de la vérification du nom d'utilisateur dans la base de données
				err := &models.CustomError{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error checking username in database",
				}
				HandleError(w, err.StatusCode, err.Message)
				return
			}
			if err_email != nil {
				// Erreur critique : Problème lors de la vérification de l'email dans la base de données
				err := &models.CustomError{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error checking email in database",
				}
				HandleError(w, err.StatusCode, err.Message)
				return
			}

			log.Printf("Username exists: %v, Email exists: %v", usernameExists, emailExists)

			// Hash password
			hashedPassword, err_password := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err_password != nil {
				// Erreur critique : Problème de hachage du mot de passe
				err := &models.CustomError{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error hashing the password",
				}
				HandleError(w, err.StatusCode, err.Message)
				return
			}
			// Check if the user already exists
			if usernameExists {
				data, err_getdata := lib.GetData(db, "null", "notlogged", "index")
				if err_getdata != "OK" {
					lib.ErrorServer(w, err_getdata)
					return
				}
				data = lib.ErrorMessage(w, data, "RegisterUsername")
				lib.RenderTemplate(w, "layout/index", "page/index", data)
			}

			// Check if the email is already taken
			if emailExists {
				data, err_getdata := lib.GetData(db, "null", "notlogged", "index")
				if err_getdata != "OK" {
					lib.ErrorServer(w, err_getdata)
					return
				}
				data = lib.ErrorMessage(w, data, "RegisterEmail")
				lib.RenderTemplate(w, "layout/index", "page/index", data)
			}

			// Insert the new user into the database as a User
			_, err_db := db.Exec("INSERT INTO User (UUID, Username, Password, Email, Role) VALUES (?, ?, ?, ?, ?)", userUUID, username, hashedPassword, email, "User")
			if err_db != nil {
				// Erreur critique : Problème lors de l'ajout de l'utilisateur dans la base de données
				err := &models.CustomError{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error adding user to the database",
				}
				HandleError(w, err.StatusCode, err.Message)
				return
			}

			// Getting the UUID from the database
			var user_uuid string
			state := `SELECT UUID FROM User WHERE Email = ?`
			err_user = db.QueryRow(state, email).Scan(&user_uuid)
			if err_user != nil {
				// Erreur non critique : Problème d'accès à l'UUID de l'utilisateur
				lib.ErrorServer(w, "Error accessing User UUID")
				return
			}

			// Attribute a session to an User
			lib.CookieSession(user_uuid, w, r)

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
