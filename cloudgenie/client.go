package cloudgenie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Client represents a CloudGenie API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new CloudGenie API client
func NewClient(baseURL string) *Client {
	log.Printf("[CloudGenie] Creating API client with base URL: %s", baseURL)
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Blueprint represents a CloudGenie blueprint
type Blueprint struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Version     string                 `json:"version"`
	Category    string                 `json:"category"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Resource represents a CloudGenie resource
type Resource struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	BlueprintID string                 `json:"blueprint_id"`
	Status      string                 `json:"status"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Properties  map[string]interface{} `json:"properties"`
}

// CreateResourceRequest represents a request to create a resource
type CreateResourceRequest struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	BlueprintID string                 `json:"blueprint_id"`
	Properties  map[string]interface{} `json:"properties"`
}

// UpdateStatusRequest represents a request to update resource status
type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// doRequest makes an HTTP request to the CloudGenie API
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
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

// HealthCheck checks the health of the CloudGenie API
func (c *Client) HealthCheck() (*HealthCheckResponse, error) {
	respBody, err := c.doRequest("GET", "/v1/healthcheck", nil)
	if err != nil {
		return nil, err
	}

	var response HealthCheckResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetBlueprints retrieves all blueprints
func (c *Client) GetBlueprints() ([]Blueprint, error) {
	respBody, err := c.doRequest("GET", "/v1/blueprints", nil)
	if err != nil {
		return nil, err
	}

	var blueprints []Blueprint
	if err := json.Unmarshal(respBody, &blueprints); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return blueprints, nil
}

// GetBlueprintByID retrieves a specific blueprint by ID
func (c *Client) GetBlueprintByID(id string) (*Blueprint, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/v1/blueprints/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var blueprint Blueprint
	if err := json.Unmarshal(respBody, &blueprint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &blueprint, nil
}

// GetResources retrieves all resources
func (c *Client) GetResources() ([]Resource, error) {
	respBody, err := c.doRequest("GET", "/v1/resources", nil)
	if err != nil {
		return nil, err
	}

	var resources []Resource
	if err := json.Unmarshal(respBody, &resources); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resources, nil
}

// GetResourceByID retrieves a specific resource by ID
func (c *Client) GetResourceByID(id string) (*Resource, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/v1/resources/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var resource Resource
	if err := json.Unmarshal(respBody, &resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resource, nil
}

// CreateResource creates a new resource
func (c *Client) CreateResource(req CreateResourceRequest) (*Resource, error) {
	respBody, err := c.doRequest("POST", "/v1/resources", req)
	if err != nil {
		return nil, err
	}

	var resource Resource
	if err := json.Unmarshal(respBody, &resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resource, nil
}

// DeleteResource deletes a resource by ID
func (c *Client) DeleteResource(id string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/v1/resources/%s", id), nil)
	return err
}

// UpdateResourceStatus updates the status of a resource
func (c *Client) UpdateResourceStatus(id string, status string) (*Resource, error) {
	req := UpdateStatusRequest{Status: status}
	respBody, err := c.doRequest("PATCH", fmt.Sprintf("/v1/resources/%s/status", id), req)
	if err != nil {
		return nil, err
	}

	var resource Resource
	if err := json.Unmarshal(respBody, &resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resource, nil
}
