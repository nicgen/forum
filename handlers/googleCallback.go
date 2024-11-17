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
	"time"

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
		// Erreur critique : Error exchanging code
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error exchanging code",
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

	// Use the access token to make API requests to Google's userinfo endpoint
	userInfoURL := "https://openidconnect.googleapis.com/v1/userinfo"
	req, errr := http.NewRequest("GET", userInfoURL, nil)
	if errr != nil {
		// Erreur critique : Error creating user info request
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error creating user info request",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
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
		// Erreur critique: Error decoding user info
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error decoding user info",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Setting up the variables we are going to set into cookies
	var username, creation_date, creation_hour string
	actual_time := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	creation_date = actual_time[0]
	creation_hour = actual_time[1]
	googleID, id_error := userInfo["sub"].(string)
	email, email_error := userInfo["email"].(string)
	username = email[:strings.Index(email, "@")]

	if !id_error {
		//Critical error : Unable to get Google ID
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unable to get Google ID",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	if !email_error {
		// Erreur non critique : Unable to get user email
		lib.ErrorServer(w, "Unable to get user email, some features may not work.")
	}

	var userID int64
	role := "User"

	// Checking if the email user in google authentication is already in the database or not
	err = db.QueryRow("SELECT ID, Role FROM User WHERE OAuthID = ? OR Email = ?", googleID, email).Scan(&userID, &role)

	// Checking if the user already exist in the database
	if err == sql.ErrNoRows {

		// Creating an username with the content in front of the @
		username := email[:strings.Index(email, "@")]

		// If the user login for the first time we generate a password for him
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
			lib.ErrorServer(w, "Error sending email, please check your inbox.")
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
			lib.ErrorServer(w, "Failed to generate UUID")
			return
		}

		// Exec function to insert a new Users with all the data we got from the token
		result, err_doesnt_exist := db.Exec(`
			INSERT INTO User (UUID, Email, Username, Password, OAuthID, Role, IsDeleted) 
			VALUES (?, ?, ?, ?, ?, ?, false)
		`, UUID, email, username, password, googleID, "User")
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
	} else if err != nil {
		// Erreur critique : Database error
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database error",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	} else {
		// User exists, update GoogleID if necessary
		_, err_exist := db.Exec("UPDATE User SET OAuthID = ? WHERE ID = ?", googleID, userID)
		if err_exist != nil {
			// Erreur non critique : Error updating user
			lib.ErrorServer(w, "Error updating user, please try again later.")
		}
	}

	//Checking if we got the user informations
	println("-------------------------------")
	println("User email: ", email)
	println("Google ID: ", googleID)
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
	lib.CookieSession(user_uuid, username, creation_date, creation_hour, email, role, w, r)

	// Redirect the user to a success page or your main application

	data, err_getdata := lib.GetData(db, user_uuid, "logged", "index", w, r)
	if err_getdata != "OK" {
		// Erreur critique : Error retrieving user data
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving user data",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// Redirect the user to a success page or your main application
	lib.RenderTemplate(w, "layout/index", "page/index", data)
}
