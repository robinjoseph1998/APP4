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

func (r *Repo) SaveInstagramToken(userID uint, token string) error {
	expiryTime := time.Now().Add(time.Duration(24) * time.Hour)
	query := "INSERT INTO connected_accounts (user_id,access_token, expires_at) VALUES (?, ?, ?, ?)"
	result := r.db.Exec(query, userID, token, expiryTime)
	if result.Error != nil {
		return fmt.Errorf("failed to save Twitter token: %v", result.Error)
	}
	return nil
}

func (r *Repo) FetchAccessTokenFromDB(userID uint, AccountName string) (string, error) {
	var token string
	query := "SELECT access_token FROM connected_accounts WHERE user_id = ? AND account_name=?"
	if err := r.db.Raw(query, userID, AccountName).Row().Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}

func (r *Repo) SaveTwitterAccount(userId uint, accountName string, accessToken string, expiresIn int) error {
	expiryTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
	query := "INSERT INTO connected_accounts (user_id, account_name, access_token, expires_at) VALUES (?, ?, ?, ?)"
	result := r.db.Exec(query, userId, accountName, accessToken, expiryTime)
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

func (r *Repo) FetchMyAccounts(userId uint) (*models.ConnectedAccounts, error) {
	var connectedAccounts *models.ConnectedAccounts
	query := "SELECT * FROM connected_accounts WHERE user_id = ?"
	if err := r.db.Raw(query, userId).Scan(&connectedAccounts).Error; err != nil {
		return nil, err
	}
	return connectedAccounts, nil
}
