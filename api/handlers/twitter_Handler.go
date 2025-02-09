package api

import (
	"APP4/api/repository"
	"APP4/database/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type OauthTwitterHandlers struct {
	Repo repository.RepoInterfaces
}

func NewOAuthTwitterHandlers(repo repository.RepoInterfaces) *OauthTwitterHandlers {
	return &OauthTwitterHandlers{
		Repo: repo}
}

var (
	CodeVerifier  = models.GenerateCodeVerifier()
	contextUserID = 0
)

func (x *OauthTwitterHandlers) OAuthTwitterLogin(c *gin.Context) {
	authURL := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=randomstring&code_challenge=%s&code_challenge_method=plain",
		os.Getenv("TWITTER_AUTH_URL"),
		os.Getenv("TWITTER_CLIENT_ID"),
		os.Getenv("TWITTER_REDIRECT_URI"),
		url.QueryEscape("tweet.read tweet.write users.read"), // Remove "offline.access" and "media.upload"
		CodeVerifier,
	)

	fmt.Println("🔗 Twitter Login URL:", authURL)
	c.Redirect(http.StatusFound, authURL)
}

func (x *OauthTwitterHandlers) OAuthTwitterCallback(c *gin.Context) {
	fmt.Println("🔄 Callback function triggered!")
	fmt.Println("Full Callback URL:", c.Request.URL.String())

	// Log all query parameters
	fmt.Println("Query Parameters:", c.Request.URL.Query())
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is missing"})
		return
	}

	fmt.Println("✅ Received Authorization Code:", code)

	data := url.Values{}
	data.Set("client_id", os.Getenv("TWITTER_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("TWITTER_CLIENT_SECRET"))
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", os.Getenv("TWITTER_REDIRECT_URI"))
	data.Set("code_verifier", CodeVerifier) // Ensure this is set correctly

	// Encode client_id and client_secret as Base64 for the Basic Auth header
	authHeader := base64.StdEncoding.EncodeToString([]byte(
		os.Getenv("TWITTER_CLIENT_ID") + ":" + os.Getenv("TWITTER_CLIENT_SECRET"),
	))

	req, _ := http.NewRequest("POST", os.Getenv("TWITTER_TOKEN_URL"), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+authHeader) // ✅ FIXED: Add Authorization header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Twitter Response:", string(body)) // Debugging purpose

	var tokenResponse models.TwitterTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	// Save token to database
	fmt.Println("ACCESS _______ TOKEN__", tokenResponse.AccessToken)
	err = x.Repo.SaveTwitterToken(tokenResponse.AccessToken, tokenResponse.RefreshToken, tokenResponse.ExpiresIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "access_token": tokenResponse.AccessToken})
}

func (x *OauthTwitterHandlers) PostTweet(c *gin.Context) {
	var requestBody struct {
		TweetText string `json:"tweet_text"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Twitter User Id From Twitter Handler: ", contextUserID)
	AccessToken, err := x.Repo.FetchAccessTokenFromDB(uint(contextUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't fetch access token"})
		return
	}
	// Step 2: Prepare the Twitter API request
	url := "https://api.twitter.com/2/tweets"
	payload := map[string]interface{}{
		"text": requestBody.TweetText,
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+AccessToken) // ✅ Get access token from Postman
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post tweet"})
		return
	}
	defer resp.Body.Close()

	// Step 3: Read Twitter API response
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Twitter Response:", string(body))

	c.JSON(http.StatusOK, gin.H{"message": "Tweet posted successfully!", "twitter_response": string(body)})
}
