package handlers

import "net/http"

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	cookie := r.Cookies()
	print("cookie in profile: ", cookie, "\n")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
