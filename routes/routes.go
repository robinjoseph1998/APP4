package routes

import (
	api "APP4/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine, ctrlInstagram *api.OauthInstagramHandlers, ctrlTwitter *api.OauthTwitterHandlers, ctrlCommon *api.CommonAuthHandlers) {
	router.POST("/app/signup", ctrlCommon.AppSignup)
	router.POST("/app/login", ctrlCommon.AppLogin)

	router.GET("/instagram/login", ctrlInstagram.OauthInstagramLogin)
	router.GET("/instagram/callback", ctrlInstagram.OauthInstagramCallback)
	router.GET("/get/instagram/profile", ctrlInstagram.FetchInstagramProfile)
	router.POST("/instagram/post/media", ctrlInstagram.PostInstagramReel)

	router.GET("/twitter/login", ctrlTwitter.OAuthTwitterLogin)
	router.GET("/twitter/callback", ctrlTwitter.OAuthTwitterCallback)

	router.POST("/twitter/post/tweet", ctrlTwitter.PostTweet)
	router.POST("/twitter/publish/media", ctrlTwitter.PostTweetWithVideo)

	router.GET("/get/twitter/profile", ctrlTwitter.FetchTwitterProfile)
	router.POST("/remove/twitter/account", ctrlTwitter.RemoveTwitterAccount)
	router.GET("/show/twitter/accounts", ctrlTwitter.ShowConnectedTwitterAccounts)

	router.POST("/publish/media/both", ctrlCommon.PostMediaToBoth)

}
