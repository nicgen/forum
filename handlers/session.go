package handlers

import (
	"net/http"
)

func CookieSession(user_uuid string, w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "session_id", // Name of the cookie
		Value: user_uuid,    // Using UUID as session token
		// Path:     "/",          // Cookie is valid for all paths
		// HttpOnly: true,         // Cannot be accessed by JavaScript
		// Secure:   true,         // Only sent over HTTPS
		// SameSite: http.SameSiteStrictMode,
		// // Expires in 24 hours
		// Expires: time.Now().Add(24 * time.Hour),
	}
	session_id := r.Cookies()

	// Set the cookie in the response header
	http.SetCookie(w, cookie)

	println("-----------------------------")
	for i := 0; i < len(session_id); i++ {
		println("Session ID: ", session_id[i].Value)
	}
	println("-----------------------------")
}
