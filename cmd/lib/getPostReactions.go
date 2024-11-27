package lib

import (
	"database/sql"
	"fmt"
	"forum/models"
)

func GetPostReactions(db *sql.DB, data map[string]interface{}, uuid string) (map[string]interface{}, string) {
	// Making query for the posts liked by the User
	state_liked := `SELECT Post_ID FROM Reaction WHERE User_UUID = ?`
	query, err_liked := db.Query(state_liked, uuid)
	fmt.Println("Error: ", err_liked)
	if err_liked != nil {
		return nil, "Error accessing user's Reactions"
	}
	defer query.Close()

	// Variables that will store the reaction's post id
	react_tab := []string{}
	reaction := ""

	for query.Next() {
		if err := query.Scan(&reaction); err != nil {
			fmt.Println("Error2: ", err)
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
	data["LikedPosts"] = posts_liked
	return data, "OK"
}
