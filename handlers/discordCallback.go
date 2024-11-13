package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/cmd/lib"
	"forum/models"
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
		// Erreur critique : Invalid state parameter
		err := &models.CustomError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid state parameter",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	_, err_db := db.Exec("DELETE FROM oauth_states WHERE state = ?", state)
	if err_db != nil {
		// Erreur non critique : Error deleting state
		ErrorServer(w, "Error deleting state, please try again later.")
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
		// Erreur critique : Error exchanging code
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error exchanging code",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		// Erreur critique : Error decoding token response
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error decoding token response",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		// Erreur critique : Invalid access token
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid access token",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	userInfoURL := "https://discord.com/api/users/@me"
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err_client_var := client.Do(req)
	if err_client_var != nil {
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

	discordID, idOk := userInfo["id"].(string)
	email, emailOk := userInfo["email"].(string)

	if !idOk || !emailOk {
		// Erreur critique : Missing user information
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Missing user information",
		}
		HandleError(w, err.StatusCode, err.Message)
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
		if err_password != nil {
			// Erreur critique : Error creating Password
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
			ErrorServer(w, "Error sending email, please check your inbox.")
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

		// Exec function to insert a new Users with all the data we got from the token
		result, err_doesnt_exist := db.Exec(`
			INSERT INTO User (UUID, Email, Username, Password, OAuthID, Role, IsDeleted) 
			VALUES (?, ?, ?, ?, ?, ?, false)
		`, UUID, email, username, password, discordID, "User")
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
	} else if err_db_check != nil {
		// Erreur critique : Database error
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database error",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	} else {
		// User exists, update GoogleID if necessary
		_, err_exist := db.Exec("UPDATE User SET OAuthID = ? WHERE ID = ?", discordID, userID)
		if err_exist != nil {
			// Erreur non critique : Error updating user
			ErrorServer(w, "Error updating user, please try again later.")
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
		// Erreur critique : Error accessing User UUID
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error accessing User UUID",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Attribute a session to an User
	CookieSession(user_uuid, w, r)

	data, err_getdata := lib.GetData(db, user_uuid, "logged", "index")
	if err_getdata != "OK" {
		// Erreur non critique : Unable to retrieve user data
		ErrorServer(w, "Unable to retrieve user data, please try again later.")
	}

	// Redirect the user to a success page or your main application
	renderTemplate(w, "layout/default", "page/index", data)

}
