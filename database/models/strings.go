package models

type ProfileResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	AccountType string `json:"account_type"`
	MediaCount  int    `json:"media_count"`
}
