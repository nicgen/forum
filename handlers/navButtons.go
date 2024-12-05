package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

// ? Function to show the login block on the index page
func NavLogin(w http.ResponseWriter, r *http.Request) {
	// Getting database data into variable
	db := lib.GetDB()

	// Getting the data map
	data := lib.GetData(db, "empty", "not logged", "index", w, r)
	data["NavLogin"] = "show"

	data = lib.ErrorMessage(w, data, "none")
	// Redirect User to the home page
	lib.RenderTemplate(w, "layout/index", "page/index", data)
}

// ? Function to show the register block on the index page
func NavRegister(w http.ResponseWriter, r *http.Request) {
	// Getting database data into variable
	db := lib.GetDB()

	// Getting the data map
	data := lib.GetData(db, "empty", "not logged", "index", w, r)
	data["NavRegister"] = "show"

	data = lib.ErrorMessage(w, data, "none")
	// Redirect User to the home page
	lib.RenderTemplate(w, "layout/index", "page/index", data)
}
