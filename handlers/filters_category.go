package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strings"
)

func Filters_Category(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	formValues := r.URL.Query()
	categories := formValues.Get("Category")

	var role, uuid string
	cookie_role, err_role := r.Cookie("role")
	cookie_uuid, err_role := r.Cookie("session_id")
	if err_role == http.ErrNoCookie {
		role = "Guest"
	} else {
		role = cookie_role.Value
	}
	if err_role == http.ErrNoCookie {
		uuid = "nill"
	} else {
		uuid = cookie_uuid.Value
	}

	data_user := map[string]interface{}{
		"Role":	role,
	}

	var posts []*models.Post
	var err_post error

	state_filters :=
		`SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID, ImagePath FROM Posts ORDER BY CreatedAt DESC`

	rows, err_post := db.Query(state_filters)

	if err_post != nil {
		lib.ErrorMessage(w, nil, err_post.Error())
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			lib.ErrorMessage(w, nil, err.Error())
		}
	}()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID, &post.ImagePath); err != nil {
		}

		category_array := strings.Split(post.Category_ID, " - ")
		is_contained := false
		for i := 0; i < len(category_array); i++ {
			if category_array[i] == categories {
				is_contained = true
			}
		}
		if is_contained {
			// Getting the Username of the person who made the post
			post.Username = lib.CheckUsername(w, post.User_UUID)
	
			if data_user["Role"] != "Guest" {
				// Check the post status with user's uuid and post id
				post.Status = lib.CheckStatus(w, uuid, post.ID, "post")
	
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
				lib.ErrorServer(w, "Error accessing user comments")
			}
	
			defer rows_comment.Close()
	
			for rows_comment.Next() {
				var comment models.Comment
				if err := rows_comment.Scan(&comment.ID, &comment.Text, &comment.Like, &comment.Dislike, &comment.CreatedAt, &comment.User_UUID, &comment.Post_ID); err != nil {
					lib.ErrorServer(w, "Error scanning posts comments")
				}
	
				time_comment := strings.Split(comment.CreatedAt.Format("2006-01-02 15:04:05"), " ")
				comment.Creation_Date = time_comment[0]
				comment.Creation_Hour = time_comment[1]
	
				// Getting the Username of the person who made the comment
				comment.Username = lib.CheckUsername(w, comment.User_UUID)
	
				if data_user["Role"] != "Guest" {
					// Check the post status with user's uuid and post id
					post.Status = lib.CheckStatus(w, uuid, post.ID, "comment")
	
					// Check if the post is from the User making the request
					if post.User_UUID == uuid {
						comment.IsAuthor = "yes"
					} else {
						comment.IsAuthor = "no"
					}
				}
				data_comment := map[string]interface{}{
					"Role":	role,
				}
				comment.Data = data_comment
	
				comments = append(comments, &comment)
			}
	
			if err := rows.Err(); err != nil {
				lib.ErrorServer(w, "Error iterating over user comments")
			}
	
			post.Comments = comments
	
			data_post := map[string]interface{}{
				"Role": role,
			}
			post.Data = data_post
			posts = append(posts, &post)
		}
	}

	if err := rows.Err(); err != nil {
		lib.ErrorMessage(w, nil, err.Error())
		return
	}
	data := lib.DataTest(w, r)
	data["Posts"] = posts
	data = lib.ErrorMessage(w, data, "none")
	lib.RenderTemplate(w, "layout/index", "page/index", data)
}
