package handlers

import "net/http"

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session_id := r.Cookies()
	println("-----------------------------")
	println("Session ID: ", session_id[0].Value)
	println("-----------------------------")
}
