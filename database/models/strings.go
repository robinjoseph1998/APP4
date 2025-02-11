package models

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProfileResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	AccountType string `json:"account_type"`
	MediaCount  int    `json:"media_count"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	UserName    string `json:"user_name"`
	ExpiresIn   int    `json:"expires_in"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateCodeVerifier() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}

func GetUserIDFromContext(c *gin.Context) (int, error) {
	userIdStr := c.GetString("userID")
	userID, err := strconv.Atoi(userIdStr)
	return userID, err
}
