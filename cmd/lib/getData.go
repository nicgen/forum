package lib

import (
	"database/sql"
	"forum/models"
	"net/http"
	"strings"
)

func GetData(db *sql.DB, uuid string, status string, page string, r *http.Request) (map[string]interface{}, string) {
	// Declaring the map we are going to return
	data := map[string]interface{}{}

	// Getting the Categories from the database
	var categories []*models.Category
	state_categories := `SELECT ID, Name FROM Categories`
	query_category, err_liked := db.Query(state_categories, uuid)
	if err_liked != nil {
		return nil, "Error accessing Categories"
	}
	defer query_category.Close()

	for query_category.Next() {
		var category models.Category
		if err := query_category.Scan(&category.ID, &category.Name); err != nil {
			return nil, "Error scanning Categories"
		}

		categories = append(categories, &category)
	}

	if status == "logged" {

		// Getting the User infos from the cookies
		cookie_username, err_username := r.Cookie("username")
		cookie_date, err_date := r.Cookie("creation_date")
		cookie_hour, err_hour := r.Cookie("creation_hour")
		cookie_email, err_email := r.Cookie("email")
		cookie_role, err_role := r.Cookie("role")

		// Checking for database requests errors
		if err_username != nil {
			return nil, "Error getting Username from the cookies"
		} else if err_date != nil {
			return nil, "Error getting Creation Date from the cookies"
		} else if err_hour != nil {
			return nil, "Error getting Creation Hour from the cookies"
		} else if err_email != nil {
			return nil, "Error getting email from the cookies"
		} else if err_role != nil {
			return nil, "Error getting User role from the cookies"
		}

		// Making query for the posts liked by the User
		state_liked := `SELECT Post_ID FROM Reaction WHERE User_UUID = ?`
		query, err_liked := db.Query(state_liked, uuid)
		if err_liked != nil {
			return nil, "Error accessing user's Reactions"
		}
		defer query.Close()

		// Variables that will store the reaction's post id
		react_tab := []string{}
		reaction := ""

		for query.Next() {
			if err := query.Scan(&reaction); err != nil {
				return nil, "Error scanning user's Reactions"
			}

			react_tab = append(react_tab, reaction)
		}

		// Ranging over the posts id to get all posts reactions
		var posts_liked []*models.Post
		state_reacted_posts := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts WHERE ID = ? ORDER BY CreatedAt DESC`
		for i := 0; i < len(react_tab); i++ {
			row_post, err_react := db.Query(state_reacted_posts, react_tab[i])
			if err_react != nil {
				return nil, "Error accessing user's liked posts"
			}
			defer row_post.Close()

			for row_post.Next() {
				var post_liked models.Post
				if err := row_post.Scan(&post_liked.ID, &post_liked.Category_ID, &post_liked.Title, &post_liked.Text, &post_liked.Like, &post_liked.Dislike, &post_liked.CreatedAt, &post_liked.User_UUID); err != nil {
					return nil, "Error scanning posts data"
				}

				posts_liked = append(posts_liked, &post_liked)
			}

			if err := row_post.Err(); err != nil {
				return nil, "Error iterating over user's liked posts"
			}
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

			time_post := strings.Split(post.CreatedAt.Format("2006-01-02 15:04:05"), " ")
			post.Creation_Date = time_post[0]
			post.Creation_Hour = time_post[1]

			// Getting the Username of the person who made the post
			state_username := `SELECT Username FROM User WHERE UUID = ?`
			err_db := db.QueryRow(state_username, post.User_UUID).Scan(&post.Username)
			if err_db != nil {
				return nil, "Error getting User's Username for the post"
			}

			// Checking if the post is liked by the User or not
			var status_post string
			state_status := `SELECT Status FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
			err_status := db.QueryRow(state_status, post.User_UUID, post.ID).Scan(&status_post)
			if err_status == sql.ErrNoRows {
				status = ""
			} else if err_status != nil {
				return nil, "Error checking post status"
			}
			post.Status = status_post

			state_comment := `SELECT ID, Text, Like, Dislike, CreatedAt, User_UUID, Post_ID FROM Comments WHERE Post_ID = ? ORDER BY CreatedAt DESC`
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

				time_comment := strings.Split(comment.CreatedAt.Format("2006-01-02 15:04:05"), " ")
				comment.Creation_Date = time_comment[0]
				comment.Creation_Hour = time_comment[1]

				// Getting the Username of the person who made the comment
				state_username := `SELECT Username FROM User WHERE UUID = ?`
				err_db := db.QueryRow(state_username, comment.User_UUID).Scan(&comment.Username)
				if err_db != nil {
					return nil, "Error getting User's Username for the comment"
				}

				// Checking if the comment is liked by the User or not
				var status_comment string
				state_status := `SELECT Status FROM Reaction WHERE User_UUID = ? AND Comment_ID = ?`
				err_status := db.QueryRow(state_status, comment.User_UUID, comment.ID).Scan(&status_comment)
				if err_status == sql.ErrNoRows {
					status = ""
				} else if err_status != nil {
					return nil, "Error checking post status"
				}
				comment.Status = status_comment

				comments = append(comments, &comment)
			}

			if err := rows_comment.Err(); err != nil {
				return nil, "Error iterating over user comments"
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
			"LikedPosts":   posts_liked,
			"AllUsers":     allUsers,
			"UUID":         uuid,
			"NavLogin":     "hide",
			"NavRegister":  "hide",
			"Categories":   categories,
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

			time_post := strings.Split(post.CreatedAt.Format("2006-01-02 15:04:05"), " ")
			post.Creation_Date = time_post[0]
			post.Creation_Hour = time_post[1]

			// Getting the Username of the person who made the post
			state_username := `SELECT Username FROM User WHERE UUID = ?`
			err_db := db.QueryRow(state_username, post.User_UUID).Scan(&post.Username)
			if err_db != nil {
				return nil, "Error getting User's Username"
			}
			state_comment := `SELECT ID, Text, Like, Dislike, CreatedAt, User_UUID, Post_ID FROM Comments WHERE Post_ID = ? ORDER BY CreatedAt DESC`
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

				time_comment := strings.Split(comment.CreatedAt.Format("2006-01-02 15:04:05"), " ")
				comment.Creation_Date = time_comment[0]
				comment.Creation_Hour = time_comment[1]

				// Getting the Username of the person who made the comment
				state_username := `SELECT Username FROM User WHERE UUID = ?`
				err_db := db.QueryRow(state_username, comment.User_UUID).Scan(&comment.Username)
				if err_db != nil {
					return nil, "Error getting User's Username for the comment"
				}

				comments = append(comments, &comment)
			}

			if err := rows.Err(); err != nil {
				return nil, "Error iterating over user posts"
			}

			post.Comments = comments
			post.Data = data_post

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
			"NavLogin":     "hide",
			"NavRegister":  "hide",
			"Categories":   categories,
		}
	}

	return data, "OK"
}
