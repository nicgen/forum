package lib

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // or your database driver
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
