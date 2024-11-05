package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Delete the used state from the database
	_, err_delete_db := db.Exec("DELETE FROM oauth_states WHERE state = ?", state)
	if err_delete_db != nil {
		log.Printf("Error deleting state: %v", err_delete_db)
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
		http.Error(w, err_request.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var result map[string]interface{}
	if err_json := json.NewDecoder(resp.Body).Decode(&result); err_json != nil {
		http.Error(w, "Error decoding token response", http.StatusInternalServerError)
		return
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		http.Error(w, "Unable to get access token", http.StatusInternalServerError)
		return
	}

	// Use the access token to make API requests to GitHub's user endpoint
	userInfoURL := "https://api.github.com/user"
	req, _ = http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Error decoding user info", http.StatusInternalServerError)
		return
	}

	githubID, id_error := userInfo["id"].(float64)
	email, email_error := userInfo["email"].(string)
	username, username_error := userInfo["login"].(string)

	if !id_error {
		http.Error(w, "Unable to get GitHub ID", http.StatusInternalServerError)
		return
	}
	if !email_error {
		// If email is not public, fetch it separately
		email = fetchGitHubEmail(accessToken)
		if email == "" {
			http.Error(w, "Unable to get user email", http.StatusInternalServerError)
			return
		}
	}
	if !username_error {
		http.Error(w, "Unable to get GitHub username", http.StatusInternalServerError)
		return
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
			http.Error(w, "Error creating password", http.StatusInternalServerError)
			return
		} else if err_hashing != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		} else if err_email != nil {
			http.Error(w, "Error sending email", http.StatusInternalServerError)
			return
		}

		//Generate Random UUID for the user
		UUID, err := uuid.NewV4()
		if err != nil {
			fmt.Printf("failed to generate UUID: %v\n", err)
			return
		}

		// Insert new user into the database
		result, err_doesnt_exist := db.Exec(`
			INSERT INTO User (UUID, Email, Username, Password, OAuthID, Role, IsDeleted) 
			VALUES (?, ?, ?, ?, ?, false, ?, false)
		`, UUID, email, username, password, int64(githubID))
		if err_doesnt_exist != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		userID, _ = result.LastInsertId()
	} else if err_db != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	} else {
		// User exists, update GitHubID if necessary
		_, err_exist := db.Exec("UPDATE User SET OAuthID = ? WHERE ID = ?", int64(githubID), userID)
		if err_exist != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
	}

	//Checking if we got the user informations
	print("-------------------------------\n")
	print("User email: ", email, "\n")
	print("Github ID: ", githubID, "\n")
	print("-------------------------------\n")

	// Redirect the user to a success page or your main application
	http.Redirect(w, r, "/", http.StatusFound)
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
