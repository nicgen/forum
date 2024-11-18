package lib

import (
	"database/sql"
	"net/http"
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
			"Role": cookie_role.Value,
		}

		// Users posts Request
		var rows *sql.Rows
		var err_post error

		if page == "profile" {
			rows, err_post = db.Query(state_posts, uuid)
			data = GetPosts(w, uuid, state_posts, rows, data, data_post)
		} else {
			rows, err_post = db.Query(state_posts)
			data = GetPosts(w, uuid, state_posts, rows, data, data_post)
		}

		if err_post != nil {
			ErrorServer(w, "Error ranging over posts")
		}

		// Storing the data into a map that can be sent into the html
		data["Username"] = cookie_username.Value
		data["Email"] = cookie_email.Value
		data["CreationDate"] = cookie_date.Value
		data["CreationHour"] = cookie_hour.Value
		data["Role"] = cookie_role.Value
		data["UUID"] = uuid
		data["NavLogin"] = "hide"
		data["NavRegister"] = "hide"

	} else {

		// Map to send to the posts map
		data_post := map[string]interface{}{
			"Role": "Guest",
		}

		// Not logged in - show all posts
		state_posts := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts ORDER BY CreatedAt DESC`
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
