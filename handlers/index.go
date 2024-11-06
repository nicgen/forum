package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

// IndexHandler handles requests to the root URL
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
	db := lib.GetDB()

	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound, "Page Not Found")
		return
	}

	// Defining variables
	err := "OK"
	data := map[string]interface{}{}

	// Checking if the User is on guest or is logged
	_, err_cookie := r.Cookie("session_id")

	// If they're not logged in
	if err_cookie == http.ErrNoCookie {
		data, err = lib.GetData(db, "not logged", "not logged", "index")
	} else {
		// Get data for logged-in user
		session_id := r.Cookies()
		data, err = lib.GetData(db, session_id[0].Value, "logged", "index")
	}

	// Checking the error returned by the GetData function
	if err != "OK" {
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "layout/index", "page/index", data)
}
