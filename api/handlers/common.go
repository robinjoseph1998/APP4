package api

import (
	"APP4/api/repository"
	igservices "APP4/api/services/instagram"
	xservices "APP4/api/services/twitter"
	"APP4/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonAuthHandlers struct {
	Repo       repository.RepoInterfaces
	XServices  xservices.TwitterServiceInterfaces
	IgServices igservices.InstagramServiceInterfaces
}

func NewCommonAuthHandlers(repo repository.RepoInterfaces, xservice xservices.TwitterServiceInterfaces, igservices igservices.InstagramServiceInterfaces) *CommonAuthHandlers {
	return &CommonAuthHandlers{
		Repo:       repo,
		XServices:  xservice,
		IgServices: igservices}
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	errorResponse := gin.H{
		"error":       true,
		"status_code": statusCode,
		"message":     message,
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

func (ctrl *CommonAuthHandlers) PostMediaToBoth(c *gin.Context) {

	igUserName := c.PostForm("ig_username")
	twtrUserName := c.PostForm("twtr_username")
	videoURL := c.PostForm("video_url")
	caption := c.PostForm("caption")

	if igUserName == "" || twtrUserName == "" || videoURL == "" || caption == "" {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request: request field empty", nil)
		return
	}

	// twitterAcccesToken, err := ctrl.Repo.FetchTwitterAccessTokenFromDB(twtrUserName)
	// if err != nil {
	// 	ErrorResponse(c, http.StatusInternalServerError, "can't fetch twitter access token", err)
	// 	return
	// }
	instagramAccessToken, err := ctrl.Repo.FetchInstagramAccessTokenFromDB(igUserName)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "can't fetch instagram access token", err)
		return
	}

	igbusinessID, err := ctrl.IgServices.GetIGBusinessID(instagramAccessToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get Instagram Business ID", err)
		return
	}
	filePath, err := ctrl.XServices.PublicUrlVedioDownloader(videoURL)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "can't download the vedio", err)
		return
	}
	mediaID, err := ctrl.XServices.InitializeMediaUpload(filePath)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize media upload", err)
		return
	}

	err = ctrl.XServices.AppendMediaUpload(mediaID, filePath)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to upload media chunks", err)
		return
	}

	err = ctrl.XServices.FinalizeMediaUpload(mediaID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to finalize media upload", err)
		return
	}

	err = ctrl.XServices.PostTweet(caption, mediaID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to post tweet", err)
		return
	}
	mediaID, err = ctrl.IgServices.UploadInstagramReel(igbusinessID, videoURL, caption, instagramAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video", "details": err.Error()})
		return
	}

	err = ctrl.IgServices.CheckVideoProcessingStatus(mediaID, instagramAccessToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Video processing failed", err)
		return
	}

	postID, err := ctrl.IgServices.PublishInstagramVideo(igbusinessID, mediaID, instagramAccessToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to publish media", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Media posted to Twitter and Instagram successfully!",
		"post_id": postID,
	})

}
