package handlers

import (
	"forum/models"
	"html/template"
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
		Title:  "Forum",
		Header: "Welcome to our Forum project.",
		Content: map[string]template.HTML{
			"Msg_raw":    "<h2>Sub-Title 02.</h1><p>paragraph</>",
			"Msg_styled": "<h1 style=\"text=color: blue;\">Title 01.</h1><p>paragraph with</br>style</>",
		},

		IsError: false,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "layout/index", "page/index", data)
}
