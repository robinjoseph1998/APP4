package api

import (
	"APP4/api/repository"
	"APP4/database/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OauthInstagramHandlers struct {
	Repo repository.RepoInterfaces
}

func NewOAuthInstagramHandlers(repo repository.RepoInterfaces) *OauthInstagramHandlers {
	return &OauthInstagramHandlers{
		Repo: repo}
}

func (ctrl *OauthInstagramHandlers) OauthInstagramLogin(c *gin.Context) {
	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")
	// Debugging: Print redirect URI in login request

	if clientID == "" || redirectURI == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing environment variables"})
		return
	}

	authURL := fmt.Sprintf(
		"https://www.facebook.com/v18.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=pages_show_list,instagram_basic,instagram_manage_insights&response_type=code",
		clientID, redirectURI,
	)
	fmt.Println("OAuth Callback Redirect URI:", strconv.Quote(redirectURI))

	c.JSON(http.StatusOK, gin.H{"login url": authURL})
}

func (ctrl *OauthInstagramHandlers) OauthInstagramCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No authorization code provided"})
		return
	}

	clientSecret := os.Getenv("INSTAGRAM_CLIENT_SECRET")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")
	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")

	// Exchange code for access token
	tokenURL := "https://api.instagram.com/oauth/access_token"
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", url.QueryEscape(redirectURI))

	data.Set("code", code)
	// Debugging: Print redirect URI
	fmt.Println("🔹 OAuth Callback Redirect URI:", strconv.Quote(redirectURI))

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)

	// Debugging: Print the full response from Instagram
	fmt.Println("Raw Instagram Response:", string(body))

	// Use a flexible map to parse JSON
	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to parse JSON",
			"response": string(body), // Debugging response
		})
		return
	}

	// Extract values safely
	accessToken, tokenOk := tokenResponse["access_token"].(string)
	userID, userIDOk := tokenResponse["user_id"].(float64) // JSON numbers default to float64

	if !tokenOk || !userIDOk {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Missing access token or user ID",
			"response": tokenResponse, // Debugging response
		})
		return
	}

	// Convert userID from float64 to uint
	savedUserID := uint(userID)

	// Save token in DB
	err = ctrl.Repo.SaveInstagramToken(savedUserID, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save access token"})
		return
	}

	// Return response with saved user ID and token
	c.JSON(http.StatusOK, gin.H{
		"message": "Access token saved successfully!",
	})

}

func (ctrl *OauthInstagramHandlers) FetchInstagramProfile(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	accessToken, err := ctrl.Repo.FetchAccessTokenFromDB(uint(userID), "instagram")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found"})
		return
	}

	profileURL := fmt.Sprintf("https://graph.instagram.com/me?fields=id,username,account_type,media_count&access_token=%s", accessToken)

	// Make GET request to Instagram API
	resp, err := http.Get(profileURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile details"})
		return
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)

	// Debugging: Print raw response
	fmt.Println("Instagram Profile Response:", string(body))

	// Parse JSON response
	var profile models.ProfileResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse profile data"})
		return
	}

	// Send response to client
	c.JSON(http.StatusOK, gin.H{
		"message":      "Profile fetched successfully!",
		"id":           profile.ID,
		"username":     profile.Username,
		"account_type": profile.AccountType,
		"media_count":  profile.MediaCount,
	})

}
