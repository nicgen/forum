package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strings"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Getting database data into a variable
	db := lib.GetDB()

	// Getting the id of the post selected
	post_id := r.URL.Query().Get("id")

	// Getting User data and posts
	data := lib.DataTest(w, r)
	data = lib.ErrorMessage(w, data, "none")

	state_comment := `SELECT ID, Text, Like, Dislike, CreatedAt, User_UUID, Post_ID FROM Comments WHERE Post_ID = ? ORDER BY CreatedAt DESC`
	// Users posts Request
	var comments []*models.Comment
	var rows_comment *sql.Rows
	rows_comment, err_comment := db.Query(state_comment, post_id)
	if err_comment != nil {
		lib.ErrorServer(w, "Error accessing user comments")
	}

	defer rows_comment.Close()

	for rows_comment.Next() {
		var comment models.Comment
		if err := rows_comment.Scan(&comment.ID, &comment.Text, &comment.Like, &comment.Dislike, &comment.CreatedAt, &comment.User_UUID, &comment.Post_ID); err != nil {
			lib.ErrorServer(w, "Error scanning posts comments")
		}

		time_comment := strings.Split(comment.CreatedAt.Format("2006-01-02 15:04:05"), " ")
		comment.Creation_Date = time_comment[0]
		comment.Creation_Hour = time_comment[1]

		// Getting the Username of the person who made the comment
		state_username := `SELECT Username FROM User WHERE UUID = ?`
		err_db := db.QueryRow(state_username, comment.User_UUID).Scan(&comment.Username)
		if err_db != nil {
			lib.ErrorServer(w, "Error getting User's Username for the comment")
		}

		comments = append(comments, &comment)
	}

	lib.RenderTemplate(w, "layout/index", "page/post", data)
}
