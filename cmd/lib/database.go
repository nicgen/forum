package lib

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Variable that will store the database
var db *sql.DB

// ? Initiate the database
func Init() error {
	var err error

	//Ouverture de la connexion DB
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("Failed to open database: %w", err)
	}

	// Set DB password if needed
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return fmt.Errorf("DB_PASSWORD is not set")
	}

	_, err = db.Exec(fmt.Sprintf("PRAGMA key = '%s'", dbPassword))
	if err != nil {
		return fmt.Errorf("failed to set database password: %w", err)
	}

	// Load and execute SQL schema
	if err := executeSQLFile("../db/schema"); err != nil {
		return fmt.Errorf("failed to execute SQL file: %w", err)
	}

	return nil
}

// Function to execute an external SQL file
func executeSQLFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Split content by semicolons to execute each statement individually
	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query != "" {
			_, err := db.Exec(query)
			if err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}
		}
	}

	return nil
}

// Test database connection
func TestDBConnection() {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection established successfully!")
}

// Send database instance
func GetDB() *sql.DB {
	return db
}

// Insert default categories
func InsertCategories() {
	categories := []string{"Test 1", "Test 2", "Test 3"}
	for _, category := range categories {
		_, err := db.Exec(`INSERT OR IGNORE INTO Categories (Name) VALUES (?)`, category)
		if err != nil {
			log.Fatalf("Error inserting category %s: %v", category, err)
		} else {
			fmt.Printf("Category '%s' inserted successfully or already exists.\n", category)
		}
	}
}

// Anonymizes user data upon account deletion request
func AnonymizeUser(uuid string) error {
	// Prepare anonymization queries
	_, err := db.Exec(`UPDATE User SET Username = 'anonymous', Email = 'anonymous@domain.com', Password = '', IsDeleted = TRUE WHERE UUID = ?`, uuid)
	if err != nil {
		return fmt.Errorf("failed to anonymize user data: %w", err)
	}

	// Anonymize posts, comments, reactions
	_, err = db.Exec(`UPDATE Posts SET User_UUID = 'anonymous' WHERE User_UUID = ?`, uuid)
	if err != nil {
		return fmt.Errorf("failed to anonymize user posts: %w", err)
	}

	_, err = db.Exec(`UPDATE Comments SET User_UUID = 'anonymous' WHERE User_UUID = ?`, uuid)
	if err != nil {
		return fmt.Errorf("failed to anonymize user comments: %w", err)
	}

	_, err = db.Exec(`UPDATE Reaction SET User_ID = NULL WHERE User_ID = (SELECT ID FROM User WHERE UUID = ?)`, uuid)
	if err != nil {
		return fmt.Errorf("failed to anonymize user reactions: %w", err)
	}

	return nil
}
