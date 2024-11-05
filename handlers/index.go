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

	// Defining variables
	err := "OK"
	data := map[string]interface{}{}

	// Checking if the User is on guest or is logged
	_, err_cookie := r.Cookie("session_id")

	// If he's not logged
	if err_cookie == http.ErrNoCookie {
		data, err = lib.GetData(db, "not logged", "not logged", "index")
	} else {
		// Checking the cookie values
		session_id := r.Cookies()

		// Getting the data values
		data, err = lib.GetData(db, session_id[0].Value, "logged", "index")
		if err != "OK" {

		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "layout/index", "page/index", data)
}
