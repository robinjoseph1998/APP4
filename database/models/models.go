package models

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type InstagramToken struct {
	UserID uint   `json:"user_id" gorm:"primaryKey"`
	Token  string `json:"token"`
}
