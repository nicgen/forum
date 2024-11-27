package lib

import (
	"database/sql"
	"fmt"
	"forum/models"
	"net/http"
	"strings"
)

// ? Function to get the posts bases on the query
func GetPosts(w http.ResponseWriter, uuid, state string, rows *sql.Rows, data, data_user map[string]interface{}) map[string]interface{} {

	var posts []*models.Post
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID, &post.ImagePath); err != nil {
			fmt.Println("test: ", err)
			ErrorServer(w, "Error scanning posts")
		}

		time_post := strings.Split(post.CreatedAt.Format("2006-01-02 15:04:05"), " ")
		post.Creation_Date = time_post[0]
		post.Creation_Hour = time_post[1]

		// Getting the Username of the person who made the post
		post.Username = CheckUsername(w, post.User_UUID)

		if data_user["Role"] != "Guest" {
			// Check the post status with user's uuid and post id
			post.Status = CheckStatus(w, uuid, post.ID, "post")

			// Check if the post is from the User making the request
			if post.User_UUID == uuid {
				post.IsAuthor = "yes"
			} else {
				post.IsAuthor = "no"
			}
		}

		state_comment := `SELECT ID, Text, Like, Dislike, CreatedAt, User_UUID, Post_ID FROM Comments WHERE Post_ID = ? ORDER BY CreatedAt DESC`
		// Users posts Request
		var comments []*models.Comment
		var rows_comment *sql.Rows
		rows_comment, err_comment := db.Query(state_comment, post.ID)
		if err_comment != nil {
			ErrorServer(w, "Error accessing user comments")
		}

		defer rows_comment.Close()

		for rows_comment.Next() {
			var comment models.Comment
			if err := rows_comment.Scan(&comment.ID, &comment.Text, &comment.Like, &comment.Dislike, &comment.CreatedAt, &comment.User_UUID, &comment.Post_ID); err != nil {
				ErrorServer(w, "Error scanning posts comments")
			}

			time_comment := strings.Split(comment.CreatedAt.Format("2006-01-02 15:04:05"), " ")
			comment.Creation_Date = time_comment[0]
			comment.Creation_Hour = time_comment[1]

			// Getting the Username of the person who made the comment
			comment.Username = CheckUsername(w, comment.User_UUID)

			if data_user["Role"] != "Guest" {
				// Check the post status with user's uuid and post id
				post.Status = CheckStatus(w, uuid, post.ID, "comment")

				// Check if the post is from the User making the request
				if post.User_UUID == uuid {
					comment.IsAuthor = "yes"
				} else {
					comment.IsAuthor = "no"
				}
			}

			data_comment := map[string]interface{}{
				"Role": "Guest",
			}
			comment.Data = data_comment

			comments = append(comments, &comment)
		}

		if err := rows.Err(); err != nil {
			ErrorServer(w, "Error iterating over user comments")
		}

		post.Comments = comments
		post.Data = data_user

		post.Data = data_user
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		ErrorServer(w, "Error iterating over posts")
	}
	data["Posts"] = posts
	return data
}
