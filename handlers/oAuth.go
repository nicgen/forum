package handlers

import (
	"forum/cmd/lib"
	"net/http"
	"net/url"
	"os"
)

const (
	GoogleRedirectURI  = "https://localhost:8080/callback"
	githubRedirectURI  = "https://localhost:8080/github/callback"
	DiscordRedirectURI = "https://localhost:8080/discord/callback"
)

func GoogleOAuthHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Generate a random state parameter to prevent CSRF attacks
	state, err_state := lib.GenerateRandomID()
	if err_state != nil {
		ErrorServer(w, "Error generating state")
	}

	// Store the state in the database
	_, err := db.Exec("INSERT INTO oauth_states (state) VALUES (?)", state)
	if err != nil {
		ErrorServer(w, "Error storing state")
	}

	// Redirect the user to the Google OAuth 2.0 authorization URL
	authURL := "https://accounts.google.com/o/oauth2/auth"

	params := url.Values{
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"response_type": {"code"},
		"redirect_uri":  {GoogleRedirectURI},
		"scope":         {"email profile"},
		"state":         {state},
	}

	http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusFound)
}

func GitHubOAuthHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// Generate a random state parameter to prevent CSRF attacks
	state, err_state := lib.GenerateRandomID()
	if err_state != nil {
		ErrorServer(w, "Error generating state")
	}

	// Store the state in the database
	_, err := db.Exec("INSERT INTO oauth_states (state) VALUES (?)", state)
	if err != nil {
		ErrorServer(w, "Error storing state")
	}

	// Redirect the user to the GitHub OAuth 2.0 authorization URL
	authURL := "https://github.com/login/oauth/authorize"

	params := url.Values{
		"client_id":    {os.Getenv("GITHUB_CLIENT_ID")},
		"redirect_uri": {githubRedirectURI},
		"scope":        {"user:email"},
		"state":        {state},
	}

	http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusFound)
}

func DiscordOAuthHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	state, err := lib.GenerateRandomID()
	if err != nil {
		ErrorServer(w, "Error generating state")
	}

	_, err = db.Exec("INSERT INTO oauth_states (state) VALUES (?)", state)
	if err != nil {
		ErrorServer(w, "Error storing state")
	}

	authURL := "https://discord.com/api/oauth2/authorize"
	params := url.Values{
		"client_id":     {os.Getenv("DISCORD_CLIENT_ID")},
		"redirect_uri":  {DiscordRedirectURI},
		"response_type": {"code"},
		"scope":         {"identify email"},
		"state":         {state},
	}

	http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusFound)
}
