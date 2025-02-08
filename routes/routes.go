package routes

import (
	api "APP4/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine, ctrl *api.AuthHandlers) {
	router.POST("create/user", ctrl.AddUser)

	router.GET("instagram/login", ctrl.OauthInstagramLogin)
	router.GET("instagram/callback", ctrl.OauthInstagramCallback)
	router.GET("instagram/profile", ctrl.FetchInstagramProfile)
}
