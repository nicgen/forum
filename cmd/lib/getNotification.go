package lib

import (
	"forum/models"
	"net/http"
)

// GetNotifications - Function to retrieve notifications for a user and store them in the data map
func GetNotifications(w http.ResponseWriter, userUUID string, data map[string]interface{}) map[string]interface{} {
	// Retrieve all notifications for the given user
	var notifications []*models.Notification

	notificationsQuery := `
		SELECT ID, Post_ID, Comment_ID, CreatedAt, IsRead	FROM Notification WHERE User_UUID = ? ORDER BY CreatedAt DESC`

	rows, err := db.Query(notificationsQuery, userUUID)
	if err != nil {
		ErrorServer(w, "Error accessing notifications")
		return data
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.Notification
		// var createdAt string

		// Scan the results into the Notification struct
		if err := rows.Scan(&notification.ID, &notification.PostID, &notification.CommentID, &notification.CreatedAt, &notification.IsRead); err != nil {
			ErrorServer(w, "Error scanning notifications")
			return data
		}

		// // Convert the createdAt string to a time.Time
		// notification.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		// if err != nil {
		// 	ErrorServer(w, "Error parsing notification timestamp")
		// 	return data
		// }

		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		ErrorServer(w, "Error iterating over notifications")
		return data
	}

	// Add notifications to the data map
	data["Notifications"] = notifications
	return data
}
