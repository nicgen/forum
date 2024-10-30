package handlers

import (
	"net/http"
	"time"
)

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// // Get session cookie
	// cookie, err := r.Cookie("session_id")
	// if err == nil {
	// 	// Delete session from database
	// 	_, err = db.Exec("DELETE FROM sessions WHERE id = $1", cookie.Value)
	// 	if err != nil {
	// 		http.Error(w, "Error logging out", http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	// Delete cookie from client
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
