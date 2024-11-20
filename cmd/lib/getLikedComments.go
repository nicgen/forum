package lib

import (
	"forum/models"
	"net/http"
	"strings"
)

// ? Function to get liked posts from an User to store it into a data map
func GetLikedComments(w http.ResponseWriter, uuid string, data map[string]interface{}) map[string]interface{} {
	// Making query for the posts liked by the User
	state_liked := `SELECT Comment_ID FROM Reaction WHERE User_UUID = ? AND Comment_ID IS NOT NULL`
	query, err_liked := db.Query(state_liked, uuid)
	if err_liked != nil {
		ErrorServer(w, "Error accessing user's Reactions")
	}
	defer query.Close()

	// Variables that will store the reaction's post id
	react_map := map[string]string{}
	var id, status string

	for query.Next() {
		if err := query.Scan(&id, &status); err != nil {
			ErrorServer(w, "Error scanning user's Reactions")
		}

		react_map[id] = status
	}

	// Ranging over the posts id to get all posts reactions
	var comments_liked []*models.Comment
	state_reacted_posts := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts WHERE ID = ? ORDER BY CreatedAt DESC`
	for key, value := range react_map {
		row_post, err_react := db.Query(state_reacted_posts, key)
		if err_react != nil {
			ErrorServer(w, "Error accessing user's liked comments")
		}
		defer row_post.Close()

		for row_post.Next() {
			var comment_liked models.Comment
			if err := row_post.Scan(&comment_liked.ID, &comment_liked.Text, &comment_liked.Like, &comment_liked.Dislike, &comment_liked.CreatedAt, &comment_liked.User_UUID); err != nil {
				ErrorServer(w, "Error scanning comments data")
			}

			time_comment := strings.Split(comment_liked.CreatedAt.Format("2006-01-02 15:04:05"), " ")
			comment_liked.Creation_Date = time_comment[0]
			comment_liked.Creation_Hour = time_comment[1]

			comment_liked.Status = value
			comments_liked = append(comments_liked, &comment_liked)
		}

		if err := row_post.Err(); err != nil {
			ErrorServer(w, "Error iterating over user's liked comments")
		}
	}
	data["LikedComments"] = comments_liked
	return data
}
