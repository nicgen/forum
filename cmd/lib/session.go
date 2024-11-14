package lib

import (
	"net/http"
	"time"
)

// ? Function to attribute a temporary cookie that will store User UUID in header
func CookieSession(user_uuid, username, creation_date, creation_hour, email, role string, w http.ResponseWriter, r *http.Request) {
	// Storing Db data into a variable
	db := GetDB()

	// Common cookie settings
	cookieBase := http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(12 * time.Hour),
	}

	// Define all cookies
	cookies := []http.Cookie{
		{
			Name:  "session_id",
			Value: user_uuid,
		},
		{
			Name:  "username",
			Value: username,
		},
		{
			Name:  "creation_date",
			Value: creation_date,
		},
		{
			Name:  "creation_hour",
			Value: creation_hour,
		},
		{
			Name:  "email",
			Value: email,
		},
		{
			Name:  "role",
			Value: role,
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

	// Setting the User as Logged in the database
	state := `UPDATE User SET IsLogged = ? WHERE UUID = ?`
	_, err_db := db.Exec(state, true, user_uuid)
	if err_db != nil {
		ErrorServer(w, "Error logging in")
	}
}
