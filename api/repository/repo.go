package repository

import (
	"APP4/database/models"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) RepoInterfaces {
	return &Repo{db: db}
}

func (r Repo) CreateUser(user models.User) error {
	return nil
}
