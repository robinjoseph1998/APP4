package repository

import "APP4/database/models"

type RepoInterfaces interface {
	CreateUser(user models.User) error
	SaveInstagramToken(userID uint, token string) error
	FetchTwitterAccessTokenFromDB(userName string) (string, error)
	SaveTwitterAccount(userId uint, accountName string, accessToken string, expiresIn int) error
	GetUserByEmail(loginRequest models.LoginRequest) (*models.User, error)
	FetchMyAccounts(userId uint) (*models.TwitterAccounts, error)
	DeleteAccountByName(accountName string) error
}
