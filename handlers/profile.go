package handlers

import (
	"forum/cmd/lib"
	"net/http"
	"strings"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	session_id := r.Cookies()

	// Setting up the variables of the User informations
	var user_username, user_email, user_creation, user_role string
	var role bool

	// Setting up the states before for the database requests
	state_uuid := `SELECT Username FROM User WHERE UUID = ?`
	state_email := `SELECT Email FROM User WHERE UUID = ?`
	state_creation := `SELECT CreatedAt FROM User WHERE UUID = ?`
	state_role := `SELECT IsModerator FROM User WHERE UUID = ?`

	// Making requests to the database
	err_uuid := db.QueryRow(state_uuid, session_id[0].Value).Scan(&user_username)
	err_email := db.QueryRow(state_email, session_id[0].Value).Scan(&user_email)
	err_creation := db.QueryRow(state_creation, session_id[0].Value).Scan(&user_creation)
	err_role := db.QueryRow(state_role, session_id[0].Value).Scan(&role)

	// Checking for database requests errors
	if err_uuid != nil {
		http.Error(w, "Error accessing User UUID", http.StatusUnauthorized)
		return
	} else if err_email != nil {
		http.Error(w, "Error accessing User EMAIL", http.StatusUnauthorized)
		return
	} else if err_creation != nil {
		http.Error(w, "Error accessing User CREATION DATE", http.StatusUnauthorized)
		return
	} else if err_role != nil {
		http.Error(w, "Error accessing User ROLE", http.StatusUnauthorized)
		return
	}

	// Checking if the User is a moderator or not
	if !role {
		user_role = "Member"
	} else {
		user_role = "Moderator"
	}

	// Spliting the creation date into 2 different values
	creation := strings.Split(user_creation, "T")
	creation_Date := creation[0]
	creation_Hour := creation[1][:len(creation[1])-1]

	// Printing the User informations
	println("Username: ", user_username)
	println("Email: ", user_email)
	println("Creation date: ", creation_Date)
	println("Creation hour: ", creation_Hour)
	println("Role: ", user_role)

	// Storing the data into a map that can be sent into the html
	data := map[string]interface{}{
		"Username":     user_username,
		"Email":        user_email,
		"CreationDate": creation_Date,
		"CreationHour": creation_Hour,
		"Role":         user_role,
	}

	renderTemplate(w, "layout/index", "page/profile", data)
}