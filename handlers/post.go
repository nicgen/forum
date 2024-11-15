package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	data := lib.DataTest(w, r)
	data = lib.ErrorMessage(w, data, "none")

	lib.RenderTemplate(w, "layout/index", "page/index", data)
}
