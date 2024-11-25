package handlers

import (
	"forum/cmd/lib"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	cookie, _ := r.Cookie("session_id")

	// Getting the User Data
	data := lib.GetData(db, cookie.Value, "logged", "profile", w, r)
	data = lib.GetComments(db, cookie.Value, data, w, r)
	data = lib.GetNotifications(w, cookie.Value, data)
	data = lib.GetReport(w, data, r)
	// Redirect User to the profile html page and sending the data to it
	lib.RenderTemplate(w, "layout/index", "page/profile", data)
}

func ProfileUserHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	user_uuid := r.URL.Query().Get("uuid")

	// Getting the User Data
	data := lib.GetData(db, user_uuid, "logged", "profile_user", w, r)
	data = lib.GetComments(db, user_uuid, data, w, r)
	data = lib.GetNotifications(w, user_uuid, data)

	// Redirect User to the profile html page and sending the data to it
	lib.RenderTemplate(w, "layout/index", "page/profile_user", data)
}