package lib

import (
	"net/http"
	"time"
)

// ? Handler that will delete the cookie of the User logged
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := GetDB()

	// Checking the cookie values
	cookie, _ := r.Cookie("session_id")

	// Unlogging the User in the database
	state := `UPDATE User SET IsLogged = ? WHERE UUID = ?`
	_, err_db := db.Exec(state, false, cookie.Value)
	if err_db != nil {
		ErrorServer(w, "Error logging out")
	}

	// Common cookie settings
	cookieBase := http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	}

	// Define all cookies
	cookies := []http.Cookie{
		{
			Name:  "session_id",
			Value: "",
		},
		{
			Name:  "username",
			Value: "",
		},
		{
			Name:  "creation_date",
			Value: "",
		},
		{
			Name:  "creation_hour",
			Value: "",
		},
		{
			Name:  "email",
			Value: "",
		},
		{
			Name:  "role",
			Value: "",
		},
	}

	// Set common properties and add cookies
	for _, cookie := range cookies {
		cookie.Path = cookieBase.Path
		cookie.HttpOnly = cookieBase.HttpOnly
		cookie.Secure = cookieBase.Secure
		cookie.SameSite = cookieBase.SameSite
		cookie.Expires = cookieBase.Expires
		http.SetCookie(w, &cookie)
	}

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
