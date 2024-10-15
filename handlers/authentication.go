package handlers

import (
	"log"
	"net/http"
	"web-starter/cmd/lib"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// Handler to get form values, store them into database after checking them
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()

	if r.Method == "GET" {
		http.ServeFile(w, r, "./templates/index.html")
		return
	}

	if r.Method == "POST" {
		//Parsing forms values
		err1 := r.ParseForm()
		if err1 != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Getting form values
		username := r.FormValue("UsernameForm")
		password := r.FormValue("PasswordForm")
		email := r.FormValue("EmailForm")

		// Generate UUID
		userUUID, errUUID := uuid.NewV4()
		if errUUID != nil {
			http.Error(w, "Error generating user UUID", http.StatusInternalServerError)
			return
		}

		var usernameExists bool
		var emailExists bool

		err_user := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&usernameExists)
		err_email := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=?)", email).Scan(&emailExists)
		is_valid_password := lib.IsValidPassword(password)
		is_valid_email := lib.IsValidEmail(email)

		//Checking if the email and password are valid
		if !is_valid_password {
			http.Error(w, "Password not valid", http.StatusBadRequest)
			return
		}
		if !is_valid_email {
			http.Error(w, "Email not valid", http.StatusBadRequest)
			return
		}

		// Hash password
		hashedPassword, err_password := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		// Checking for errors
		if err_user != nil {
			http.Error(w, "Error checking username", http.StatusInternalServerError)
			return
		}
		if err_email != nil {
			http.Error(w, "Error checking email", http.StatusInternalServerError)
			return
		}
		if err_password != nil {
			http.Error(w, "Error hashing the password", http.StatusInternalServerError)
			return
		}

		// Check if the user already exists
		if usernameExists {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		// Check if the email is already taken
		if emailExists {
			http.Error(w, "Email already used", http.StatusConflict)
			return
		}

		// Insérer le nouvel utilisateur dans la base de données
		_, err_db := db.Exec("INSERT INTO User (UUID, Username, Password, Email) VALUES (?, ?, ?, ?)", userUUID, username, hashedPassword, email)
		if err_db != nil {
			log.Printf("Erreur lors de l'ajout de l'utilisateur à la base de données: %v", err_db)
			http.Error(w, "Erreur lors de l'ajout de l'utilisateur à la base de données", http.StatusInternalServerError)
			return
		}

		// Rediriger vers une page de succès
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
	}
}
