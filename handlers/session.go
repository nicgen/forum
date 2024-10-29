package handlers

import (
	"net/http"
	"time"
)

func AttributeSession(user_uuid string, w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user_uuid,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	// print("cookie in login: ", &http.Cookie{
	// 	Name:     "session_token",
	// 	Value:    sessionToken,
	// 	Expires:  time.Now().Add(24 * time.Hour),
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode,
	// }, "\n")
}
