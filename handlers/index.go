package handlers

import (
	"forum/models"
	"net/http"
)

// IndexHandler handles requests to the root URL
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// * generate your error message
		// err := &models.CustomError{
		// 	StatusCode: http.StatusNotFound,
		// 	Message:    "Page Not Found",
		// }
		// Use HandleError to send the error response
		// HandleError(w, err.StatusCode, err.Message)
		// return
		// * alt. use the auto-generated error code & message
		HandleError(w, http.StatusNotFound, "Page Not Found")
		return
	}

	data := models.PageData{
		Title:  "Home",
		Header: "Welcome to My Web App",
		Content: map[string]interface{}{
			"Message":   "This is the home page.",
			"Paragraph": "This is a new paragraph",
		},
		IsError: false,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "index", data)
}
