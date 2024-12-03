package lib

import (
	"database/sql"
	"forum/models" // Assurez-vous d'importer le package "models"
	"net/http"
	"strings"
)

func GetNotifications(w http.ResponseWriter, userUUID string, data map[string]interface{}) map[string]interface{} {
	var notifications []*models.Notification
	var totalCount int

	// Requête pour compter le nombre total de notifications pour l'utilisateur
	countQuery := `SELECT COUNT(*) FROM Notification WHERE User_UUID = ?`
	err := db.QueryRow(countQuery, userUUID).Scan(&totalCount)
	if err != nil {
		ErrorServer(w, "Error counting notifications")
		return data
	}

	// Requête pour récupérer les notifications avec distinction post/commentaire
	notificationsQuery := `
        SELECT 
            ID, 
            Post_ID, 
            Reaction_ID, 
            Comment_ID,
            CASE WHEN Comment_ID IS NOT NULL THEN true ELSE false END AS IsOnComment,
            CreatedAt, 
            IsRead 
        FROM Notification 
        WHERE User_UUID = ? 
        ORDER BY CreatedAt DESC
    `
	rows, err := db.Query(notificationsQuery, userUUID)
	if err != nil {
		ErrorServer(w, "Error accessing notifications")
		return data
	}
	defer rows.Close()

	// Parcourir les résultats
	for rows.Next() {
		var notification models.Notification
		var tempReactionID sql.NullInt64
		var tempCommentID sql.NullInt64
		var isOnComment bool

		// Scanner les données
		if err := rows.Scan(&notification.ID, &notification.PostID, &tempReactionID, &tempCommentID, &isOnComment, &notification.CreatedAt, &notification.IsRead); err != nil {
			ErrorServer(w, "Error scanning notifications")
			return data
		}

		// Gérer les valeurs NULL
		if tempReactionID.Valid {
			id := int(tempReactionID.Int64)
			notification.ReactionID = &id
		}
		if tempCommentID.Valid {
			id := int(tempCommentID.Int64)
			notification.CommentID = &id
		}

		time_comment := strings.Split(notification.CreatedAt.Format("2006-01-02 15:04:05"), " ")
		notification.Creation_Date = time_comment[0]
		notification.Creation_Hour = time_comment[1]

		// Affecter IsOnComment
		notification.IsOnComment = isOnComment

		// Ajouter la notification à la liste
		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		ErrorServer(w, "Error iterating notifications")
		return data
	}

	// Ajouter les notifications et le total à la map des données
	data["Notifications"] = notifications
	data["TotalCount"] = totalCount

	return data
}
