package lib

import (
	"database/sql"
	"fmt"
	"forum/models" // Assurez-vous d'importer le package "models"
	"net/http"
)

func GetNotifications(w http.ResponseWriter, userUUID string, useReactionID bool, data map[string]interface{}) map[string]interface{} {
	var notifications []*models.Notification
	var totalCount int

	// Requête pour compter le nombre total de notifications pour l'utilisateur
	countQuery := `SELECT COUNT(*) FROM Notification WHERE User_UUID = ?`
	err := db.QueryRow(countQuery, userUUID).Scan(&totalCount)
	if err != nil {
		ErrorServer(w, "Error counting notifications")
		return data
	}

	// Construire la requête SQL pour récupérer les notifications
	columnToSelect := "Reaction_ID"
	if !useReactionID {
		columnToSelect = "Comment_ID"
	}

	notificationsQuery := fmt.Sprintf(`
        SELECT ID, Post_ID, %s, CreatedAt, IsRead FROM Notification WHERE User_UUID = ? ORDER BY CreatedAt DESC`, columnToSelect)

	// Exécuter la requête pour récupérer les notifications
	rows, err := db.Query(notificationsQuery, userUUID)
	if err != nil {
		ErrorServer(w, "Error accessing notifications")
		return data
	}
	defer rows.Close()

	// Récupérer les notifications
	for rows.Next() {
		var notification models.Notification
		var tempID sql.NullInt64 // Utilisation de NullInt64 pour gérer les NULL

		// Scanner les données de la base
		if err := rows.Scan(&notification.ID, &notification.PostID, &tempID, &notification.CreatedAt, &notification.IsRead); err != nil {
			ErrorServer(w, "Error scanning notifications")
			return data
		}

		// Affecter ReactionID ou CommentID
		if tempID.Valid {
			id := int(tempID.Int64)
			if useReactionID {
				notification.ReactionID = &id
			} else {
				notification.CommentID = &id
			}
		}

		// Ajouter la notification à la liste
		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		ErrorServer(w, "Error iterating notifications")
		return data
	}

	// Ajouter les notifications et le total à la map des données
	data["Notifications"] = notifications
	data["TotalCount"] = totalCount // Assurez-vous que TotalCount est bien assigné

	return data
}
