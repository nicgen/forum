package handlers

import (
	"net/http"
	"time"
)

// ? Handler that will delete the cookie of the User logged
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
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
