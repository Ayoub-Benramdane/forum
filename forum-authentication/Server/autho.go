package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	structs "forum/Data"
)

var Google = structs.ServiceAuth{
	ClientID:     "1029413446644-869kbbrp8h9na82m119n67mh06sqeh20.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURI:  "http://localhost:4444/callback/auth",
	AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
	TokenURL:     "https://oauth2.googleapis.com/token",
	UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
}

var Github = structs.ServiceAuth{
	ClientID:     "Ov23livNQhGyN67zpsmY",
	ClientSecret: "747f43600c227f9b93b41e956b360a4bfcdfde9b",
	RedirectURI:  "http://localhost:4444/callback/auth",
	AuthURL:      "https://github.com/login/oauth/authorize",
	TokenURL:     "https://github.com/login/oauth/access_token",
	UserInfoURL:  "https://api.github.com/user",
}

var Service structs.ServiceAuth

// handleGoogleLogin redirects the user to Google's OAuth 2.0 authorization endpoint.
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("service") == "github" {
		Service = Github
	} else {
		Service = Google
	}
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=email&prompt=consent",
		Service.AuthURL, Service.ClientID, Service.RedirectURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback processes the OAuth callback from Google.
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization code from the query parameters.
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token.
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", Service.ClientID)
	data.Set("client_secret", Service.ClientSecret)
	data.Set("redirect_uri", Service.RedirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(Service.TokenURL, data)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var accessToken string
	if Service == Github {
		body, _ := io.ReadAll(resp.Body)
		body1 := string(body)
		b, _ := url.ParseQuery(body1)
		accessToken = b.Get("access_token")
	} else {
		var tokenData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
			http.Error(w, "Failed to parse token response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		accessToken1, ok := tokenData["access_token"].(string)
		if !ok {
			http.Error(w, "Failed to retrieve access token", http.StatusInternalServerError)
			return
		}
		accessToken = accessToken1
	}
	defer resp.Body.Close()

	// Fetch the user's profile information.
	req, _ := http.NewRequest("GET", Service.UserInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse and display the user profile.
	var UserInfo structs.GoogleInfoUser
	var UserInfo1 structs.GithubInfoUser
	userInfo, _ := io.ReadAll(resp.Body)
	if Service == Google {
		if err = json.Unmarshal(userInfo, &UserInfo); err != nil {
			Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
			return
		}
		username := strings.Split(UserInfo.Email, "@")[0]
		RegisterPostAuth(username, UserInfo.Email, UserInfo.Id)
		LoginPostAuth(w, r, username, UserInfo.Id)
	} else {
		if err = json.Unmarshal(userInfo, &UserInfo1); err != nil {
			Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
			return
		}
		password := strconv.Itoa(UserInfo1.Id)
		RegisterPostAuth(UserInfo1.Login, UserInfo1.Login, password)
		LoginPostAuth(w, r, UserInfo1.Login, password)
	}
}
