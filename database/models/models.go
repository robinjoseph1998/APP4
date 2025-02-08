package models

import "time"

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type InstagramAccessToken struct {
	UserID uint   `json:"user_id" gorm:"primaryKey"`
	Token  string `json:"token"`
}

type TwitterAccessToken struct {
	UserID       uint      `json:"user_id" gorm:"primaryKey"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
