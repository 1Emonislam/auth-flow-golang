package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

var googleOauthConfig *oauth2.Config

func init() {
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleLogin initiates the Google OAuth flow
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the callback from Google OAuth
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Get user info from Google
	client := googleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Here you would typically:
	// 1. Check if user exists in your database
	// 2. If not, create new user
	// 3. Generate JWT token
	// 4. Return token to client

	// For now, just return the user info
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully authenticated with Google",
		"user":    userInfo,
	})
}
