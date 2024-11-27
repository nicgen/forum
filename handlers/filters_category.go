package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strings"
)

func Filters_Category(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	formValues := r.URL.Query()
	categories := formValues.Get("Category")

	var posts []*models.Post
	var err_post error

	state_filters :=
		`SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID
	FROM Posts `

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
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
		}

		category_array := strings.Split(post.Category_ID, " - ")
		is_contained := false
		for i := 0; i < len(category_array); i++ {
			if category_array[i] == categories {
				is_contained = true
			}
		}
		if is_contained {
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
