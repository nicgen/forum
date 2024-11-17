package lib

import (
	"database/sql"
	"forum/models"
	"net/http"
	"strings"
)

// ? Function to generate a map containing the informations to show on the html
func GetData(db *sql.DB, uuid string, status string, page string, w http.ResponseWriter, r *http.Request) map[string]interface{} {
	// Declaring the map we are going to return
	data := map[string]interface{}{}

	// Getting categories info on the data map
	data = GetCategories(w, data)

	if status == "logged" {

		// Getting the User infos from the cookies
		cookie_username, err_username := r.Cookie("username")
		cookie_date, err_date := r.Cookie("creation_date")
		cookie_hour, err_hour := r.Cookie("creation_hour")
		cookie_email, err_email := r.Cookie("email")
		cookie_role, err_role := r.Cookie("role")

		// Checking for database requests errors
		if err_username != nil {
			ErrorServer(w, "Error getting Username from the cookies")
		} else if err_date != nil {
			ErrorServer(w, "Error getting Creation Date from the cookies")
		} else if err_hour != nil {
			ErrorServer(w, "Error getting Creation Hour from the cookies")
		} else if err_email != nil {
			ErrorServer(w, "Error getting email from the cookies")
		} else if err_role != nil {
			ErrorServer(w, "Error getting User role from the cookies")
		}

		// Storing the post liked into the data map
		data = GetLikedPosts(w, uuid, data)

		// Storing the list of Users into the data map if the role is Admin
		data = GetListOfUsers(w, cookie_role.Value, data)

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
			ErrorServer(w, "Error accessing user posts")
		}
		defer rows.Close()

		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
				ErrorServer(w, "Error scanning user posts")
			}

			time_post := strings.Split(post.CreatedAt.Format("2006-01-02 15:04:05"), " ")
			post.Creation_Date = time_post[0]
			post.Creation_Hour = time_post[1]

			// Getting the Username of the person who made the post
			post.Username = CheckUsername(w, post.User_UUID)

			// Check the post status with user's uuid and post id
			post.Status = CheckStatus(w, post.User_UUID, post.ID, "post")

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

				// Check the comment status with user's uuid and comment id
				comment.Status = CheckStatus(w, comment.User_UUID, comment.ID, "comment")

				comments = append(comments, &comment)
			}

			if err := rows_comment.Err(); err != nil {
				ErrorServer(w, "Error iterating over user comments")
			}

			post.Comments = comments
			post.Data = data_post
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			ErrorServer(w, "Error iterating over user posts")
		}

		// Storing the data into a map that can be sent into the html
		data["Username"] = cookie_username.Value
		data["Email"] = cookie_email.Value
		data["CreationDate"] = cookie_date.Value
		data["CreationHour"] = cookie_hour.Value
		data["Role"] = cookie_role.Value
		data["Posts"] = posts
		data["UUID"] = uuid
		data["NavLogin"] = "hide"
		data["NavRegister"] = "hide"

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
			ErrorServer(w, "Error accessing posts")
		}
		defer rows.Close()

		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
				ErrorServer(w, "Error scanning posts")
			}

			time_post := strings.Split(post.CreatedAt.Format("2006-01-02 15:04:05"), " ")
			post.Creation_Date = time_post[0]
			post.Creation_Hour = time_post[1]

			// Getting the Username of the person who made the post
			state_username := `SELECT Username FROM User WHERE UUID = ?`
			err_db := db.QueryRow(state_username, post.User_UUID).Scan(&post.Username)
			if err_db != nil {
				ErrorServer(w, "Error getting User's Username")
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
				state_username := `SELECT Username FROM User WHERE UUID = ?`
				err_db := db.QueryRow(state_username, comment.User_UUID).Scan(&comment.Username)
				if err_db != nil {
					ErrorServer(w, "Error getting User's Username for the comment")
				}

				comments = append(comments, &comment)
			}

			if err := rows.Err(); err != nil {
				ErrorServer(w, "Error iterating over user posts")
			}

			post.Comments = comments
			post.Data = data_post

			post.Data = data_post
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			ErrorServer(w, "Error iterating over posts")
		}

		data["Username"] = nil
		data["Email"] = nil
		data["CreationDate"] = nil
		data["CreationHour"] = nil
		data["Role"] = "Guest"
		data["Posts"] = posts
		data["NavLogin"] = "hide"
		data["NavRegister"] = "hide"
	}

	return data
}
