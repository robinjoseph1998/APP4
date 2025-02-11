package services

type TwitterServiceInterfaces interface {
	InitializeMediaUpload(filePath string) (string, error)
	AppendMediaUpload(mediaID, filePath string) error
	FinalizeMediaUpload(mediaID string) error
	CheckMediaProcessingStatus(mediaID string) (string, int, error)
	PostTweet(caption, mediaID string) error
	PublicUrlVedioDownloader(videoURL string) (string, error)
}
