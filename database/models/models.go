package models

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TwitterAccounts struct {
	AccountID   uint      `json:"account_id" gorm:"primaryKey"`
	UserID      int       `json:"user_id"`
	UserName    string    `json:"user_name"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type InstagramAccounts struct {
	AccountID   uint      `json:"account_id" gorm:"primaryKey"`
	UserID      int       `json:"user_id"`
	BusinessID  string    `json:"business_id"`
	UserName    string    `json:"user_name"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}
