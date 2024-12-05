package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strings"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {

	db := lib.GetDB()

	if r.Method == http.MethodPost {

		categoryName := strings.TrimSpace(r.FormValue("categoryName"))

		if categoryName == "" {
			//Erreur critique : nom de catégorie ne peut pas etre vide
			err := &models.CustomError{
				StatusCode: http.StatusBadRequest,
				Message:    "Veuillez remplir le champ 'nom de catégorie'",
			}
			HandleError(w, err.StatusCode, err.Message)
		}

		insertQuery := `INSERT INTO Categories (Name) VALUES (?)`
		_, err := db.Exec(insertQuery, categoryName)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				//Erreur critique : Category already exist
				err := &models.CustomError{
					StatusCode: http.StatusConflict,
					Message:    "Category Already exist",
				}
				HandleError(w, err.StatusCode, err.Message)
			} else {
				// Erreur non critique : Echec de la création de la catégorie
				lib.ErrorServer(w, "Échec de la création de la catégorie, veuillez réessayer plus tard.")
			}
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
}
