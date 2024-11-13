package lib

import (
	"database/sql"
	"fmt"
	"forum/models"
	"strings"
	"time"
)

func GetData(db *sql.DB, uuid string, status string, page string) (map[string]interface{}, string) {
	// Declaring the map we are going to return
	data := map[string]interface{}{}

	if status == "logged" {
		// Setting up the variables of the User informations
		var user_username, user_email, user_creation, user_role string

		// Setting up the states for the database requests
		state_uuid := `SELECT Username FROM User WHERE UUID = ?`
		state_email := `SELECT Email FROM User WHERE UUID = ?`
		state_creation := `SELECT CreatedAt FROM User WHERE UUID = ?`
		state_role := `SELECT Role FROM User WHERE UUID = ?`

		// Making requests to the database
		err_uuid := db.QueryRow(state_uuid, uuid).Scan(&user_username)
		err_email := db.QueryRow(state_email, uuid).Scan(&user_email)
		err_creation := db.QueryRow(state_creation, uuid).Scan(&user_creation)
		err_role := db.QueryRow(state_role, uuid).Scan(&user_role)

		// Checking for database requests errors
		if err_uuid != nil {
			return nil, "Error accessing User UUID"
		} else if err_email != nil {
			return nil, "Error accessing User EMAIL"
		} else if err_creation != nil {
			return nil, "Error accessing User CREATION DATE"
		} else if err_role != nil {
			return nil, "Error accessing User ROLE"
		}

		// Posts Query based on page
		var state_posts string
		if page == "profile" {
			state_posts = `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt FROM Posts WHERE User_UUID = ? ORDER BY CreatedAt DESC`
		} else if page == "index" {
			state_posts = `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt FROM Posts ORDER BY CreatedAt DESC`
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
			if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt); err != nil {
				return nil, "Error scanning user posts"
			}

			fmt.Println("---------------------------")
			fmt.Println("Actual time: ", time.Now())
			fmt.Println("post creation date: ", post.CreatedAt)
			fmt.Println("---------------------------")
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			return nil, "Error iterating over user posts"
		}

		// Spliting the creation date into 2 different values
		creation := strings.Split(user_creation, "T")
		creation_Date := creation[0]
		creation_Hour := creation[1][:len(creation[1])-1]

		// Retrieve all users for admin view
		var allUsers []models.User
		if user_role == "Admin" {
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
			"Username":     user_username,
			"Email":        user_email,
			"CreationDate": creation_Date,
			"CreationHour": creation_Hour,
			"Role":         user_role,
			"Posts":        posts,
			"AllUsers":     allUsers,
		}
	} else {
		// Not logged in - show all posts
		state_posts := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt FROM Posts ORDER BY CreatedAt DESC`
		var posts []*models.Post
		rows, err := db.Query(state_posts)
		if err != nil {
			return nil, "Error accessing posts"
		}
		defer rows.Close()

		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt); err != nil {
				return nil, "Error scanning posts"
			}
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
