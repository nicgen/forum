package lib

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tmpl *template.Template

func init() {
	rootDir := "./templates"
	funcMap := template.FuncMap{}
	var err error
	tmpl, err = findAndParseTemplates(rootDir, funcMap)
	if err != nil {
		log.Fatal(err)
	}
}

func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := os.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

func RenderTemplate(w http.ResponseWriter, layoutName, tmplName string, data interface{}) {

	// Convertir data en map si ce n'est pas déjà le cas
	var dataMap map[string]interface{}
	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	default:
		dataMap = map[string]interface{}{}
	}

	// Exécution du template spécifique d'abord
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, tmplName, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Créer la structure layoutData
	layoutData := struct {
		Content     template.HTML
		Data        interface{}
		UserContent interface{}
		Post        interface{} // Changé de *models.Post à interface{}
	}{
		Content:     template.HTML(buf.String()),
		Data:        data,
		UserContent: data,
		Post:        dataMap["Post"], // Utilisez dataMap["Post"] sans type assertion
	}

	// Exécution du layout template
	err = tmpl.ExecuteTemplate(w, layoutName, layoutData)
	if err != nil {
		log.Printf("Error executing layout template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
