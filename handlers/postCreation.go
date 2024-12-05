package handlers

import (
	"forum/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"forum/cmd/lib"

	"github.com/gofrs/uuid/v5"
)

const (
	maxUploadSize = 20971520 // 20MB en octets
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
		if strings.Contains(err.Error(), "request body too large") {
			data := lib.DataTest(w, r)
			data = lib.ErrorMessage(w, data, "ImageContent")

			// Redirigez vers la page d'index au lieu de la page de création de post
			lib.RenderTemplate(w, "layout/index", "page/index", data)
		}
		//Erreur critique : erreur lors du téléchargemement
		err := &models.CustomError{
			StatusCode: http.StatusBadRequest,
			Message:    "Erreur lors du téléchargement",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Récupération du cookie de session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		//Erreur critique: session non trouvée
		err := &models.CustomError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Session not found, please log in again.",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Vérifie si le champ "image" est présent
	file, handler, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		//Erreur critique: Erreur lors du téléchargement
		err := &models.CustomError{
			StatusCode: http.StatusBadRequest,
			Message:    "Erreur lors du téléchargement",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var filename string
	if err == http.ErrMissingFile {
		// Pas d'image téléchargée, continuez sans image
		filename = "" // Pas de fichier
	} else {
		// Traitement de l'image
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			//Erreur critique : erreur de lecture
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Erreur de lecture",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}

		filetype := http.DetectContentType(buff)

		if !allowedImageTypes[filetype] {
			//Erreur non critique : type de fichier non autorisé
			http.Error(w, "Type de fichier non autorisé", http.StatusBadRequest)
			return
		}

		// Réinitialise le curseur du fichier
		_, err = file.Seek(0, 0)
		if err != nil {
			// Erreur critique : erreur de lecture
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Erreur de lecture",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}

		// Vérifie la taille du fichier
		if handler.Size > maxUploadSize {
			//Erreur non critique : image trop volumineuse
			data := lib.DataTest(w, r)
			data = lib.ErrorMessage(w, data, "ImageContent")

			// Rendre la page d'index avec le message d'erreur
			lib.RenderTemplate(w, "layout/index", "page/index", data)
		}

		// Génère un nom de fichier unique
		ext := filepath.Ext(handler.Filename)
		filename = uuid.Must(uuid.NewV4()).String() + ext

		// Crée le dossier d'upload s'il n'existe pas avec des permissions spécifiques
		uploadDir := "./static/uploads/"
		errMkdir := os.MkdirAll(uploadDir, 0755)
		if errMkdir != nil {
			// Erreur critique : impossible de créer le dossier d'upload
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Impossible de créer le dossier d'upload",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}

		// Chemin complet du fichier
		filePath := filepath.Join(uploadDir, filename)

		// Crée le fichier avec des permissions spécifiques
		out, errFile := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if errFile != nil {
			// Erreur critique : impossible de créer le fichier
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Impossible de créer le fichier",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}
		defer out.Close()

		// Copie le fichier
		_, err = io.Copy(out, file)
		if err != nil {
			//Erreur critique: erreur de copie
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Erreur de copie de fichier",
			}
			HandleError(w, err.StatusCode, err.Message)
		}
	}

	// Récupère les autres données du formulaire
	db := lib.GetDB()

	// Récupération des données du formulaire
	errParse := r.ParseForm()
	if errParse != nil {
		http.Error(w, "Unable to parse form data", http.StatusInternalServerError)
		return
	}

	// Récupération des valeurs du formulaire
	title := r.FormValue("post_title")
	text := r.FormValue("post_text")
	selectedCategories := r.Form["categories"]

	var categories string
	for i := 0; i < len(selectedCategories); i++ {
		if i != len(selectedCategories)-1 {
			categories += selectedCategories[i] + " - "
		} else {
			categories += selectedCategories[i]
		}
	}

	// Chemin relatif pour la base de données
	relativePath := filename

	// Insère le post dans la base de données
	statePost := `INSERT INTO Posts (User_UUID, Title, Category_ID, Text, Like, Dislike, CreatedAt, ImagePath) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, errDb := db.Exec(statePost, cookie.Value, title, categories, text, 0, 0, time.Now(), relativePath)

	if errDb != nil {
		// Erreur critique : échec de l'insertion du post
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error inserting new post, please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)

		// Supprime le fichier uploadé en cas d'erreur
		if filename != "" {
			os.Remove(filepath.Join("./static/uploads/", filename))
		}
		return
	}

	// Redirige vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
