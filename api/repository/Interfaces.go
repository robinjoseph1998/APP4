package repository

import "APP4/database/models"

type RepoInterfaces interface {
	CreateUser(user models.User) error
	SaveInstagramToken(userID uint, token string) (uint, string, error)
	FetchAccessTokenFromDB(userID uint) (string, error)

	SaveTwitterToken(accessToken string, refreshToken string, expiresIn int) error
	GetUserByEmail(loginRequest models.LoginRequest) (*models.User, error)
}
