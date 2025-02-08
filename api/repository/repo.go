package repository

import (
	"APP4/database/models"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) RepoInterfaces {
	return &Repo{db: db}
}

func (r *Repo) CreateUser(user models.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repo) SaveInstagramToken(userID uint, token string) (uint, string, error) {
	// Debugging: Print the token before saving
	fmt.Println("Saving Token for User ID:", userID)
	fmt.Println("Access Token:", token)

	instaToken := models.InstagramAccessToken{
		UserID: userID,
		Token:  token,
	}

	// Save token to DB
	err := r.db.Create(&instaToken).Error
	if err != nil {
		fmt.Println("Error Saving Token:", err)
		return 0, "", err
	}

	fmt.Println("✅ Access Token Saved Successfully!")
	return instaToken.UserID, instaToken.Token, nil
}

func (r *Repo) FetchAccessTokenFromDB(userID uint) (string, error) {
	token := ""
	if err := r.db.Find(&token, userID).Error; err != nil {
		return "", err
	}
	return token, nil
}

func (r *Repo) SaveTwitterToken(accessToken string, refreshToken string, expiresIn int) error {
	expiryTime := time.Now().Add(time.Duration(expiresIn) * time.Second) // Calculate expiry time

	query := `INSERT INTO twitter_access_tokens (access_token, refresh_token, expires_at) VALUES (?, ?, ?)`
	result := r.db.Exec(query, accessToken, refreshToken, expiryTime)
	if result.Error != nil {
		log.Println("Error saving Twitter token:", result.Error)
		return fmt.Errorf("failed to save Twitter token: %v", result.Error)
	}
	return nil
}
