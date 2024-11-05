package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"net/http"
)

// IndexHandler handles requests to the root URL
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Storing database into a variable
	db := lib.GetDB()

	if r.URL.Path != "/" {
		// * generate your error message
		// err := &models.CustomError{
		// 	StatusCode: http.StatusNotFound,
		// 	Message:    "Page Not Found",
		// }
		// Use HandleError to send the error response
		// HandleError(w, err.StatusCode, err.Message)
		// return
		// * alt. use the auto-generated error code & message
		HandleError(w, http.StatusNotFound, "Page Not Found")
		return
	}

	// data := models.PageData{
	// 	Title:  "Forum",
	// 	Header: "Welcome to our Forum project.",
	// 	Content: map[string]template.HTML{
	// 		"Msg_raw":    "<h2>Sub-Title 02.</h1><p>paragraph</>",
	// 		"Msg_styled": "<h1 style=\"text=color: blue;\">Title 01.</h1><p>paragraph with</br>style</>",
	// 	},

	// 	IsError: false,
	// }

	// Users posts Request
	state_posts := `SELECT ID, Category_ID, Title, Text, Like, CreatedAt FROM Posts ORDER BY CreatedAt DESC`
	var posts []*models.Post
	rows, err := db.Query(state_posts)
	if err != nil {
		http.Error(w, "Error accessing user posts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.CreatedAt); err != nil {
			http.Error(w, "Error scanning user posts: "+err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over user posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Posts": posts,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "layout/index", "page/index", data)
}
