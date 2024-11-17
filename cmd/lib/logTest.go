package lib

import (
	"database/sql"

	"net/http"
)

func DataTest(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	// Checking if the User is on guest or is logged
	cookie, err_cookie := r.Cookie("session_id")
	data := map[string]interface{}{}
	err_data := ""
	// If they're not logged in
	if err_cookie == http.ErrNoCookie {
		data, err_data = GetData(db, "not logged", "not logged", "index", w, r)
		if err_data != "OK" {
			ErrorServer(w, err_data)
		}
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
			data, err_data = GetData(db, cookie.Value, "logged", "index", w, r)
			if err_data != "OK" {
				ErrorServer(w, err_data)
			}
		}
	}
	return data
}
