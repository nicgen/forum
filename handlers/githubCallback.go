package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/cmd/lib"
	"forum/models"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	_ "github.com/mattn/go-sqlite3"
)

func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Extract the authorization code and state from the query parameters
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Verify the state to prevent CSRF attacks
	var storedState string
	err_check_db := db.QueryRow("SELECT state FROM oauth_states WHERE state = ?", state).Scan(&storedState)
	if err_check_db != nil || state != storedState {
		// Erreur critique : Invalid state parameter
		err := &models.CustomError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid state parameter",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Delete the used state from the database
	_, err_delete_db := db.Exec("DELETE FROM oauth_states WHERE state = ?", state)
	if err_delete_db != nil {
		// Erreur critique : Error deleting state
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting state",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Exchange the authorization code for an access token
	tokenURL := "https://github.com/login/oauth/access_token"

	// Set-up the request parameters
	values := url.Values{}
	values.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	values.Set("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	values.Set("code", code)
	values.Set("redirect_uri", githubRedirectURI)

	// Sending the POST request
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		// Erreur critique : Error creating request
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error creating request",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err_request := client.Do(req)
	if err_request != nil {
		// Erreur critique : Error sending request
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error sending request",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var result map[string]interface{}
	if err_json := json.NewDecoder(resp.Body).Decode(&result); err_json != nil {
		// Erreur critique : Error decoding token response
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error decoding token response",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		// Erreur critique : Unable to get access token
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unable to get access token",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Use the access token to make API requests to GitHub's user endpoint
	userInfoURL := "https://api.github.com/user"
	req, err = http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		// Erreur critique : Error creating user info request
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error creating user info request",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		// Erreur critique : Error fetching user info
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching user info",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		// Erreur critique : Error decoding user info
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error decoding user info",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Setting up the variables we are going to set into cookies
	var creation_date, creation_hour string
	githubID, id_error := userInfo["id"].(float64)
	email, email_error := userInfo["email"].(string)
	username, username_error := userInfo["login"].(string)
	actual_time := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	creation_date = actual_time[0]
	creation_hour = actual_time[1]

	if !id_error {
		// Erreur critique : Unable to get GitHub ID
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unable to get GitHub ID",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	if !email_error {
		// If email is not public, fetch it separately
		email = fetchGitHubEmail(accessToken)
		if email == "" {
			// Erreur non critique : Unable to get user email
			lib.ErrorServer(w, "Unable to get user email")
		}
	}
	if !username_error {
		// Erreur critique : Unable to get GitHub username
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unable to get GitHub Username",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	var userID int64
	role := "User"

	// Checking if the user is already in the database
	err_db := db.QueryRow("SELECT ID, Role FROM User WHERE OAuthID = ? OR Email = ?", int64(githubID), email).Scan(&userID, &role)

	// If the user doesn't exist, create a new one
	if err_db == sql.ErrNoRows {
		// Generate a password for the user
		password, err_password := lib.GeneratePassword(16)
		if err_password != nil {
			// Erreur critique : Error creating password
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error creating password",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}
		err_email := lib.SendEmail(email, "Your forum password", password)
		if err_email != nil {
			// Erreur non critique : Error sending email
			lib.ErrorServer(w, "Error sending email")
		}

		password, err_hashing := lib.HashPassword(password)
		if err_hashing != nil {
			// Erreur critique : Error hashing password
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error hashing password",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}

		//Generate Random UUID for the user
		UUID, err := uuid.NewV4()
		if err != nil {
			// Erreur critique : Failed to generate UUID
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to generate UUID",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}

		// Insert new user into the database
		result, err_doesnt_exist := db.Exec(`
			INSERT INTO User (UUID, Email, Username, Password, OAuthID, Role, IsDeleted) 
			VALUES (?, ?, ?, ?, ?, ?, false)
		`, UUID, email, username, password, int64(githubID), "User")
		if err_doesnt_exist != nil {
			// Erreur critique : Error creating user
			err := &models.CustomError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error creating user",
			}
			HandleError(w, err.StatusCode, err.Message)
			return
		}
		userID, _ = result.LastInsertId()
	} else if err_db != nil {
		// Erreur critique : Error checking user in database
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error checking user in database",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	} else {
		// User exists, update GitHubID if necessary
		_, err_exist := db.Exec("UPDATE User SET OAuthID = ? WHERE ID = ?", int64(githubID), userID)
		if err_exist != nil {
			// Erreur non critique : Error updating user
			lib.ErrorServer(w, "Error updating user")
		}
	}

	//Checking if we got the user informations
	println("-------------------------------")
	println("User email: ", email)
	println("Github ID: ", githubID)
	println("-------------------------------")

	// Getting the UUID from the database
	var user_uuid string
	state_uuid := `SELECT UUID FROM User WHERE Email = ?`
	err_user := db.QueryRow(state_uuid, email).Scan(&user_uuid)
	if err_user != nil {
		// Erreur non critique : Error accessing User UUID
		lib.ErrorServer(w, "Error accessing User UUID")
	}

	// Attribute a session to an User
	lib.CookieSession(user_uuid, username, creation_date, creation_hour, email, role, w, r)

	data := lib.GetData(db, user_uuid, "logged", "index", w, r)

	// Redirect the user to a success page or your main application
	lib.RenderTemplate(w, "layout/index", "page/index", data)

}

// Helper function to fetch GitHub email if it's not public
func fetchGitHubEmail(accessToken string) string {
	emailURL := "https://api.github.com/user/emails"
	req, _ := http.NewRequest("GET", emailURL, nil)
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching GitHub email: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		log.Printf("Error decoding GitHub email response: %v", err)
		return ""
	}

	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email
		}
	}

	return ""
}
