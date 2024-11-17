package lib

import (
	"database/sql"
	"net/http"
)

// ? Function to check if a post is liked or not by an User
func CheckStatus(w http.ResponseWriter, uuid, id, status_type string) string {
	var status_comment, state_status string
	if status_type == "post" {
		state_status = `SELECT Status FROM Reaction WHERE User_UUID = ? AND Post_ID = ?`
	} else {
		state_status = `SELECT Status FROM Reaction WHERE User_UUID = ? AND Comment_ID = ?`
	}
	err_status := db.QueryRow(state_status, uuid, id).Scan(&status_comment)
	if err_status == sql.ErrNoRows {
		status_comment = ""
	} else if err_status != nil {
		ErrorServer(w, "Error checking post status")
	}
	return status_comment
}

// ? Function to get Username from post/comment creator
func CheckUsername(w http.ResponseWriter, uuid string) string {
	var username string
	// Getting the Username of the person who made the post
	state_username := `SELECT Username FROM User WHERE UUID = ?`
	err_db := db.QueryRow(state_username, uuid).Scan(&username)
	if err_db != nil {
		ErrorServer(w, "Error getting User's Username for the post")
	}
	return username
}
