package api

import (
	"APP4/api/repository"
	services "APP4/api/services/instagram"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OauthInstagramHandlers struct {
	Repo     repository.RepoInterfaces
	Services services.InstagramServiceInterfaces
}

func NewOAuthInstagramHandlers(repo repository.RepoInterfaces, services services.InstagramServiceInterfaces) *OauthInstagramHandlers {
	return &OauthInstagramHandlers{
		Repo:     repo,
		Services: services}
}

func (ig *OauthInstagramHandlers) OauthInstagramLogin(c *gin.Context) {
	clientID := os.Getenv("INSTAGRAM_APP_ID")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")

	if clientID == "" || redirectURI == "" {
		ErrorResponse(c, http.StatusInternalServerError, "Missing environment variables", nil)
		return
	}
	authURL := fmt.Sprintf(
		"https://www.facebook.com/v18.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=pages_show_list,instagram_basic,instagram_manage_insights&response_type=code",
		clientID, redirectURI,
	)
	fmt.Println("AuthUrl", authURL)
	c.Redirect(http.StatusFound, authURL)
}

func (ig *OauthInstagramHandlers) OauthInstagramCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		ErrorResponse(c, http.StatusBadRequest, "Authorization code is missing", nil)
		return
	}
	exchangeApiURL := "https://graph.facebook.com/v22.0/oauth/access_token"
	clientID := os.Getenv("INSTAGRAM_APP_ID")
	clientSecret := os.Getenv("INSTAGRAM_APP_SECRET")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")

	data := fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s",
		clientID, clientSecret, redirectURI, code,
	)
	resp, err := http.Post(exchangeApiURL, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(data)))
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to exchange token", err)
		return
	}
	defer resp.Body.Close()
	var AccessTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &AccessTokenResponse); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to parse token response", err)
		return
	}

	businessId, err := ig.Services.GetIGBusinessID(AccessTokenResponse.AccessToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "failed to get instagram business id", err)
	}

	strId := strconv.Itoa(contextUserID)
	err = ig.Repo.SaveInstgramAccount(contextUserID, "instagram"+strId, businessId, AccessTokenResponse.AccessToken, 5)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to save token", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instagram account connected to the app successfully"})
}

func (ig *OauthInstagramHandlers) FetchInstagramProfile(c *gin.Context) {
	userName := c.PostForm("user_name")
	if userName == "" {
		ErrorResponse(c, http.StatusBadRequest, "invalid Request", nil)
		return
	}
	accessToken, err := ig.Repo.FetchInstagramAccessTokenFromDB(userName)
	if err != nil || accessToken == "" {
		ErrorResponse(c, http.StatusInternalServerError, "can't fetch access token", err)
		return
	}
	businessId, err := ig.Services.GetIGBusinessID(accessToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "failed to get instagram business id", err)
	}

	iGprofileApi := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?fields=id,username,profile_picture_url&access_token=%s", businessId, accessToken)
	resp, err := http.Get(iGprofileApi)
	if err != nil || resp.StatusCode != http.StatusOK {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch profile", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to read the response", err)
		return
	}

	var profileData map[string]interface{}
	if err := json.Unmarshal(body, &profileData); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to parse JSON response", err)
		return
	}
	c.JSON(http.StatusOK, profileData)
}

func (ig *OauthInstagramHandlers) PostInstagramReel(c *gin.Context) {
	userName := c.PostForm("user_name")
	if userName == "" {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request: user_name is required", nil)
		return
	}

	postCaption := c.PostForm("caption")
	if postCaption == "" {
		ErrorResponse(c, http.StatusBadRequest, "Post description is required", nil)
		return
	}

	videoURL := c.PostForm("video_url")
	if videoURL == "" {
		ErrorResponse(c, http.StatusBadRequest, "Failed to retrieve media file", nil)
		return
	}

	accessToken, err := ig.Repo.FetchInstagramAccessTokenFromDB(userName)
	if err != nil || accessToken == "" {
		ErrorResponse(c, http.StatusInternalServerError, "Can't fetch access token", err)
		return
	}

	businessID, err := ig.Services.GetIGBusinessID(accessToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get Instagram Business ID", err)
		return
	}
	mediaID, err := ig.Services.UploadInstagramReel(businessID, videoURL, postCaption, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video", "details": err.Error()})
		return
	}

	err = ig.Services.CheckVideoProcessingStatus(mediaID, accessToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Video processing failed", err)
		return
	}

	postID, err := ig.Services.PublishInstagramVideo(businessID, mediaID, accessToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to publish media", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Media posted successfully!",
		"post_id": postID,
	})
}

func (ctrl *OauthInstagramHandlers) ShowConnectedInstagramAccounts(c *gin.Context) {
	userId := c.PostForm("user_id")
	if userId == "" {
		ErrorResponse(c, http.StatusBadRequest, "invalid user id", nil)
		return
	}
	accounts, err := ctrl.Repo.FetchMyInstagramAccounts(contextUserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "can't fetch connected acccount details", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Instagram_Accounts": accounts.UserName,
	})
}
