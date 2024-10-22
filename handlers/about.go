package handlers

import (
	"forum/models"
	"html/template"
	"net/http"
)

// IndexHandler handles requests to the /about URL
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		Title:  "About",
		Header: "About Us",
		Content: map[string]template.HTML{
			"Msg_raw":    "<h1>Title 01.</h1><p>paragraph</>",
			"Msg_styled": "<h1 style=\"text=color: blue;\">Title 01.</h1><p>paragraph with</br>style</>",
		},
		IsError: false,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "layout/alt", "page/about", data)
}
