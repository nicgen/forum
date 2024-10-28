package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// load all templates
var templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))

// ensures all templates are parsed only once when the application starts
func init() {
	var err error
	// Parse all templates in the templates directory
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

// renderer function (handles different layouts)
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	layout := "layout"

	// ? ALT. if you want to use another layout for a page (here: the error page)
	// if tmpl == "error" {
	// 	// Use a different layout
	// 	layout = "layout_alt"
	// }

	// execute the specific template first
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute the layout, passing the executed template as content
	layoutData := struct {
		Content template.HTML
		Data    interface{}
	}{
		Content: template.HTML(buf.String()),
		Data:    data,
	}

	err = templates.ExecuteTemplate(w, layout, layoutData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
