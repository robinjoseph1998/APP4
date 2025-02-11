package repository

import "APP4/database/models"

type RepoInterfaces interface {
	CreateUser(user models.User) error
	SaveInstagramToken(userID uint, token string) error
	FetchTwitterAccessTokenFromDB(userName string) (string, error)
	FetchInstagramAccessTokenFromDB(userName string) (string, error)
	SaveTwitterAccount(userId uint, userName string, accessToken string, expiresIn int) error
	GetUserByEmail(loginRequest models.LoginRequest) (*models.User, error)
	FetchMyTwitterAccounts(userId int) (*models.TwitterAccounts, error)
	DeleteTwitterAccount(userName string) error
	SaveInstgramAccount(userId int, userName string, businessId string, accessToken string, expiresIn int) error
	FetchMyInstagramAccounts(userId int) (*models.InstagramAccounts, error)
}
