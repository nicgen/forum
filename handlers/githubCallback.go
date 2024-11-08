package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/cmd/lib"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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
		ErrorServer(w, "Invalid state parameter")
	}

	// Delete the used state from the database
	_, err_delete_db := db.Exec("DELETE FROM oauth_states WHERE state = ?", state)
	if err_delete_db != nil {
		ErrorServer(w, "Error deleting state")
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
	req, _ := http.NewRequest("POST", tokenURL, strings.NewReader(values.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err_request := client.Do(req)
	if err_request != nil {
		ErrorServer(w, err_request.Error())
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var result map[string]interface{}
	if err_json := json.NewDecoder(resp.Body).Decode(&result); err_json != nil {
		ErrorServer(w, "Error decoding token response")
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		ErrorServer(w, "Unable to get access token")
	}

	// Use the access token to make API requests to GitHub's user endpoint
	userInfoURL := "https://api.github.com/user"
	req, _ = http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		ErrorServer(w, err.Error())
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		ErrorServer(w, "Error decoding user info")
	}

	githubID, id_error := userInfo["id"].(float64)
	email, email_error := userInfo["email"].(string)
	username, username_error := userInfo["login"].(string)

	if !id_error {
		ErrorServer(w, "Unable to get GitHub ID")
	}
	if !email_error {
		// If email is not public, fetch it separately
		email = fetchGitHubEmail(accessToken)
		if email == "" {
			ErrorServer(w, "Unable to get user email")
		}
	}
	if !username_error {
		ErrorServer(w, "Unable to get GitHub username")
	}

	var userID int64

	// Checking if the user is already in the database
	err_db := db.QueryRow("SELECT ID FROM User WHERE OAuthID = ? OR Email = ?", int64(githubID), email).Scan(&userID)

	// If the user doesn't exist, create a new one
	if err_db == sql.ErrNoRows {
		// Generate a password for the user
		password, err_password := lib.GeneratePassword(16)
		err_email := lib.SendEmail(email, "Your forum password", password)
		password, err_hashing := lib.HashPassword(password)

		// Checking for password errors
		if err_password != nil {
			ErrorServer(w, "Error creating password")
		} else if err_hashing != nil {
			ErrorServer(w, "Error hashing password")
		} else if err_email != nil {
			ErrorServer(w, "Error sending email")
		}

		//Generate Random UUID for the user
		UUID, err := uuid.NewV4()
		if err != nil {
			ErrorServer(w, "failed to generate UUID")
		}

		// Insert new user into the database
		result, err_doesnt_exist := db.Exec(`
			INSERT INTO User (UUID, Email, Username, Password, OAuthID, Role, IsDeleted) 
			VALUES (?, ?, ?, ?, ?, ?, false)
		`, UUID, email, username, password, int64(githubID), "User")
		if err_doesnt_exist != nil {
			ErrorServer(w, "Error creating user")
		}
		userID, _ = result.LastInsertId()
	} else if err_db != nil {
		ErrorServer(w, "Database error")
	} else {
		// User exists, update GitHubID if necessary
		_, err_exist := db.Exec("UPDATE User SET OAuthID = ? WHERE ID = ?", int64(githubID), userID)
		if err_exist != nil {
			ErrorServer(w, "Error updating user")
		}
	}

	//Checking if we got the user informations
	print("-------------------------------")
	print("User email: ", email)
	print("Github ID: ", githubID)
	print("-------------------------------")

	// Getting the UUID from the database
	var user_uuid string
	state_uuid := `SELECT UUID FROM User WHERE Email = ?`
	err_user := db.QueryRow(state_uuid, email).Scan(&user_uuid)
	if err_user != nil {
		ErrorServer(w, "Error accessing User UUID")
	}

	// Attribute a session to an User
	CookieSession(user_uuid, w, r)

	data, err_getdata := lib.GetData(db, user_uuid, "logged", "index")
	if err_getdata != "OK" {
		ErrorServer(w, err_getdata)
	}

	// Redirect the user to a success page or your main application
	renderTemplate(w, "layout/default", "page/index", data)
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
