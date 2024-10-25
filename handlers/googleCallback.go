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

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
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
	tokenURL := "https://oauth2.googleapis.com/token"

	// Set-up the request parameters
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", GoogleRedirectURI)
	values.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	values.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))

	// Sending the POST request
	resp, err_request := http.PostForm(tokenURL, values)
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

	// Use the access token to make API requests to Google's userinfo endpoint
	userInfoURL := "https://openidconnect.googleapis.com/v1/userinfo"
	req, _ := http.NewRequest("GET", userInfoURL, nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
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

	googleID, id_error := userInfo["sub"].(string)
	email, email_error := userInfo["email"].(string)

	if !id_error {
		http.Error(w, "Unable to get Google ID", http.StatusInternalServerError)
		return
	}
	if !email_error {
		http.Error(w, "Unable to get user email", http.StatusInternalServerError)
		return
	}

	var userID int64

	// Checking if the email user in google authentication is already in the database or not
	err = db.QueryRow("SELECT ID FROM User WHERE OAuthID = ? OR Email = ?", googleID, email).Scan(&userID)

	// Checking if the user already exist in the database
	if err == sql.ErrNoRows {

		// Creating an username with the content in front of the @
		username := email[:strings.Index(email, "@")]

		// If the user login for the first time we generate a password for him
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

		// Exec function to insert a new Users with all the data we got from the token
		result, err_doesnt_exist := db.Exec(`
			INSERT INTO User (UUID, Email, Username, Password, OAuthID, IsSuperUser, IsModerator, IsDeleted) 
			VALUES (?, ?, ?, ?, ?, false, false, false)
		`, UUID, email, username, password, googleID)
		if err_doesnt_exist != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		userID, _ = result.LastInsertId()
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	} else {
		// User exists, update GoogleID if necessary
		_, err_exist := db.Exec("UPDATE User SET OAuthID = ? WHERE ID = ?", googleID, userID)
		if err_exist != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
	}

	//Checking if we got the user informations
	print("-------------------------------\n")
	print("User email: ", email, "\n")
	print("Google ID: ", googleID, "\n")
	print("-------------------------------\n")

	// Redirect the user to a success page or your main application
	http.Redirect(w, r, "/", http.StatusFound)
}
