package handlers

import (
	"fmt"
	"forum/cmd/lib"
	"net/http"
	"strings"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {

	db := lib.GetDB()

	if r.Method == http.MethodPost {

		categoryName := strings.TrimSpace(r.FormValue("categoryName"))

		if categoryName == "" {
			http.Error(w, "Category name cannot be empty", http.StatusBadRequest)
			return
		}

		insertQuery := `INSERT INTO Categories (Name) VALUES (?)`
		_, err := db.Exec(insertQuery, categoryName)
		if err != nil {

			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				http.Error(w, "Category already exists", http.StatusConflict)
			} else {
				http.Error(w, "Failed to create category", http.StatusInternalServerError)
				fmt.Println("Database error:", err)
			}
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
}
