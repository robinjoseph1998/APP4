package api

import (
	"APP4/api/repository"
	services "APP4/api/services/twitter"
	"APP4/database/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type OauthTwitterHandlers struct {
	Repo     repository.RepoInterfaces
	services services.TwitterServiceInterfaces
}

func NewOAuthTwitterHandlers(repo repository.RepoInterfaces, service services.TwitterServiceInterfaces) *OauthTwitterHandlers {
	return &OauthTwitterHandlers{
		Repo:     repo,
		services: service}
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
		url.QueryEscape("tweet.read tweet.write users.read"), // Removing "offline.access" and "media.upload"
		CodeVerifier,
	)
	fmt.Println("Auth URL", authURL)
	c.Redirect(http.StatusFound, authURL)
}

func (x *OauthTwitterHandlers) OAuthTwitterCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		ErrorResponse(c, http.StatusBadRequest, "Authorization code is missing", nil)
		return
	}

	data := url.Values{}
	data.Set("client_id", os.Getenv("TWITTER_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("TWITTER_CLIENT_SECRET"))
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", os.Getenv("TWITTER_REDIRECT_URI"))
	data.Set("code_verifier", CodeVerifier)

	authHeader := base64.StdEncoding.EncodeToString([]byte(
		os.Getenv("TWITTER_CLIENT_ID") + ":" + os.Getenv("TWITTER_CLIENT_SECRET"),
	))

	req, _ := http.NewRequest("POST", os.Getenv("TWITTER_TOKEN_URL"), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to exchange token", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var tokenResponse models.TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to parse token response", err)
		return
	}
	err = x.Repo.SaveTwitterAccount(uint(contextUserID), "twitter", tokenResponse.AccessToken, tokenResponse.ExpiresIn)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to save token", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (x *OauthTwitterHandlers) PostTweet(c *gin.Context) {
	var requestBody struct {
		TweetText string `json:"tweet_text"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	AccessToken, err := x.Repo.FetchAccessTokenFromDB(uint(contextUserID), "twitter")
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "can't fetch access token", err)
		return
	}
	url := "https://api.twitter.com/2/tweets"
	payload := map[string]interface{}{
		"text": requestBody.TweetText,
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to post tweet", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	c.JSON(http.StatusOK, gin.H{"message": "Tweet posted successfully!", "twitter_response": string(body)})
}

func (x *OauthTwitterHandlers) PostTweetWithVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Failed to retrieve video file", err)
		return
	}
	tweetText := c.PostForm("tweet_text")
	if tweetText == "" {
		ErrorResponse(c, http.StatusBadRequest, "Tweet text is required", nil)
		return
	}
	saveDir := "uploads/videos"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	filePath := filepath.Join(saveDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to save video", err)
		return
	}
	mediaID, err := x.services.InitializeMediaUpload(filePath)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize media upload", err)
		return
	}

	err = x.services.AppendMediaUpload(mediaID, filePath)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to upload media chunks", err)
		return
	}

	err = x.services.FinalizeMediaUpload(mediaID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to finalize media upload", err)
		return
	}

	err = x.services.PostTweet(tweetText, mediaID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to post tweet", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tweet with video posted successfully!"})
}
