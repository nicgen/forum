package handlers

import (
	"forum/cmd/lib"
	"net/http"
	"time"
)

// ? Handler that will delete the cookie of the User logged
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := lib.GetDB()

	// Checking the cookie values
	cookie, _ := r.Cookie("session_id")

	// Unlogging the User in the database
	state := `UPDATE User SET IsLogged = ? WHERE UUID = ?`
	_, err_db := db.Exec(state, false, cookie.Value)
	if err_db != nil {
		ErrorServer(w, "Error logging out")
	}

	// Overwrite the cookie with one that expire instantly
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
	})

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
