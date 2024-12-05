package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

// IndexHandler handles requests to the /about URL
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := lib.DataTest(w, r)
	data = lib.ErrorMessage(w, data, "none")
	lib.RenderTemplate(w, "layout/index", "page/about", data)
}
