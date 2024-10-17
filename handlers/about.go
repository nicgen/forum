package handlers

import (
	"forum/models"
	"net/http"
)

// IndexHandler handles requests to the /about URL
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		Title:  "About",
		Header: "About Us",
		Content: map[string]interface{}{
			"Message": "Learn more about our company and mission.",
		},
		IsError: false,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "about", data)
}
