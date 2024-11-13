package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"net/http"
)

// IndexHandler handles requests to the root URL
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// * generate your error message
		err := &models.CustomError{
			StatusCode: http.StatusNotFound,
			Message:    "Page Not Found",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	data := lib.DataTest(w, r)
	data = lib.ErrorMessage(w, data, "none")

	lib.RenderTemplate(w, "layout/index", "page/index", data)
}
