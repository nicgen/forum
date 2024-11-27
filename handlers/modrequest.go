package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"html/template"
	"log"
	"net/http"
)

// Handles POST requests to update the `IsRequest` column for a specific user.
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST; reject any other methods.
	if r.Method != http.MethodPost {
		lib.ErrorServer(w, "Invalid request method.") // Custom error handling function
		return
	}

	// Retrieve the user's UUID from the form data.
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		// Erreur non critique : UUID required
		lib.ErrorServer(w, "User UUID is required.") // Respond if UUID is missing.
		return
	}

	// Connect to the database using a custom utility function.
	db := lib.GetDB()

	// Prepare and execute the SQL query to set `IsRequest` to TRUE.
	requestQuery := `UPDATE User SET IsRequest = TRUE WHERE UUID = ?`
	_, err := db.Exec(requestQuery, userUUID)
	if err != nil {
		//Erreur critique : échec mise a jour du statut de la demande
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error updating request status.Please try again later",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Redirect the user to their profile after successfully updating the status.
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Fetches the `IsRequest` status for a specific user and renders it on a webpage.
func GetIsRequestStatus(w http.ResponseWriter, r *http.Request) {
	// Extract the UUID from the form data.
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		//Erreur non critique : User UUID required
		http.Error(w, "User UUID is required.", http.StatusBadRequest) // Handle missing UUID.
		return
	}

	// Connect to the database.
	db := lib.GetDB()

	var isRequest bool // Variable to store the result of the query.

	// Query to retrieve the `IsRequest` status for the given UUID.
	query := `SELECT IsRequest FROM User WHERE UUID = ?`

	err := db.QueryRow(query, userUUID).Scan(&isRequest)
	if err != nil {
		if err == sql.ErrNoRows {
			//Erreur non critique : User not found
			http.Error(w, "User not found.", http.StatusNotFound)
		} else {
			//Erreur critique : échec de la récuperation du status IsRequest
			log.Printf("Error fetching IsRequest status for UUID %s: %v", userUUID, err)
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching data from the database. Please try again later.",
			}
			HandleError(w, err.StatusCode, err.Message)
		}
		return
	}

	// Parse the HTML template for rendering the response.
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		// Erreur critique : échec de l'analyse du modèle
		log.Printf("Error parsing template: %v", err)
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Create a user model with the `IsRequest` status.
	data := models.User{
		IsRequest: isRequest,
	}

	// Render the template with the user data and send it to the response.
	if err := tmpl.Execute(w, data); err != nil {
		// Erreur critique : échec de l'exécution du modèle
		log.Printf("Error executing template: %v", err)
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
	}
}
