package routes

import (
	api "APP4/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine, ctrlInstagram *api.OauthInstagramHandlers, ctrlTwitter *api.OauthTwitterHandlers) {
	router.POST("create/user", ctrlInstagram.AddUser)

	router.GET("instagram/login", ctrlInstagram.OauthInstagramLogin)
	router.GET("instagram/callback", ctrlInstagram.OauthInstagramCallback)
	router.GET("instagram/profile", ctrlInstagram.FetchInstagramProfile)

	router.GET("twitter/login", ctrlTwitter.OAuthTwitterLogin)
	router.GET("twitter/callback", ctrlTwitter.OAuthTwitterCallback)

	router.POST("/tweet", ctrlTwitter.PostTweet)

}
