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
		err1 := r.ParseForm()
		if err1 != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Récupérer les données du formulaire
		username := r.FormValue("UsernameForm")
		password := r.FormValue("PasswordForm")
		email := r.FormValue("EmailForm")

		// Valider les données
		if username == "" || password == "" || email == "" {
			http.Error(w, "Les champs nom d'utilisateur, mot de passe et email sont obligatoires", http.StatusBadRequest)
			return
		}

		// Generate UUID
		userUUID := uuid.New().String()

		// Get current time
		createdAt := time.Now()

		// Set default role
		defaultRole := "user"

		// Vérifier si l'utilisateur existe déjà
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE Username=?)", username).Scan(&exists)
		if err != nil {
			log.Printf("Erreur lors de la vérification de l'utilisateur: %v", err)
			http.Error(w, "Erreur lors de la vérification de l'utilisateur", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Nom d'utilisateur déjà pris", http.StatusConflict)
			return
		}

		// Hash du mot de passe
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Erreur lors du hash du mot de passe: %v", err)
			http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
			return
		}

		// Insérer le nouvel utilisateur dans la base de données
		_, err = db.Exec("INSERT INTO User (UUID, Username, Password, CreatedAt, Role, Email) VALUES (?, ?, ?, ?, ?, ?)", userUUID, username, hashedPassword, createdAt, defaultRole, email)
		if err != nil {
			log.Printf("Erreur lors de l'ajout de l'utilisateur à la base de données: %v", err)
			http.Error(w, "Erreur lors de l'ajout de l'utilisateur à la base de données", http.StatusInternalServerError)
			return
		}

		// Rediriger vers une page de succès
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
	}
}
