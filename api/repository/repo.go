package repository

import (
	"APP4/database/models"
	"fmt"

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

	instaToken := models.InstagramToken{
		UserID: userID,
		Token:  token,
	}

	// Save token to DB
	err := r.db.Create(&instaToken).Error
	if err != nil {
		fmt.Println("❌ Error Saving Token:", err)
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
