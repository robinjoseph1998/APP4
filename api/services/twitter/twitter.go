package services

import (
	"APP4/api/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
)

type TwitterServices struct {
	Repo repository.RepoInterfaces
}

func NewTwitterServices(repo repository.RepoInterfaces) TwitterServiceInterfaces {
	return &TwitterServices{
		Repo: repo}
}

func getOAuth1Client() *http.Client {
	consumerKey := os.Getenv("TWITTER_API_KEY")
	consumerSecret := os.Getenv("TWITTER_API_SECRET_KEY")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	return config.Client(oauth1.NoContext, token)
}

// INIT - Initialize media upload
func (xservice *TwitterServices) InitializeMediaUpload(filePath string) (string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}
	totalBytes := fileInfo.Size()

	// Twitter API URL
	urlStr := "https://upload.twitter.com/1.1/media/upload.json"

	// Form data
	data := url.Values{}
	data.Set("command", "INIT")
	data.Set("media_type", "video/mp4")
	data.Set("media_category", "tweet_video") // Required for tweet videos
	data.Set("total_bytes", strconv.FormatInt(totalBytes, 10))

	// Create HTTP request
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Use OAuth1 client for authentication
	client := getOAuth1Client()

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("INIT Response Status:", resp.Status)
	fmt.Println("INIT Response Body:", string(body))

	// Check if request failed
	// if resp.StatusCode != http.StatusOK {
	// 	return "", fmt.Errorf("twitter INIT API error: %s - %s", resp.Status, string(body))
	// }

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}
	mediaID, ok := result["media_id_string"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format, missing media_id")
	}
	return mediaID, nil
}

// 2️⃣ APPEND - Upload video in chunks
func (xservice *TwitterServices) AppendMediaUpload(mediaID, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	const chunkSize = 5 * 1024 * 1024 // 5MB per chunk
	buffer := make([]byte, chunkSize)

	client := getOAuth1Client()

	for segmentIndex := 0; ; segmentIndex++ {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file: %v", err)
		}
		if bytesRead == 0 {
			break // Done reading the file
		}

		// Twitter API URL for APPEND
		urlStr := "https://upload.twitter.com/1.1/media/upload.json"

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("command", "APPEND")
		_ = writer.WriteField("media_id", mediaID)
		_ = writer.WriteField("segment_index", strconv.Itoa(segmentIndex))
		part, _ := writer.CreateFormFile("media", "video.mp4")
		part.Write(buffer[:bytesRead])
		writer.Close()

		req, _ := http.NewRequest("POST", urlStr, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send APPEND request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("APPEND error: %s - %s", resp.Status, string(body))
		}
	}
	return nil
}

// FINALIZE - Complete the upload
func (xservice *TwitterServices) FinalizeMediaUpload(mediaID string) error {
	urlStr := "https://upload.twitter.com/1.1/media/upload.json"
	data := url.Values{}
	data.Set("command", "FINALIZE")
	data.Set("media_id", mediaID)

	client := getOAuth1Client()

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to finalize upload: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse FINALIZE response: %v", err)
	}
	if processingInfo, ok := result["processing_info"].(map[string]interface{}); ok {
		state := processingInfo["state"].(string)
		checkAfterSecs := int(processingInfo["check_after_secs"].(float64))

		for state != "succeeded" {
			fmt.Println("Waiting for video processing... Checking again in", checkAfterSecs, "seconds.")
			time.Sleep(time.Duration(checkAfterSecs) * time.Second)

			state, checkAfterSecs, err = xservice.CheckMediaProcessingStatus(mediaID)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("✅ Video processing complete!")
	return nil
}

// Helper function to check video processing status
func (xservice *TwitterServices) CheckMediaProcessingStatus(mediaID string) (string, int, error) {
	urlStr := "https://upload.twitter.com/1.1/media/upload.json?command=STATUS&media_id=" + mediaID
	client := getOAuth1Client()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create STATUS request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send STATUS request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, fmt.Errorf("failed to parse STATUS response: %v", err)
	}

	if processingInfo, ok := result["processing_info"].(map[string]interface{}); ok {
		state := processingInfo["state"].(string)
		checkAfterSecs := 5 // Default retry time

		if val, ok := processingInfo["check_after_secs"].(float64); ok {
			checkAfterSecs = int(val)
		}

		return state, checkAfterSecs, nil
	}

	return "succeeded", 0, nil
}

// POST TWEET - Attach media_id to tweet
func (xservice *TwitterServices) PostTweet(status, mediaID string) error {
	urlStr := "https://api.twitter.com/2/tweets"

	client := getOAuth1Client()
	payload := map[string]interface{}{
		"text": status,
		"media": map[string]interface{}{
			"media_ids": []string{mediaID},
		},
	}
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send tweet: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("twitter Tweet API error: %s - %s", resp.Status, string(body))
	}
	return nil
}
