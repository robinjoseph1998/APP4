package api

import (
	"APP4/api/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

func (ctrl *OauthInstagramHandlers) OauthInstagramCallback(c *gin.Context) {
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
		TokenType   string `json:"token_type"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &AccessTokenResponse); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to parse token response", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": AccessTokenResponse.AccessToken})
}

func (ctrl *OauthInstagramHandlers) FetchInstagramProfile(c *gin.Context) {

	access_token := "EAA07LhpjPJIBOZCdOTte4OC33N4A9CSb21cVZAUocxBAZCXwz8mHVWzql1NVcmHLTMGFzmrUhj2v5zFrSShYPBZASjYliaZCWQWmjVm5ZByQs0vWvFKqZB9OLlVynhxyISRP93QSUag8TAdH31sciAwNTj1AZCEjXm2YlmfTXsVva91sEefpLKZASZCEjbvsbhUhJKZAD6xCg0gRjKzidEWZA47EZASz5dbyQaaPv1YytGbw3JVE3CbN9xDAXXFvwel7HGfI07wZDZD"
	if access_token == "" {
		ErrorResponse(c, http.StatusInternalServerError, "can't fetch access token in .env", nil)
		return
	}
	business_account_id := "17841472666905757"
	igProfileEndpoint := fmt.Sprintf("https://graph.facebook.com/v22.0/%s?fields=name,username&access_token=%s", business_account_id, access_token)
	req, err := http.NewRequest("GET", igProfileEndpoint, nil)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create request", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to send http request", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to read the response", err)
		return
	}
	fmt.Println("Raw Response:", string(body))
	var profileData map[string]interface{}
	if err := json.Unmarshal(body, &profileData); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to parse JSON response", err)
		return
	}
	c.JSON(http.StatusOK, profileData)

}
