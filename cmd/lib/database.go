package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Variable that will store the database
var db *sql.DB

// Initiate the database
func Init() error {
	var err error

	// Open the DB connection
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("Failed to open database: %w", err)
	}

	// Set database password
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return fmt.Errorf("DB_PASSWORD is not set")
	}

	// Define the database password
	_, err = db.Exec(fmt.Sprintf("PRAGMA key = '%s'", dbPassword))
	if err != nil {
		return fmt.Errorf("failed to set database password: %w", err)
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

// Send database into handlers
func GetDB() *sql.DB {
	return db
}

// Create the database
func CreateTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS User(
   ID INTEGER PRIMARY KEY AUTOINCREMENT,
   UUID VARCHAR(255) NOT NULL UNIQUE,
   Email VARCHAR(50) NOT NULL UNIQUE,
   Username VARCHAR(25) NOT NULL UNIQUE,
   Password VARCHAR(100),
   OAuthID VARCHAR(255) UNIQUE,
   Role TEXT NOT NULL CHECK (Role IN ('Admin', 'User', 'Moderator', 'DeleteUser')),
   IsLogged BOOL DEFAULT FALSE,
   IsDeleted BOOL DEFAULT FALSE,
   IsRequest BOOL DEFAULT FALSE,
   CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );`,

		// 		`CREATE TABLE IF NOT EXISTS Admin (
		//   ID INTEGER PRIMARY KEY AUTOINCREMENT,
		//   User_UUID VARCHAR(255) NOT NULL,
		//   Mod_ID INTEGER NOT NULL,
		//   RequestMod_ID INTEGER NOT NULL,
		//   FOREIGN KEY (User_UUID) REFERENCES User(UUID),
		//   FOREIGN KEY (RequestMod_ID) REFERENCES RequestMod(ID),
		//   FOREIGN KEY (Mod_ID) REFERENCES Moderator(ID)
		// );`,

		// 		`CREATE TABLE IF NOT EXISTS Moderator (
		//   ID INTEGER PRIMARY KEY AUTOINCREMENT,
		//   User_UUID VARCHAR(255) NOT NULL,
		//   IsAdmin BOOL DEFAULT FALSE,
		//   ACCESS_GIVEN DATETIME DEFAULT CURRENT_TIMESTAMP,
		//   ACCESS_REVOKED DATETIME DEFAULT CURRENT_TIMESTAMP,
		//   FOREIGN KEY (User_UUID) REFERENCES User(UUID)
		// );`,

		`CREATE TABLE IF NOT EXISTS Categories (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Name VARCHAR(50) UNIQUE
);`,

		`CREATE TABLE IF NOT EXISTS Posts (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_UUID VARCHAR(255) NOT NULL,
  Title TEXT NOT NULL,
  Category_ID INTEGER NOT NULL,
  Text TEXT,
  Like INTEGER,
  Dislike INTEGER,
  CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (User_UUID) REFERENCES User(UUID),
  FOREIGN KEY (Category_ID) REFERENCES Categories(ID)
);`,

		`CREATE TABLE IF NOT EXISTS Post_Categories (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Post_ID INTEGER NOT NULL, 
  Categories_ID INTEGER NOT NULL,
  FOREIGN KEY(Post_ID) REFERENCES Posts(ID),
  FOREIGN KEY(Categories_ID) REFERENCES Categories(ID)
);`,

		`CREATE TABLE IF NOT EXISTS Comments (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_UUID VARCHAR(255) NOT NULL,
  Post_ID INTEGER NOT NULL,
  Text TEXT NOT NULL,
  Like INTEGER,
  Dislike INTEGER,
  CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (Post_ID) REFERENCES Posts(ID) ON DELETE CASCADE,
  FOREIGN KEY (User_UUID) REFERENCES User(UUID)
);`,

		`CREATE TABLE IF NOT EXISTS Report (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_UUID VARCHAR(255) NOT NULL,
  Username VARCHAR(255) NOT NULL,
  Post_ID INTEGER NOT NULL,
  Title TEXT NOT NULL,
  Respons_Text TEXT,
  FOREIGN KEY (User_UUID) REFERENCES User(UUID)
);`,

		// 		`CREATE TABLE IF NOT EXISTS RequestMod (
		//   ID INTEGER PRIMARY KEY AUTOINCREMENT,
		//   User_UUID VARCHAR(255) NOT NULL,
		//   Reason TEXT NOT NULL,
		//   FOREIGN KEY (User_UUID) REFERENCES User(UUID)
		// );`,

		`CREATE TABLE IF NOT EXISTS Reaction (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Post_ID INTEGER,
  Comment_ID INTEGER,
  User_UUID VARCHAR(255) NOT NULL,
  Status VARCHAR(255) NOT NULL,
  FOREIGN KEY (Post_ID) REFERENCES Posts(ID) ON DELETE CASCADE,
  FOREIGN KEY (Comment_ID) REFERENCES Comments(ID) ON DELETE CASCADE,
  FOREIGN KEY (User_UUID) REFERENCES User(UUID),
  CHECK ((Post_ID IS NULL AND Comment_ID IS NOT NULL) OR (Post_ID IS NOT NULL AND Comment_ID IS NULL))
);`,

		`CREATE TABLE IF NOT EXISTS Notification (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_UUID VARCHAR(255) NOT NULL,
  Reaction_ID INTEGER,
  Post_ID INTEGER,
  Comment_ID INTEGER,
  CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  IsRead BOOL,
  FOREIGN KEY (Comment_ID) REFERENCES Comments(ID),
  FOREIGN KEY (User_UUID) REFERENCES User(UUID),
  FOREIGN KEY (Reaction_ID) REFERENCES Reaction(ID),
  FOREIGN KEY (Post_ID) REFERENCES Posts(ID)
);`,

		`CREATE TABLE IF NOT EXISTS Image (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  FilePath TEXT,
  Post_ID INTEGER,
  FOREIGN KEY (Post_ID) REFERENCES Posts(ID) ON DELETE CASCADE
);`,

		`CREATE TABLE IF NOT EXISTS oauth_states (
    state TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`,
	}
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Tables created successfully.")

}
