package lib

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
)

// ? Function to generate a map containing the informations to show on the html
func GetData(db *sql.DB, uuid string, status string, page string, w http.ResponseWriter, r *http.Request) map[string]interface{} {
	// Declaring the map we are going to return
	data := map[string]interface{}{}

	// Getting categories info on the data map
	data = GetCategories(w, data)

	if status == "logged" {

		// Storing the post and comments liked into the data map
		data = GetLikedPosts(w, uuid, data)
		data = GetLikedComments(w, uuid, data)

		// Setting up the User infos
		var username, date, hour, email, role string

		if page == "profile" || page == "index" {

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

			// Giving values to the variables set before
			username = cookie_username.Value
			date = cookie_date.Value
			hour = cookie_hour.Value
			email = cookie_email.Value
			role = cookie_role.Value

		} else if page == "profile_user" {

			var createdAt time.Time
			state_info := `SELECT Username, CreatedAt, Email, Role FROM User WHERE UUID = ?`
			err_info := db.QueryRow(state_info, uuid).Scan(&username, &createdAt, &email, &role)
			if err_info != nil {
				ErrorServer(w, "Error getting User infos")
			}

			// Checking the cookie values
			data["User_UUID"] = "user_profile"

			// Storing date informations into the map
			time_comment := strings.Split(createdAt.Format("2006-01-02 15:04:05"), " ")
			date = time_comment[0]
			hour = time_comment[1]
		}

		// Storing the list of Users into the data map if the role is Admin
		data = GetListOfUsers(w, role, data)

		// Posts Query based on page
		var state_posts string
		if page == "profile" || page == "profile_user" {
			state_posts = `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID, ImagePath FROM Posts WHERE User_UUID = ? ORDER BY CreatedAt DESC`
		} else if page == "index" {
			state_posts = `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID, ImagePath FROM Posts ORDER BY CreatedAt DESC`
		}

		data_post := map[string]interface{}{
			"Role": role,
		}

		// Users posts Request
		var rows *sql.Rows
		var err_post error

		if page == "profile" || page == "profile_user" {
			rows, err_post = db.Query(state_posts, uuid)
		} else {
			rows, err_post = db.Query(state_posts)
		}

		data = GetPosts(w, uuid, state_posts, rows, data, data_post)

		if err_post != nil {
			ErrorServer(w, "Error ranging over posts")
		}

		// Storing the data into a map that can be sent into the html
		data["Username"] = username
		data["Email"] = email
		data["CreationDate"] = date
		data["CreationHour"] = hour
		data["Role"] = role
		data["UUID"] = uuid
		data["NavLogin"] = "hide"
		data["NavRegister"] = "hide"
		data = ErrorMessage(w, data, "none")
		if page != "profile_user" {
			data["User_UUID"] = uuid
		}

	} else {

		// Map to send to the posts map
		data_post := map[string]interface{}{
			"Role": "Guest",
		}

		// Not logged in - show all posts
		state_posts := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID, ImagePath FROM Posts ORDER BY CreatedAt DESC`
		rows, err := db.Query(state_posts)
		if err != nil {
			ErrorServer(w, "Error accessing posts")
		}

		// Setting up the map
		data = GetPosts(w, "", state_posts, rows, data, data_post)
		data["Role"] = "Guest"
		data["NavLogin"] = "hide"
		data["NavRegister"] = "hide"
	}

	return data
}
