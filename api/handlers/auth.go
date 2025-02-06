package api

import (
	"APP4/api/repository"

	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	Repo repository.RepoInterfaces
}

func NewAuthHandlers(repo repository.RepoInterfaces) *AuthHandlers {
	return &AuthHandlers{
		Repo: repo}
}

func (ctrl AuthHandlers) AddUser(c *gin.Context) {

}
