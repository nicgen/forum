package main

import (
	"fmt"
	"forum/cmd/RateLimit"
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
	lib.Admin()

	// Configuration du routeur
	mux := setupMux()
	rateLimiter := RateLimit.NewRateLimiter()
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

	// Authentication
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	mux.HandleFunc("/nav-register", handlers.NavRegister)
	mux.HandleFunc("/nav-login", handlers.NavLogin)

	// Tier Authentication
	mux.HandleFunc("/google", handlers.GoogleOAuthHandler)
	mux.HandleFunc("/github", handlers.GitHubOAuthHandler)
	mux.HandleFunc("/discord", handlers.DiscordOAuthHandler)
	mux.HandleFunc("/callback", handlers.GoogleCallbackHandler)
	mux.HandleFunc("/github/callback", handlers.GitHubCallbackHandler)
	mux.HandleFunc("/discord/callback", handlers.DiscordCallbackHandler)

	// DB Requests
	mux.HandleFunc("/profile", handlers.AuthMiddleware(handlers.ProfileHandler))
	mux.HandleFunc("/post", handlers.AuthMiddleware(handlers.PostHandler))
	mux.HandleFunc("/create-post", handlers.AuthMiddleware(handlers.CreatePostHandler))
	mux.HandleFunc("/comment", handlers.AuthMiddleware(handlers.CommentHandler))
	mux.HandleFunc("/like", handlers.AuthMiddleware(handlers.LikeHandler))
	mux.HandleFunc("/admin/update-role", handlers.AuthMiddleware(handlers.UpdateUserToModerator))
	mux.HandleFunc("/admin/remove-role", handlers.AuthMiddleware(handlers.RemoveModerator))
	mux.HandleFunc("/admin/delete-user", handlers.AuthMiddleware(handlers.DeleteUser))
	mux.HandleFunc("/admin/admindelete-user", handlers.AuthMiddleware(handlers.AdminDeleteUser))

	// Basic Web handlers
	mux.HandleFunc("/about", handlers.AboutHandler)
	mux.HandleFunc("/error", handlers.ForceDirectError) // !for testing purpose only (not for production)
	mux.HandleFunc("/500", handlers.Force500Handler)    // !for testing purpose only (not for production)

	return mux
}

// Configuration du serveur
func setupServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              "localhost:8080", // Écoute sur le port HTTPS
		Handler:           handlers.WithErrorHandling(handler),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
