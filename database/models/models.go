package models

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConnectedAccounts struct {
	AccountID   uint      `json:"account_id" gorm:"primaryKey"`
	UserID      int       `json:"user_id"`
	AccountName string    `json:"account_name"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}
