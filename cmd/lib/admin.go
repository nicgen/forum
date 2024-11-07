package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func Admin() {
	LoadEnv(".env")

	// Retrieve the admin password from environment variables
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD is not defined in the .env file")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing the password: %v", err)
	}

	// Check if the admin user exists and insert if necessary
	createAdminUserIfNotExists(db, string(hashedPassword))
}

func createAdminUserIfNotExists(db *sql.DB, hashedPassword string) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE Username = ?)", "admin").Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking for admin user: %v", err)
	}

	if !exists {
		// Example UUID and email for the admin user
		userUUID := "123e4567-e89b-12d3-a456-426614174000"
		email := "admin@admin.com"
		username := "admin"

		_, err := db.Exec(
			"INSERT INTO User (UUID, Username, Password, Email, Role) VALUES (?, ?, ?, ?, ?)",
			userUUID, username, hashedPassword, email, "Admin",
		)
		if err != nil {
			log.Fatalf("Error inserting admin user: %v", err)
		}
		fmt.Println("Admin user created successfully.")
	} else {
		fmt.Println("Admin user already exists.")
	}
}
