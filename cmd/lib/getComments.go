package lib

import (
	"database/sql"
	"fmt"
	"forum/models"
	"net/http"
	"strings"
)

// ? Function to get all the comments a specific User made, including the post infos
func GetComments(db *sql.DB, uuid string, data map[string]interface{}, w http.ResponseWriter, r *http.Request) map[string]interface{} {
	state_comment := `SELECT ID, Text, Like, Dislike, CreatedAt, User_UUID, Post_ID FROM Comments WHERE User_UUID = ? ORDER BY CreatedAt DESC`
	// Users posts Request
	var comments []*models.Comment
	var rows_comment *sql.Rows
	rows_comment, err_comment := db.Query(state_comment, uuid)
	if err_comment != nil {
		ErrorServer(w, "Error accessing user comments")
	}

	defer rows_comment.Close()

	for rows_comment.Next() {
		var comment models.Comment
		if err := rows_comment.Scan(&comment.ID, &comment.Text, &comment.Like, &comment.Dislike, &comment.CreatedAt, &comment.User_UUID, &comment.Post_ID); err != nil {
			ErrorServer(w, "Error scanning posts comments")
		}
		var post models.Post
		fmt.Println("id: ", comment.Post_ID)
		state_post := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts WHERE ID = ? ORDER BY CreatedAt DESC`
		err_comment := db.QueryRow(state_post, comment.Post_ID).Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID)
		if err_comment != nil {
			fmt.Println("error: ", err_comment)
			ErrorServer(w, "Error getting post infos")
		}

		time_comment := strings.Split(comment.CreatedAt.Format("2006-01-02 15:04:05"), " ")
		comment.Creation_Date = time_comment[0]
		comment.Creation_Hour = time_comment[1]

		data_comment := map[string]interface{}{
			"Role":     "Guest",
			"PostInfo": post,
		}
		comment.Data = data_comment

		comments = append(comments, &comment)
	}
	data["ProfileComments"] = comments
	return data
}
