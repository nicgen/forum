package lib

import (
	"database/sql"
	"net/http"
)

func AuthData(uuid, username, email, role, date, hour string, w http.ResponseWriter) map[string]interface{} {
	data := map[string]interface{}{
		"Username": username,
		"Email": email,
		"CreationDate": date,
		"CreationHour": hour,
		"Role": role,
		"UUID": uuid,
		"NavLogin": "hide",
		"NavRegister": "hide",
		"User_UUID": uuid,
	}
	var rows *sql.Rows
	var err_post error
	data_post := map[string]interface{}{
		"Role": role,
	}
	state_posts := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID, ImagePath FROM Posts ORDER BY CreatedAt DESC`
	rows, err_post = db.Query(state_posts)
	data = GetPosts(w, uuid, state_posts, rows, data, data_post)
	if err_post != nil {
		ErrorServer(w, "Error ranging over posts")
	}

	data = ErrorMessage(w, data, "none")
	data = GetCategories(w, data)
	data = GetListOfUsers(w, role, data)
	data = GetLikedPosts(w, uuid, data)
	data = GetLikedComments(w, uuid, data)

	return data
}