package routes

import (
	api "APP4/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine, ctrlInstagram *api.OauthInstagramHandlers, ctrlTwitter *api.OauthTwitterHandlers, ctrlCommon *api.CommonAuthHandlers) {
	router.POST("/app/signup", ctrlCommon.AppSignup)
	router.POST("/app/login", ctrlCommon.AppLogin)

	router.GET("/show/twitter/accounts", ctrlCommon.ShowConnectedTwitterAccounts)

	router.GET("/instagram/login", ctrlInstagram.OauthInstagramLogin)
	router.GET("/instagram/callback", ctrlInstagram.OauthInstagramCallback)
	router.GET("/fetch/instagram/profile", ctrlInstagram.FetchInstagramProfile)

	router.GET("/twitter/login", ctrlTwitter.OAuthTwitterLogin)
	router.GET("/twitter/callback", ctrlTwitter.OAuthTwitterCallback)

	router.POST("/twitter/post/tweet", ctrlTwitter.PostTweet)
	router.POST("/twitter/publish/video", ctrlTwitter.PostTweetWithVideo)
	router.GET("/fetch/twitter/profile", ctrlTwitter.FetchTwitterProfile)
	router.POST("/remove/account", ctrlCommon.RemoveAccounts)

}
