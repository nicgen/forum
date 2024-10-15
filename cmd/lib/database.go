package lib

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
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

func GetDB() *sql.DB {
	return db
}

func CreateTables() {
	tables := []string{

		`CREATE TABLE IF NOT EXISTS User(
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  UUID VARCHAR(255) NOT NULL UNIQUE,
  Email VARCHAR(50) NOT NULL UNIQUE,
  Username VARCHAR(25) NOT NULL UNIQUE,
  Password VARCHAR(100),
  IsSuperUser    BOOL, 
  IsModerator BOOL, 
  IsDeleted BOOL, 
  CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);`,

		`CREATE TABLE IF NOT EXISTS Admin (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_ID INTEGER NOT NULL,
  Mod_ID INTEGER NOT NULL,
  RequestMod_ID INTEGER NOT NULL,
  FOREIGN KEY (User_ID) REFERENCES User(ID),
  FOREIGN KEY (RequestMod_ID) REFERENCES RequestMod(ID),
  FOREIGN KEY (Mod_ID) REFERENCES Moderateur(ID)
);`,

		`Create TABLE IF NOT EXISTS Moderateur (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_ID INT NOT NULL,
  IsAdmin BOOL,
  ACCESS_GIVEN DATETIME DEFAULT CURRENT_TIMESTAMP,
  ACCESS_REVOKED DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (User_ID) REFERENCES User(ID)
);`,

		`CREATE TABLE IF NOT EXISTS Categories (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Name VARCHAR(50) UNIQUE
);`,

		`CREATE TABLE IF NOT EXISTS Posts (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_ID INTEGER NOT NULL,
  Title TEXT NOT NULL,
  Category_ID INTEGER NOT NULL,
  Texte TEXT,
  CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (User_ID) REFERENCES User(ID),
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
  User_ID INTEGER NOT NULL,
  Post_ID INTEGER NOT NULL,
  Texte TEXT NOT NULL,
  CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (Post_ID) REFERENCES Posts(ID) ON DELETE CASCADE,
  FOREIGN KEY (User_ID) REFERENCES User(ID)
);`,

		`CREATE TABLE IF NOT EXISTS Report (
  ID INTEGER PRIMARY KEY AUTOINCREMENT ,
  Reported_ID INTEGER NOT NULL, 
  User_ID INTEGER NOT NULL,
  Reported_Reason INTEGER NOT NULL,
  Reported_Texte TEXT,
  FOREIGN KEY (User_ID) REFERENCES User(ID)
);`,

		`CREATE TABLE IF NOT EXISTS RequestMod (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_ID INTEGER NOT NULL,
  Reason TEXT NOT NULL,
  FOREIGN KEY (User_ID) REFERENCES User(ID)
);`,

		`CREATE TABLE IF NOT EXISTS Reaction (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Post_ID INTEGER,
  Comment_ID INTEGER,
  User_ID INTEGER,
  IsLike BOOL,
  FOREIGN KEY (Post_ID) REFERENCES Posts(ID) ON DELETE CASCADE,
  FOREIGN KEY (Comment_ID) REFERENCES Comments(ID) ON DELETE CASCADE,
  FOREIGN KEY (User_ID) REFERENCES User(ID)
  CHECK ((Post_ID is NULL AND Comment_ID IS NOT NULL)OR(Post_ID IS NOT NULL AND Comment_ID IS NULL))
);`,

		`CREATE TABLE IF NOT EXISTS Notification (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  User_ID INTEGER NOT NULL,
  Reaction_ID INTEGER,
  Post_ID INTEGER,
  Comment_ID INTEGER,
  CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  IsRead Bool,
  FOREIGN KEY(Comment_ID) REFERENCES Comments(ID),
  FOREIGN KEY(User_ID) REFERENCES User(ID),
  FOREIGN KEY(Reaction_ID) REFERENCES Reaction(ID)
  FOREIGN KEY(Post_ID) REFERENCES Posts(ID)
);`,

		`CREATE TABLE IF NOT EXISTS Image (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  FilePath TEXT,
  Post_ID INTEGER,
  FOREIGN KEY (Post_ID) REFERENCES Posts(ID) ON DELETE CASCADE

    );`,
	}
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Tables créées avec succès.")

	InsertCategories()
}

func InsertCategories() {
	categories := []string{"Test 1", "Test 2", "Test 3"}

	for _, category := range categories {
		_, err := db.Exec(`INSERT OR IGNORE INTO Categories (Name) VALUES (?)`, category) //créer les catégories ou ignore si elles existent déjà
		if err != nil {
			log.Fatalf("Error inserting category %s: %v", category, err)
		} else {
			fmt.Printf("Catégorie '%s' insérée avec succès ou déjà existante.\n", category)
		}
	}
}
