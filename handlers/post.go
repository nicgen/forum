package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	post_id := r.URL.Query().Get("id")

	var post models.Post
	query := `
        SELECT 
            p.ID, 
            c.Name, 
            p.Title, 
            p.Text, 
            p.Like, 
            p.Dislike, 
            p.CreatedAt, 
            u.Username, 
            p.ImagePath,
            p.User_UUID,
            u.Role
        FROM Posts p
        JOIN Categories c ON p.Category_ID = c.ID
        JOIN User u ON p.User_UUID = u.UUID
        WHERE p.ID = ?
    `
	err := db.QueryRow(query, post_id).Scan(
		&post.ID,
		&post.Category_ID,
		&post.Title,
		&post.Text,
		&post.Like,
		&post.Dislike,
		&post.CreatedAt,
		&post.Username,
		&post.ImagePath,
		&post.User_UUID,
		&post.Role,
	)
	if err != nil {
		log.Printf("Debug - Post Details: ID=%d, Title=%s, ImagePath=%s", post.ID, post.Title, post.ImagePath)
		log.Printf("Error fetching post details: %v", err)
		lib.ErrorServer(w, "Error fetching post details")
		return
	}

	// Récupération des commentaires
	state_comment := `
        SELECT 
            c.ID, 
            c.Text, 
            c.Like, 
            c.Dislike, 
            c.CreatedAt, 
            u.Username,
            c.User_UUID
        FROM Comments c
        JOIN User u ON c.User_UUID = u.UUID
        WHERE c.Post_ID = ?
        ORDER BY c.CreatedAt DESC
    `

	rows_comment, err_comment := db.Query(state_comment, post_id)
	if err_comment != nil {
		log.Printf("Error querying comments: %v", err_comment)
		lib.ErrorServer(w, "Error accessing user comments")
		return
	}
	defer rows_comment.Close()

	var comments []*models.Comment
	for rows_comment.Next() {
		var comment models.Comment
		var createdAt time.Time

		err := rows_comment.Scan(
			&comment.ID,
			&comment.Text,
			&comment.Like,
			&comment.Dislike,
			&createdAt,
			&comment.Username,
			&comment.User_UUID,
		)

		if err != nil {
			log.Printf("Error scanning comment: %v", err)
			lib.ErrorServer(w, "Error scanning posts comments")
			return
		}

		comment.CreatedAt = createdAt
		comments = append(comments, &comment)
	}

	post.Comments = comments

	// Préparation des données pour le template
	data := lib.DataTest(w, r)

	if post.ImagePath != "" {
		log.Printf("Image path is not empty: %s", post.ImagePath)
	} else {
		log.Printf("Warning: Image path is empty for post %d", post.ID)
	}

	// Vérification du statut de like/dislike
	cookie, _ := r.Cookie("session_id")
	if cookie != nil {
		// Convertissez post.ID en string
		post.Status = lib.CheckStatus(w, cookie.Value, strconv.Itoa(post.ID), "post")

		// Vérifier si l'utilisateur est l'auteur
		if post.User_UUID == cookie.Value {
			post.IsAuthor = "yes"
		} else {
			post.IsAuthor = "no"
		}

		// Récupérer les données utilisateur
		post.Data = data
	}

	data["Post"] = post
	log.Printf("Rendering post with ImagePath: %s", post.ImagePath)
	lib.RenderTemplate(w, "layout/index", "page/post", data)
}
