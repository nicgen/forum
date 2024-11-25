package handlers

import (
	"forum/cmd/lib"
	"forum/models"
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
		lib.ErrorServer(w, "Error fetching post details")
		return
	}

	// Récupération des commentaires
	stateComment := `
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

	rowsComment, errComment := db.Query(stateComment, post_id)
	if errComment != nil {
		lib.ErrorServer(w, "Error accessing comments")
		return
	}
	defer rowsComment.Close()

	var comments []*models.Comment
	for rowsComment.Next() {
		var comment models.Comment
		var createdAt time.Time

		err := rowsComment.Scan(
			&comment.ID,
			&comment.Text,
			&comment.Like,
			&comment.Dislike,
			&createdAt,
			&comment.Username,
			&comment.User_UUID,
		)

		if err != nil {
			lib.ErrorServer(w, "Error scanning comments")
			return
		}

		comment.CreatedAt = createdAt
		comments = append(comments, &comment)
	}

	post.Comments = comments

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
		data := lib.DataTest(w, r)
		post.Data = data
	}

	// Préparez les données pour le template
	data := lib.DataTest(w, r)
	data["Post"] = post
	data["Comments"] = comments

	lib.RenderTemplate(w, "layout/index", "page/index", data)
}
