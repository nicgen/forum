package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"forum/cmd/lib"

	"github.com/gofrs/uuid/v5"
)

const (
	maxUploadSize = 20971520 // 20MB in octets
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Limite de taille de requête
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse le formulaire multipart
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, "Fichier trop volumineux", http.StatusBadRequest)
		return
	}

	// Récupère le fichier
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Erreur lors du téléchargement", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Vérifie le type de fichier
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, "Erreur de lecture", http.StatusInternalServerError)
		return
	}
	filetype := http.DetectContentType(buff)

	if !allowedImageTypes[filetype] {
		http.Error(w, "Type de fichier non autorisé", http.StatusBadRequest)
		return
	}

	// Réinitialise le curseur du fichier
	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Erreur de lecture", http.StatusInternalServerError)
		return
	}

	// Génère un nom de fichier unique
	ext := filepath.Ext(handler.Filename)
	filename := uuid.Must(uuid.NewV4()).String() + ext

	// Crée le dossier d'upload s'il n'existe pas avec des permissions spécifiques
	uploadDir := "./static/uploads/"
	errMkdir := os.MkdirAll(uploadDir, 0755)
	if errMkdir != nil {
		http.Error(w, "Impossible de créer le dossier d'upload", http.StatusInternalServerError)
		return
	}

	// Chemin complet du fichier
	filepath := filepath.Join(uploadDir, filename)

	// Crée le fichier avec des permissions spécifiques
	out, errFile := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if errFile != nil {
		http.Error(w, "Impossible de créer le fichier", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Copie le fichier
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Erreur de copie", http.StatusInternalServerError)
		return
	}

	// Récupère les autres données du formulaire
	db := lib.GetDB()

	// Parse the form data (including query parameters and form body)
	err_parse := r.ParseForm()
	if err_parse != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, "Unable to parse form data", http.StatusInternalServerError)
		return
	}

	// Getting the cookie (containing the UUID)
	cookie, _ := r.Cookie("session_id")

	// Storing the form values into variables
	title := r.FormValue("post_title")
	text := r.FormValue("post_text")
	selectedCategories := r.Form["categories"] // This gives you a slice of strings

	var categories string
	for i := 0; i < len(selectedCategories); i++ {
		if i != len(selectedCategories)-1 {
			categories += selectedCategories[i] + ", "
		} else {
			categories += selectedCategories[i]
		}
	}

	// Chemin relatif pour la base de données
	relativePath := filename

	// Insère le post dans la base de données
	state_post := `INSERT INTO Posts (User_UUID, Title, Category_ID, Text, Like, Dislike, CreatedAt, ImagePath) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err_db := db.Exec(state_post, cookie.Value, title, categories, text, 0, 0, time.Now(), relativePath)

	if err_db != nil {

		// Supprime le fichier uploadé en cas d'erreur
		os.Remove(filepath)
		lib.ErrorServer(w, "Erreur lors de l'insertion du post")
		return
	}

	// Redirige vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
