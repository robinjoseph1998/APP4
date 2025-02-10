package api

import (
	"APP4/api/repository"
	"APP4/database/models"
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
		ErrorResponse(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	var loginRequest models.LoginRequest
	loginRequest.Email = request.Email
	loginRequest.Password = request.Password

	user, err := ctrl.Repo.GetUserByEmail(loginRequest)
	if err != nil || user != nil {
		ErrorResponse(c, http.StatusBadRequest, "user already exists", err)
		return
	}

	if err := ctrl.Repo.CreateUser(request); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "can't create the user", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created successully"})
}

func (ctrl *CommonAuthHandlers) AppLogin(c *gin.Context) {
	var loginRequest models.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "invalid request body", err)
		return
	}
	user, err := ctrl.Repo.GetUserByEmail(loginRequest)
	if err != nil && user == nil {
		ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", err)
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
		c.JSON(http.StatusOK, gin.H{
			"message": "Logged in successfully",
		})
	}
}

func (ctrl *CommonAuthHandlers) ShowConnectedTwitterAccounts(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"Twitter_Accounts": accountsConnected.UserName})
}

func (ctrl *CommonAuthHandlers) RemoveAccounts(c *gin.Context) {
	accountName := c.PostForm("account_name")
	if accountName == "" {
		ErrorResponse(c, http.StatusBadRequest, "invalid account name", nil)
		return
	}

	err := ctrl.Repo.DeleteAccountByName(accountName)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "account deletion failed", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Account Name": accountName, "Status": "Deleted Successfully"})
}
