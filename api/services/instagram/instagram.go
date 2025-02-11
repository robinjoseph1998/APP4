package instgram

import (
	"APP4/api/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type InstagramServices struct {
	Repo repository.RepoInterfaces
}

func NewInstagramServices(repo repository.RepoInterfaces) InstagramServiceInterfaces {
	return &InstagramServices{
		Repo: repo}
}

func (igs *InstagramServices) GetIGBusinessID(accessToken string) (string, error) {
	// Step 1: Get Facebook Page ID
	facebookIDapi := fmt.Sprintf("https://graph.facebook.com/v19.0/me/accounts?access_token=%s", accessToken)

	req, err := http.NewRequest("GET", facebookIDapi, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create facebook_id request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("facebook id request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var FBResponse struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &FBResponse); err != nil {
		return "", fmt.Errorf("failed to parse Facebook Page ID response: %v", err)
	}

	if len(FBResponse.Data) == 0 {
		return "", fmt.Errorf("instagram Business ID not found")
	}

	pageID := FBResponse.Data[0].ID

	// Step 2: Get Instagram Business Account ID
	instagramIDapi := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?fields=instagram_business_account&access_token=%s", pageID, accessToken)

	req, err = http.NewRequest("GET", instagramIDapi, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create instagram_id request: %v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("instagram business id request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var IGResponse struct {
		InstagramBusinessAccount struct {
			ID string `json:"id"`
		} `json:"instagram_business_account"`
	}

	if err := json.Unmarshal(body, &IGResponse); err != nil {
		return "", fmt.Errorf("failed to parse IG Business ID response: %v", err)
	}

	if IGResponse.InstagramBusinessAccount.ID == "" {
		return "", fmt.Errorf("instagram Business ID not found")
	}

	return IGResponse.InstagramBusinessAccount.ID, nil
}
func (igs *InstagramServices) UploadInstagramReel(businessID, videoURL, caption, accessToken string) (string, error) {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/media", businessID)

	form := url.Values{}
	form.Set("media_type", "REELS")
	form.Set("video_url", videoURL)
	form.Set("caption", caption)
	form.Set("access_token", accessToken)

	resp, err := http.PostForm(apiURL, form)
	if err != nil {
		return "", fmt.Errorf("failed to upload video: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("📌 Instagram Video Upload Response:", string(body))

	var uploadResp struct {
		ID string `json:"id"`
	}
	json.Unmarshal(body, &uploadResp)

	if uploadResp.ID == "" {
		return "", fmt.Errorf("failed to get media ID from response: %s", string(body))
	}
	return uploadResp.ID, nil
}

func (ig *InstagramServices) CheckVideoProcessingStatus(mediaID, accessToken string) error {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?fields=status_code&access_token=%s", mediaID, accessToken)

	for {
		resp, err := http.Get(apiURL)
		if err != nil {
			return fmt.Errorf("failed to check video processing status: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Println("📌 Processing Status Response:", string(body))

		var statusResp struct {
			StatusCode string `json:"status_code"`
		}
		json.Unmarshal(body, &statusResp)

		if statusResp.StatusCode == "FINISHED" {
			return nil
		} else if statusResp.StatusCode == "ERROR" {
			return fmt.Errorf("video processing failed")
		}

		time.Sleep(5 * time.Second)
	}
}

func (ig *InstagramServices) PublishInstagramVideo(businessID, mediaID, accessToken string) (string, error) {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/media_publish", businessID)

	form := url.Values{}
	form.Set("creation_id", mediaID)
	form.Set("access_token", accessToken)

	fmt.Println("Api URL: ", apiURL)
	resp, err := http.PostForm(apiURL, form)
	if err != nil {
		return "", fmt.Errorf("failed to publish video: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("📌 Instagram Video Publish Response:", string(body))

	var publishResp struct {
		ID string `json:"id"`
	}
	json.Unmarshal(body, &publishResp)

	if publishResp.ID == "" {
		return "", fmt.Errorf("failed to publish video: %s", string(body))
	}
	return publishResp.ID, nil
}
