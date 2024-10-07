package lib

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

// InitDB initialise la connexion à la base de données
func InitDB(dataSourceName string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture de la base de données : %v", err)
	}

	// Vérifiez la connexion
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la connexion à la base de données : %v", err)
	}

	log.Println("Connexion à la base de données établie avec succès")
	return db, nil
}

// GetDB retourne le pointeur vers la connexion à la base de données
func GetDB() *sql.DB {
	return db
}
