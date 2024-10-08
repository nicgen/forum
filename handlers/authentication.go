package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
}

func TestDBConnection() {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection established successfully!")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Afficher le formulaire d'inscription
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
		userUUID := uuid.New().String()

		// Get current time
		createdAt := time.Now()

		// Set default role
		defaultRole := "user"

		var usernameExists bool
		var emailExists bool

		err_user := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&usernameExists)
		err_email := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=?)", email).Scan(&emailExists)

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
		_, err_db := db.Exec("INSERT INTO User (UUID, Username, Password, CreatedAt, Role, Email) VALUES (?, ?, ?, ?, ?, ?)", userUUID, username, hashedPassword, createdAt, defaultRole, email)
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
