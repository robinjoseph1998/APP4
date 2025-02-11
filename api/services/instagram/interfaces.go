package instgram

type InstagramServiceInterfaces interface {
	GetIGBusinessID(accessToken string) (string, error)
	UploadInstagramReel(businessID, videoURL, caption, accessToken string) (string, error)
	CheckVideoProcessingStatus(mediaID, accessToken string) error
	PublishInstagramVideo(businessID, mediaID, accessToken string) (string, error)
}
