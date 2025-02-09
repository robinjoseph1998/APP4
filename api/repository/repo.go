package repository

import (
	"APP4/database/models"
	"fmt"
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
	instaToken := models.InstagramAccessToken{
		UserID: userID,
		Token:  token,
	}
	err := r.db.Create(&instaToken).Error
	if err != nil {
		return 0, "", err
	}
	return instaToken.UserID, instaToken.Token, nil
}

func (r *Repo) FetchAccessTokenFromDB(userID uint) (string, error) {
	var token string
	query := "SELECT access_token FROM twitter_access_tokens WHERE user_id = ?"
	if err := r.db.Raw(query, userID).Row().Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}

func (r *Repo) SaveTwitterToken(accessToken string, refreshToken string, expiresIn int) error {
	expiryTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
	query := `INSERT INTO twitter_access_tokens (access_token, refresh_token, expires_at) VALUES (?, ?, ?)`
	result := r.db.Exec(query, accessToken, refreshToken, expiryTime)
	if result.Error != nil {
		return fmt.Errorf("failed to save Twitter token: %v", result.Error)
	}
	return nil
}

func (r *Repo) GetUserByEmail(loginRequest models.LoginRequest) (*models.User, error) {
	var userDetails *models.User
	query := "SELECT * FROM users WHERE email = ? AND password = ?"
	if err := r.db.Raw(query, loginRequest.Email, loginRequest.Password).Scan(&userDetails).Error; err != nil {
		return nil, err
	}

	return userDetails, nil
}
