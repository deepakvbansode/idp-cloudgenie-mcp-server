package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/client"
)

type BlueprintTool struct {
	cgClient *client.CGClient
}

// NewBlueprintTool creates a new BlueprintTool instance
func NewBlueprintTool(cgClient *client.CGClient) *BlueprintTool {
	return &BlueprintTool{
		cgClient: cgClient,
	}
}

func (c *BlueprintTool) GetBlueprints(ctx context.Context) ([]Blueprint, error) {
	respBody, err := c.cgClient.DoRequest("GET", "/v1/blueprints", nil)
	if err != nil {
		return nil, err
	}

	var blueprints []Blueprint
	if err := json.Unmarshal(respBody, &blueprints); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return blueprints, nil
}