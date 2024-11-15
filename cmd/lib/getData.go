package lib

import (
	"database/sql"
	"fmt"
	"forum/models"
	"net/http"
)

func GetData(db *sql.DB, uuid string, status string, page string, r *http.Request) (map[string]interface{}, string) {
	// Declaring the map we are going to return
	data := map[string]interface{}{}

	if status == "logged" {

		// Getting the User infos from the cookies
		cookie_username, err_username := r.Cookie("username")
		cookie_date, err_date := r.Cookie("creation_date")
		cookie_hour, err_hour := r.Cookie("creation_hour")
		cookie_email, err_hour := r.Cookie("email")
		cookie_role, err_hour := r.Cookie("role")

		// Checking for database requests errors
		if err_username != nil {
			return nil, "Error getting Username from the cookies"
		} else if err_date != nil {
			return nil, "Error getting Creation Date from the cookies"
		} else if err_hour != nil {
			return nil, "Error getting Creation Hour from the cookies"
		}

		// Posts Query based on page
		var state_posts string
		if page == "profile" {
			state_posts = `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts WHERE User_UUID = ? ORDER BY CreatedAt DESC`
		} else if page == "index" {
			state_posts = `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts ORDER BY CreatedAt DESC`
		}

		data_post := map[string]interface{}{
			"Username":     cookie_username.Value,
			"Email":        cookie_email.Value,
			"CreationDate": cookie_date.Value,
			"CreationHour": cookie_hour.Value,
			"Role":         cookie_role.Value,
		}

		// Users posts Request
		var posts []*models.Post
		var rows *sql.Rows
		var err_post error

		if page == "profile" {
			rows, err_post = db.Query(state_posts, uuid)
		} else {
			rows, err_post = db.Query(state_posts)
		}

		if err_post != nil {
			return nil, "Error accessing user posts"
		}
		defer rows.Close()

		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
				return nil, "Error scanning user posts"
			}

			// Getting the Username of the person who made the post
			state_username := `SELECT Username FROM User WHERE UUID = ?`
			err_db := db.QueryRow(state_username, post.User_UUID).Scan(&post.Username)
			if err_db != nil {
				return nil, "Error getting User's Username for the post"
			}

			state_comment := `SELECT ID, Text, Like, Dislike, CreatedAt, User_UUID, Post_ID FROM Comments WHERE Post_ID = ? ORDER BY CreatedAt DESC`
			fmt.Println("post_id getdata: ", post.ID)
			// Users posts Request
			var comments []*models.Comment
			var rows_comment *sql.Rows
			rows_comment, err_comment := db.Query(state_comment, post.ID)
			if err_comment != nil {
				return nil, "Error accessing user comments"
			}

			defer rows_comment.Close()

			for rows_comment.Next() {
				var comment models.Comment
				if err := rows_comment.Scan(&comment.ID, &comment.Text, &comment.Like, &comment.Dislike, &comment.CreatedAt, &comment.User_UUID, &comment.Post_ID); err != nil {
					return nil, "Error scanning posts comments"
				}

				// Getting the Username of the person who made the comment
				state_username := `SELECT Username FROM User WHERE UUID = ?`
				err_db := db.QueryRow(state_username, comment.User_UUID).Scan(&comment.Username)
				if err_db != nil {
					return nil, "Error getting User's Username for the comment"
				}

				// comment.Post_ID = post.ID

				// post.Comments = data_post
				comments = append(comments, &comment)
				fmt.Println("current struct: ", comment)
				fmt.Println("comments struct: ", comments)
			}

			if err := rows.Err(); err != nil {
				return nil, "Error iterating over user posts"
			}

			post.Comments = comments
			post.Data = data_post
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			return nil, "Error iterating over user posts"
		}

		// Retrieve all users for admin view
		var allUsers []models.User
		if cookie_role.Value == "Admin" {
			allUsersQuery := `SELECT UUID, Username, Email, Role FROM User`
			rows, err := db.Query(allUsersQuery)
			if err != nil {
				return nil, "Error accessing user list"
			}
			defer rows.Close()

			for rows.Next() {
				var user models.User
				if err := rows.Scan(&user.UUID, &user.Username, &user.Email, &user.Role); err != nil {
					return nil, "Error scanning users"
				}
				allUsers = append(allUsers, user)
			}

			if err := rows.Err(); err != nil {
				return nil, "Error iterating over users"
			}
		}

		// Storing the data into a map that can be sent into the html
		data = map[string]interface{}{
			"Username":     cookie_username.Value,
			"Email":        cookie_email.Value,
			"CreationDate": cookie_date.Value,
			"CreationHour": cookie_hour.Value,
			"Role":         cookie_role.Value,
			"Posts":        posts,
			"AllUsers":     allUsers,
			"UUID":         uuid,
		}
	} else {

		data_post := map[string]interface{}{
			"Username":     nil,
			"Email":        nil,
			"CreationDate": nil,
			"CreationHour": nil,
			"Role":         "Guest",
		}

		// Not logged in - show all posts
		state_posts := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts ORDER BY CreatedAt DESC`
		var posts []*models.Post
		rows, err := db.Query(state_posts)
		if err != nil {
			return nil, "Error accessing posts"
		}
		defer rows.Close()

		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
				return nil, "Error scanning posts"
			}

			// Getting the Username of the person who made the post
			state_username := `SELECT Username FROM User WHERE UUID = ?`
			err_db := db.QueryRow(state_username, post.User_UUID).Scan(&post.Username)
			if err_db != nil {
				return nil, "Error getting User's Username"
			}

			post.Data = data_post
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			return nil, "Error iterating over posts"
		}

		data = map[string]interface{}{
			"Username":     nil,
			"Email":        nil,
			"CreationDate": nil,
			"CreationHour": nil,
			"Role":         "Guest",
			"Posts":        posts,
		}
	}

	return data, "OK"
}
