package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// Err crit : cookie de session non trouvé
		err := &models.CustomError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Session not found, please log in again",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Getting the User Data
	data := lib.GetData(db, cookie.Value, "logged", "profile", w, r)
	if data == nil {
		//Err non crit : Impossible de recuperer les données.
		lib.ErrorServer(w, "Unable to retrieve user data, please try again")
	}

	data = lib.GetComments(db, cookie.Value, data, w, r)
	if data == nil {
		//Err non crit : Impossible de recup les commentaires
		lib.ErrorServer(w, "Unable to retrieve comments, please try again later.")
	}
	data = lib.GetNotifications(w, cookie.Value, data)
	data = lib.GetReport(w, data, r)
	if data == nil {
		//Err non critique : Impossible de recup les notifications
		lib.ErrorServer(w, "Unable to retrieve notifications, please try again later")
	}

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
	// data = lib.GetNotifications(w, user_uuid, data)

	// Redirect User to the profile html page and sending the data to it
	lib.RenderTemplate(w, "layout/index", "page/profile", data)
}
