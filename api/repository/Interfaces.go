package repository

import "APP4/database/models"

type RepoInterfaces interface {
	CreateUser(user models.User) error
}
