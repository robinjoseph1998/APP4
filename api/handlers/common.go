package api

import (
	"APP4/api/repository"
	"APP4/database/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommonAuthHandlers struct {
	Repo repository.RepoInterfaces
}

func NewCommonAuthHandlers(repo repository.RepoInterfaces) *CommonAuthHandlers {
	return &CommonAuthHandlers{
		Repo: repo}
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	errorResponse := gin.H{
		"error":   true,
		"message": message,
	}
	if err != nil {
		errorResponse["details"] = err.Error()
	}
	c.JSON(statusCode, errorResponse)
}

func (ctrl *CommonAuthHandlers) AppSignup(c *gin.Context) {
	var request models.User
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var loginRequest models.LoginRequest

	loginRequest.Email = request.Email
	loginRequest.Password = request.Password

	user, err := ctrl.Repo.GetUserByEmail(loginRequest)
	if err != nil || user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	if err := ctrl.Repo.CreateUser(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created successully"})
}

func (ctrl *CommonAuthHandlers) AppLogin(c *gin.Context) {
	var loginRequest models.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := ctrl.Repo.GetUserByEmail(loginRequest)
	if err != nil && user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if user != nil {
		// strUserId := strconv.Itoa(int(user.ID))
		// token, err := middleware.GenToken(strUserId, user.Email, c)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "token genration failed"})
		// 	return
		// }
		contextUserID = int(user.ID)
		fmt.Println("Twitter User Id from Common Handler", contextUserID)
		c.JSON(http.StatusOK, gin.H{
			"message": "Logged in successfully",
		})
	}
}

func (ctrl *CommonAuthHandlers) ShowConnectedAccounts(c *gin.Context) {
	userId := c.PostForm("user_id")
	if userId == "" {
		ErrorResponse(c, http.StatusBadRequest, "invalid user id", nil)
		return
	}
	uintUserId, err := strconv.Atoi(userId)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "user id conversion failed", nil)
		return
	}
	accountsConnected, err := ctrl.Repo.FetchMyAccounts(uint(uintUserId))
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "can't fetch connected acccount details", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Accounts": accountsConnected.AccountName})
}

func (ctrl *CommonAuthHandlers) RemoveAccounts(c *gin.Context) {

}
