package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
)

// IndexHandler handles requests to the root URL
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
	db := lib.GetDB()

	if r.URL.Path != "/" {
		// * generate your error message
		err := &models.CustomError{
			StatusCode: http.StatusNotFound,
			Message:    "Page Not Found",
		}
		// Use HandleError to send the error response
		HandleError(w, err.StatusCode, err.Message)
		// return
		// * alt. use the auto-generated error code & message
		// HandleError(w, http.StatusNotFound, "Page Not Found")
		return
	}

	// Defining variables
	err := "OK"
	data := map[string]interface{}{}

	// Checking if the User is on guest or is logged
	cookie, err_cookie := r.Cookie("session_id")

	// If they're not logged in
	if err_cookie == http.ErrNoCookie {
		data, err = lib.GetData(db, "not logged", "not logged", "index")
	} else {
		var id int
		// Checking if the UUID is containned in the database
		state_check := `SELECT ID FROM User WHERE UUID = ?`
		err_check := db.QueryRow(state_check, cookie.Value).Scan(&id)

		// If the UUID is not contained in db, get rid of that cookie and redirect to homepage
		if err_check == sql.ErrNoRows {
			LogoutHandler(w, r)
		} else {
			// Else, we show the User the index page of Logged User
			data, err = lib.GetData(db, cookie.Value, "logged", "index")
		}
	}

	// Checking the error returned by the GetData function
	if err != "OK" {
		lib.ErrorServer(w, err)
	}

	data = lib.ErrorMessage(w, data, "none")

	renderTemplate(w, "layout/index", "page/index", data)
}
