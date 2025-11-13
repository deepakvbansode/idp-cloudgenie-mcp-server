package client

import (
	"bytes"
	
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type CGClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new CloudGenie API client
func NewCGClient(baseURL string) *CGClient {
	log.Printf("[CloudGenie] Creating API client with base URL: %s", baseURL)
	return &CGClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// DoRequest makes an HTTP request to the CloudGenie API (exported)
func (c *CGClient) DoRequest(method, path string, body interface{}) ([]byte, error) {
	url := c.BaseURL + path
	log.Printf("[CloudGenie API] %s %s", method, url)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			log.Printf("[CloudGenie API] Failed to marshal request: %v", err)
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
		log.Printf("[CloudGenie API] Request body: %s", string(jsonBody))
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		log.Printf("[CloudGenie API] Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Printf("[CloudGenie API] Request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[CloudGenie API] Response status: %d", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[CloudGenie API] Failed to read response: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("[CloudGenie API] API error (status %d): %s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[CloudGenie API] Request successful, response size: %d bytes", len(respBody))
	return respBody, nil
}