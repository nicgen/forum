package main

import (
	"forum/cmd/config"
	"forum/cmd/lib"
	"log"
)

func main() {
	// Load environment variables from .env file
	lib.LoadEnv(".env")

	// Request password
	if !lib.AskForPassword() {
		log.Println("Incorrect password")
		return
	}

	// Initialize the database
	if err := lib.Init(); err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}

	// Test database connection
	lib.TestDBConnection()

	// Configure the router and server
	mux := config.SetupMux()
	rateLimiter := lib.NewRateLimiter()
	server := config.SetupServer(rateLimiter.Limit(mux))

	// Configure HTTPS
	lib.SetupHTTPS(server)

	log.Printf("Server starting on https://%s...\n", server.Addr)

	// Start HTTPS server
	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
