package lib

import (
	"forum/models"
	"net/http"
)

// ? Function to get categories from the database and store it into the data map
func GetCategories(w http.ResponseWriter, data map[string]interface{}) map[string]interface{} {
	// Getting the Categories from the database
	var categories []*models.Category
	state_categories := `SELECT ID, Name FROM Categories`
	query_category, err_liked := db.Query(state_categories)
	if err_liked != nil {
		ErrorServer(w, "Error accessing Categories")
	}
	defer query_category.Close()

	// Ranging over categories
	for query_category.Next() {
		var category models.Category
		if err := query_category.Scan(&category.ID, &category.Name); err != nil {
			ErrorServer(w, "Error scanning Categories")
		}

		categories = append(categories, &category)
	}
	data["Categories"] = categories
	return data
}
