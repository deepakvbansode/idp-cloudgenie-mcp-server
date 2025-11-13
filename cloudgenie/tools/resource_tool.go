package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/client"
)

// ResourceTool handles resource-related operations
type ResourceTool struct {
	cgClient *client.CGClient
}

// NewResourceTool creates a new ResourceTool instance
func NewResourceTool(cgClient *client.CGClient) *ResourceTool {
	return &ResourceTool{
		cgClient: cgClient,
	}
}


// CreateResource creates a new resource via the CloudGenie API
func (rt *ResourceTool) CreateResource(ctx context.Context, in CreateResourceInput) (*CreateResourceOutput, error) {
	log.Printf("[ResourceTool] Creating resource: name=%s, blueprint=%s", in.Name, in.BlueprintName)

	respBody, err := rt.cgClient.DoRequest("POST", "/v1/resources", in)
	if err != nil {
		return nil, err
	}

	var resource Resource
	if err := json.Unmarshal(respBody, &resource); err != nil {
		log.Printf("[ResourceTool] Failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Printf("[ResourceTool] Resource created successfully:  name=%s, status=%s", 
		resource.Name, resource.Status)

	return &CreateResourceOutput{
		Name:   resource.Name,
		Status: resource.Status,
	}, nil
}

func (rt *ResourceTool) GetResources() ([]Resource, error) {
	respBody, err := rt.cgClient.DoRequest("GET", "/v1/resources", nil)
	if err != nil {
		return nil, err
	}

	var resources []Resource
	if err := json.Unmarshal(respBody, &resources); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resources, nil
}

func (rt *ResourceTool) GetResourceByName(name string) (*Resource, error) {
	respBody, err := rt.cgClient.DoRequest("GET", fmt.Sprintf("/v1/resources/%s", name), nil)
	if err != nil {
		return nil, err
	}

	var resource Resource
	if err := json.Unmarshal(respBody, &resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resource, nil
}