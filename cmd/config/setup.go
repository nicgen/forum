// config/setup.go
package config

import (
	"forum/handlers" // Assuming handlers is your package for route handlers
	"net/http"
	"time"
)

// SetupMux configures the HTTP router for the application
func SetupMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up routes
	mux.HandleFunc("/", handlers.IndexHandler)

	// Authentication
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	// OAuth Authentication
	mux.HandleFunc("/google", handlers.GoogleOAuthHandler)
	mux.HandleFunc("/github", handlers.GitHubOAuthHandler)
	mux.HandleFunc("/discord", handlers.DiscordOAuthHandler)
	mux.HandleFunc("/callback", handlers.GoogleCallbackHandler)
	mux.HandleFunc("/github/callback", handlers.GitHubCallbackHandler)
	mux.HandleFunc("/discord/callback", handlers.DiscordCallbackHandler)

	// DB Requests
	mux.HandleFunc("/profile", handlers.AuthMiddleware(handlers.ProfileHandler))
	mux.HandleFunc("/post", handlers.AuthMiddleware(handlers.PostHandler))

	// Basic Web handlers
	mux.HandleFunc("/about", handlers.AboutHandler)
	mux.HandleFunc("/error", handlers.ForceDirectError) // For testing purpose only (not for production)
	mux.HandleFunc("/500", handlers.Force500Handler)    // For testing purpose only (not for production)

	return mux
}

// SetupServer configures the HTTP server with custom settings
func SetupServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              "localhost:8080", // Listen on HTTPS port
		Handler:           handlers.WithErrorHandling(handler),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
