package main

import (
	"fmt"
	"forum/cmd/lib"
	"forum/handlers"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Chargement des variables d'environnement
	lib.LoadEnv(".env") // variable env

	// Demande de mot de passe
	if !askForPassword() {
		log.Println("MDP incorrect")
		return
	}

	// Initialisation de la base de données
	if err := lib.Init(); err != nil {
		log.Fatalf("Erreur d'initialisation de DB: %v", err)
	}

	// Création des tables et test de la connexion
	lib.CreateTables()
	lib.TestDBConnection()

	// Configuration du routeur
	mux := setupMux()
	rateLimiter := lib.NewRateLimiter()
	server := setupServer(rateLimiter.Limit(mux))

	// Configuration HTTPS
	lib.SetupHTTPS(server) // Configure HTTPS

	log.Printf("Server starting on https://%s...\n", server.Addr)

	// lib.OpenBrowser("https://localhost:3131")

	// Démarrage du serveur HTTPS
	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil { // Utilisation de ListenAndServeTLS pour HTTPS
		log.Fatalf("Error starting server: %v", err)
	}
}

// Demande de mot de passe via bool
func askForPassword() bool {
	var password string
	fmt.Print("Entrez le MDP: ")
	fmt.Scanln(&password)

	storedPassword := os.Getenv("DB_PASSWORD")

	return password == storedPassword
}

// Configuration du routeur
func setupMux() *http.ServeMux {
	// Création d'un nouveau ServeMux
	mux := http.NewServeMux()

	// Servir des fichiers statiques
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up routes
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/google", handlers.GoogleOAuthHandler)
	mux.HandleFunc("/github", handlers.GitHubOAuthHandler)
	mux.HandleFunc("/discord", handlers.DiscordOAuthHandler)
	mux.HandleFunc("/callback", handlers.GoogleCallbackHandler)
	mux.HandleFunc("/github/callback", handlers.GitHubCallbackHandler)
	mux.HandleFunc("/discord/callback", handlers.DiscordCallbackHandler)
	mux.HandleFunc("/profile", handlers.ProfileHandler)
	mux.HandleFunc("/about", handlers.AboutHandler)
	mux.HandleFunc("/error", handlers.ForceDirectError) // !for testing purpose only (not for production)
	mux.HandleFunc("/500", handlers.Force500Handler)    // !for testing purpose only (not for production)

	return mux
}

// Configuration du serveur
func setupServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              "localhost:8080", // Écoute sur le port HTTPS
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
