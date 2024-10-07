package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Afficher le formulaire d'inscription
		http.ServeFile(w, r, "./templates/register.html") // Chemin vers votre formulaire d'inscription
		return
	}

	if r.Method == "POST" {
		// Récupérer les données du formulaire
		username := r.FormValue("UsernameForm")
		password := r.FormValue("PasswordForm")
		email := r.FormValue("EmailForm")

		// TODO Valider les données
		if username == "" || password == "" || email == "" {
			http.Error(w, "Les champs nom d'utilisateur, mot de passe et email sont obligatoires", http.StatusBadRequest)
			return
		}

		// Vérifier si l'utilisateur existe déjà
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&exists)
		if err != nil {
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
			http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
			return
		}

		// Insérer le nouvel utilisateur dans la base de données
		_, err = db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, string(hashedPassword), email)
		if err != nil {
			log.Printf("Erreur lors de l'ajout de l'utilisateur à la base de données: %v", err)
			http.Error(w, "Erreur lors de l'ajout de l'utilisateur à la base de données", http.StatusInternalServerError)
			return
		}

		//Vérifier si l'email est déjà pris

		// Rediriger vers une page de succès
		http.Redirect(w, r, "/success", http.StatusSeeOther)
	} else {
		http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
	}
}
