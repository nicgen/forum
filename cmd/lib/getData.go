package lib

import (
	"database/sql"
	"forum/models"
	"net/http"
)

func GetData(db *sql.DB, uuid string, status string, page string, r *http.Request) (map[string]interface{}, string) {
	// Declaring the map we are going to return
	data := map[string]interface{}{}

	if status == "logged" {

		// Making the database request for the User role
		var user_role, user_email string
		state_email := `SELECT Email FROM User WHERE UUID = ?`
		state_role := `SELECT Role FROM User WHERE UUID = ?`
		err_email := db.QueryRow(state_email, uuid).Scan(&user_email)
		err_role := db.QueryRow(state_role, uuid).Scan(&user_role)

		// Getting the User infos from the cookies
		cookie_username, err_username := r.Cookie("username")
		cookie_date, err_date := r.Cookie("creation_date")
		cookie_hour, err_hour := r.Cookie("creation_hour")

		// Checking for database requests errors
		if err_username != nil {
			return nil, "Error getting Username from the cookies"
		} else if err_date != nil {
			return nil, "Error getting Creation Date from the cookies"
		} else if err_hour != nil {
			return nil, "Error getting Creation Hour from the cookies"
		} else if err_role != nil {
			return nil, "Error getting User role from the database"
		} else if err_email != nil {
			return nil, "Error getting User email from the database"
		}

		// Posts Query based on page
		var state_posts string
		if page == "profile" {
			state_posts = `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts WHERE User_UUID = ? ORDER BY CreatedAt DESC`
		} else if page == "index" {
			state_posts = `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts ORDER BY CreatedAt DESC`
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
				return nil, "Error getting User's Username"
			}
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			return nil, "Error iterating over user posts"
		}

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
			"Username":     cookie_username.Value,
			"Email":        user_email,
			"CreationDate": cookie_date.Value,
			"CreationHour": cookie_hour.Value,
			"Role":         user_role,
			"Posts":        posts,
			"AllUsers":     allUsers,
			"UUID":         uuid,
		}
	} else {
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
