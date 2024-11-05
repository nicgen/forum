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
)

func DiscordCallbackHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Verify the state to prevent CSRF attacks
	var storedState string
	err_check_db := db.QueryRow("SELECT state FROM oauth_states WHERE state = ?", state).Scan(&storedState)
	if err_check_db != nil || state != storedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	_, err_db := db.Exec("DELETE FROM oauth_states WHERE state = ?", state)
	if err_db != nil {
		log.Printf("Error deleting state: %v", err_db)
	}

	tokenURL := "https://discord.com/api/oauth2/token"
	values := url.Values{
		"client_id":     {os.Getenv("DISCORD_CLIENT_ID")},
		"client_secret": {os.Getenv("DISCORD_CLIENT_SECRET")},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {DiscordRedirectURI},
	}

	resp, err_post := http.PostForm(tokenURL, values)
	if err_post != nil {
		http.Error(w, "Error exchanging code", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		http.Error(w, "Error decoding token response", http.StatusInternalServerError)
		return
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		http.Error(w, "Invalid access token", http.StatusInternalServerError)
		return
	}

	userInfoURL := "https://discord.com/api/users/@me"
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err_client_var := client.Do(req)
	if err_client_var != nil {
		http.Error(w, "Error fetching user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Error decoding user info", http.StatusInternalServerError)
		return
	}

	discordID, idOk := userInfo["id"].(string)
	email, emailOk := userInfo["email"].(string)

	if !idOk || !emailOk {
		http.Error(w, "Missing user information", http.StatusInternalServerError)
		return
	}

	var userID int64

	// Checking if the email user in google authentication is already in the database or not
	err_db_check := db.QueryRow("SELECT ID FROM User WHERE OAuthID = ? OR Email = ?", discordID, email).Scan(&userID)

	// Checking if the user already exist in the database
	if err_db_check == sql.ErrNoRows {

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
			INSERT INTO User (UUID, Email, Username, Password, OAuthID, Role, IsDeleted) 
			VALUES (?, ?, ?, ?, ?, ?, false)
		`, UUID, email, username, password, discordID, "User")
		if err_doesnt_exist != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		userID, _ = result.LastInsertId()
	} else if err_db_check != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	} else {
		// User exists, update GoogleID if necessary
		_, err_exist := db.Exec("UPDATE User SET OAuthID = ? WHERE ID = ?", discordID, userID)
		if err_exist != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
	}

	//Checking if we got the user informations
	println("-------------------------------")
	println("User email: ", email)
	println("Discord ID: ", discordID)
	println("-------------------------------")

	// Getting the UUID from the database
	var user_uuid string
	state_uuid := `SELECT UUID FROM User WHERE Email = ?`
	err_user := db.QueryRow(state_uuid, email).Scan(&user_uuid)
	if err_user != nil {
		http.Error(w, "Error accessing User UUID", http.StatusUnauthorized)
		return
	}

	// Attribute a session to an User
	CookieSession(user_uuid, w, r)

	http.Redirect(w, r, "/", http.StatusFound)
}
