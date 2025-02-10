package services

type TwitterServiceInterfaces interface {
	InitializeMediaUpload(filePath string) (string, error)
	AppendMediaUpload(mediaID, filePath string) error
	FinalizeMediaUpload(mediaID string) error
	CheckMediaProcessingStatus(mediaID string) (string, int, error)
	PostTweet(status, mediaID string) error
}
