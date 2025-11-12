package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie"
	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/mcp"
)

func main() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	
	log.Println("==========================================")
	log.Println("CloudGenie MCP Server Starting...")
	log.Println("==========================================")
	
	// Get CloudGenie backend URL from environment variable or use default
	backendURL := os.Getenv("CLOUDGENIE_BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:50051" // Default local backend
	}
	log.Printf("[MAIN] CloudGenie backend URL: %s", backendURL)

	// Create CloudGenie API client
	log.Println("[MAIN] Creating CloudGenie API client...")
	client := cloudgenie.NewClient(backendURL)

	// Create a new MCP server
	log.Println("[MAIN] Initializing MCP server...")
	server := mcp.NewServer("idp-cloudgenie-mcp-server", "1.0.0")

	// Register CloudGenie API tools
	log.Println("[MAIN] Registering CloudGenie API tools...")
	registerTools(server, client)

	// Register example resources
	log.Println("[MAIN] Registering resources...")
	registerResources(server)

	// Register example prompts
	log.Println("[MAIN] Registering prompts...")
	registerPrompts(server)

	log.Println("[MAIN] CloudGenie MCP Server initialized successfully")
	log.Printf("[MAIN] Backend URL: %s", backendURL)

	// Start HTTP server if MCP_HTTP_PORT is set, else use stdio
	if port := os.Getenv("MCP_HTTP_PORT"); port != "" {
		addr := ":" + port
		log.Printf("[MAIN] Starting in HTTP mode on %s", addr)
		log.Println("==========================================")
		if err := server.StartHTTP(addr); err != nil {
			log.Fatalf("[MAIN] HTTP server error: %v", err)
		}
	} else {
		log.Println("[MAIN] Starting in stdio mode")
		log.Println("==========================================")
		if err := server.Start(); err != nil {
			log.Fatalf("[MAIN] Server error: %v", err)
		}
	}
}

func registerTools(server *mcp.Server, client *cloudgenie.Client) {
	// Health Check tool
	healthCheckTool := mcp.Tool{
		Name:        "cloudgenie_health_check",
		Description: "Checks the health status of the CloudGenie backend API",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
	}

	server.RegisterTool(healthCheckTool, func(args map[string]interface{}) ([]mcp.Content, error) {
		health, err := client.HealthCheck()
		if err != nil {
			return nil, fmt.Errorf("health check failed: %w", err)
		}

		return []mcp.Content{{
			Type: "text",
			Text: fmt.Sprintf("Status: %s\nMessage: %s", health.Status, health.Message),
		}}, nil
	})

	// Get Blueprints tool
	getBlueprintsTool := mcp.Tool{
		Name:        "cloudgenie_get_blueprints",
		Description: "Retrieves all available blueprints from CloudGenie",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
	}

	server.RegisterTool(getBlueprintsTool, func(args map[string]interface{}) ([]mcp.Content, error) {
		blueprints, err := client.GetBlueprints()
		if err != nil {
			return nil, fmt.Errorf("failed to get blueprints: %w", err)
		}

		result := fmt.Sprintf("Found %d blueprints:\n\n", len(blueprints))
		for i, bp := range blueprints {
			result += fmt.Sprintf("%d. %s (ID: %s)\n", i+1, bp.Name, bp.ID)
			result += fmt.Sprintf("   Description: %s\n", bp.Description)
			result += fmt.Sprintf("   Version: %s | Category: %s\n", bp.Version, bp.Category)
			if len(bp.Tags) > 0 {
				result += fmt.Sprintf("   Tags: %v\n", bp.Tags)
			}
			result += "\n"
		}

		return []mcp.Content{{
			Type: "text",
			Text: result,
		}}, nil
	})

	// Get Blueprint by ID tool
	getBlueprintByIDTool := mcp.Tool{
		Name:        "cloudgenie_get_blueprint",
		Description: "Retrieves detailed information about a specific blueprint by ID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"id": {
					Type:        "string",
					Description: "The unique identifier of the blueprint",
				},
			},
			Required: []string{"id"},
		},
	}

	server.RegisterTool(getBlueprintByIDTool, func(args map[string]interface{}) ([]mcp.Content, error) {
		id, ok := args["id"].(string)
		if !ok {
			return nil, fmt.Errorf("id argument must be a string")
		}

		blueprint, err := client.GetBlueprintByID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get blueprint: %w", err)
		}

		result := fmt.Sprintf("Blueprint Details:\n\n")
		result += fmt.Sprintf("ID: %s\n", blueprint.ID)
		result += fmt.Sprintf("Name: %s\n", blueprint.Name)
		result += fmt.Sprintf("Description: %s\n", blueprint.Description)
		result += fmt.Sprintf("Version: %s\n", blueprint.Version)
		result += fmt.Sprintf("Category: %s\n", blueprint.Category)
		if len(blueprint.Tags) > 0 {
			result += fmt.Sprintf("Tags: %v\n", blueprint.Tags)
		}
		if len(blueprint.Metadata) > 0 {
			result += fmt.Sprintf("Metadata: %v\n", blueprint.Metadata)
		}

		return []mcp.Content{{
			Type: "text",
			Text: result,
		}}, nil
	})

	// Get Resources tool
	getResourcesTool := mcp.Tool{
		Name:        "cloudgenie_get_resources",
		Description: "Retrieves all resources from CloudGenie",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
	}

	server.RegisterTool(getResourcesTool, func(args map[string]interface{}) ([]mcp.Content, error) {
		resources, err := client.GetResources()
		if err != nil {
			return nil, fmt.Errorf("failed to get resources: %w", err)
		}

		result := fmt.Sprintf("Found %d resources:\n\n", len(resources))
		for i, res := range resources {
			result += fmt.Sprintf("%d. %s (ID: %s)\n", i+1, res.Name, res.ID)
			result += fmt.Sprintf("   Type: %s | Status: %s\n", res.Type, res.Status)
			result += fmt.Sprintf("   Blueprint ID: %s\n", res.BlueprintID)
			result += fmt.Sprintf("   Created: %s | Updated: %s\n", res.CreatedAt, res.UpdatedAt)
			result += "\n"
		}

		return []mcp.Content{{
			Type: "text",
			Text: result,
		}}, nil
	})

	// Get Resource by ID tool
	getResourceByIDTool := mcp.Tool{
		Name:        "cloudgenie_get_resource",
		Description: "Retrieves detailed information about a specific resource by ID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"id": {
					Type:        "string",
					Description: "The unique identifier of the resource",
				},
			},
			Required: []string{"id"},
		},
	}

	server.RegisterTool(getResourceByIDTool, func(args map[string]interface{}) ([]mcp.Content, error) {
		id, ok := args["id"].(string)
		if !ok {
			return nil, fmt.Errorf("id argument must be a string")
		}

		resource, err := client.GetResourceByID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource: %w", err)
		}

		result := fmt.Sprintf("Resource Details:\n\n")
		result += fmt.Sprintf("ID: %s\n", resource.ID)
		result += fmt.Sprintf("Name: %s\n", resource.Name)
		result += fmt.Sprintf("Type: %s\n", resource.Type)
		result += fmt.Sprintf("Status: %s\n", resource.Status)
		result += fmt.Sprintf("Blueprint ID: %s\n", resource.BlueprintID)
		result += fmt.Sprintf("Created At: %s\n", resource.CreatedAt)
		result += fmt.Sprintf("Updated At: %s\n", resource.UpdatedAt)
		if len(resource.Properties) > 0 {
			result += fmt.Sprintf("Properties: %v\n", resource.Properties)
		}

		return []mcp.Content{{
			Type: "text",
			Text: result,
		}}, nil
	})

	// Create Resource tool
	createResourceTool := mcp.Tool{
		Name:        "cloudgenie_create_resource",
		Description: "Creates a new resource in CloudGenie",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"name": {
					Type:        "string",
					Description: "Name of the resource",
				},
				"type": {
					Type:        "string",
					Description: "Type of the resource (e.g., 'vm', 'database', 'storage')",
				},
				"blueprint_id": {
					Type:        "string",
					Description: "ID of the blueprint to use",
				},
				"properties": {
					Type:        "object",
					Description: "Additional properties for the resource (JSON object)",
				},
			},
			Required: []string{"name", "type", "blueprint_id"},
		},
	}

	server.RegisterTool(createResourceTool, func(args map[string]interface{}) ([]mcp.Content, error) {
		name, nameOk := args["name"].(string)
		resourceType, typeOk := args["type"].(string)
		blueprintID, bpOk := args["blueprint_id"].(string)

		if !nameOk || !typeOk || !bpOk {
			return nil, fmt.Errorf("name, type, and blueprint_id are required")
		}

		req := cloudgenie.CreateResourceRequest{
			Name:        name,
			Type:        resourceType,
			BlueprintID: blueprintID,
			Properties:  make(map[string]interface{}),
		}

		if props, ok := args["properties"].(map[string]interface{}); ok {
			req.Properties = props
		}

		resource, err := client.CreateResource(req)
		if err != nil {
			return nil, fmt.Errorf("failed to create resource: %w", err)
		}

		result := fmt.Sprintf("Resource created successfully!\n\n")
		result += fmt.Sprintf("ID: %s\n", resource.ID)
		result += fmt.Sprintf("Name: %s\n", resource.Name)
		result += fmt.Sprintf("Type: %s\n", resource.Type)
		result += fmt.Sprintf("Status: %s\n", resource.Status)
		result += fmt.Sprintf("Blueprint ID: %s\n", resource.BlueprintID)

		return []mcp.Content{{
			Type: "text",
			Text: result,
		}}, nil
	})

	// Delete Resource tool
	deleteResourceTool := mcp.Tool{
		Name:        "cloudgenie_delete_resource",
		Description: "Deletes a resource from CloudGenie",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"id": {
					Type:        "string",
					Description: "The unique identifier of the resource to delete",
				},
			},
			Required: []string{"id"},
		},
	}

	server.RegisterTool(deleteResourceTool, func(args map[string]interface{}) ([]mcp.Content, error) {
		id, ok := args["id"].(string)
		if !ok {
			return nil, fmt.Errorf("id argument must be a string")
		}

		err := client.DeleteResource(id)
		if err != nil {
			return nil, fmt.Errorf("failed to delete resource: %w", err)
		}

		return []mcp.Content{{
			Type: "text",
			Text: fmt.Sprintf("Resource %s deleted successfully", id),
		}}, nil
	})

	// Update Resource Status tool
	updateResourceStatusTool := mcp.Tool{
		Name:        "cloudgenie_update_resource_status",
		Description: "Updates the status of a resource in CloudGenie",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"id": {
					Type:        "string",
					Description: "The unique identifier of the resource",
				},
				"status": {
					Type:        "string",
					Description: "New status for the resource (e.g., 'running', 'stopped', 'error')",
				},
			},
			Required: []string{"id", "status"},
		},
	}

	server.RegisterTool(updateResourceStatusTool, func(args map[string]interface{}) ([]mcp.Content, error) {
		id, idOk := args["id"].(string)
		status, statusOk := args["status"].(string)

		if !idOk || !statusOk {
			return nil, fmt.Errorf("id and status are required")
		}

		resource, err := client.UpdateResourceStatus(id, status)
		if err != nil {
			return nil, fmt.Errorf("failed to update resource status: %w", err)
		}

		result := fmt.Sprintf("Resource status updated successfully!\n\n")
		result += fmt.Sprintf("ID: %s\n", resource.ID)
		result += fmt.Sprintf("Name: %s\n", resource.Name)
		result += fmt.Sprintf("New Status: %s\n", resource.Status)
		result += fmt.Sprintf("Updated At: %s\n", resource.UpdatedAt)

		return []mcp.Content{{
			Type: "text",
			Text: result,
		}}, nil
	})
}

func registerResources(server *mcp.Server) {
	// Example resource - server info
	server.RegisterResource(mcp.Resource{
		URI:         "cloudgenie://server/info",
		Name:        "Server Information",
		Description: "Information about the CloudGenie MCP server",
		MimeType:    "text/plain",
	})

	// Example resource - capabilities
	server.RegisterResource(mcp.Resource{
		URI:         "cloudgenie://server/capabilities",
		Name:        "Server Capabilities",
		Description: "List of capabilities supported by this server",
		MimeType:    "application/json",
	})

	// Example resource - documentation
	server.RegisterResource(mcp.Resource{
		URI:         "cloudgenie://docs/api",
		Name:        "API Documentation",
		Description: "API documentation for CloudGenie MCP server",
		MimeType:    "text/markdown",
	})
}

func registerPrompts(server *mcp.Server) {
	log.Println("[MAIN] Registering prompts...")

	// 1. CAPABILITIES PROMPT - "How can you help me?"
	server.RegisterPrompt(mcp.Prompt{
		Name:        "cloudgenie_capabilities",
		Description: `Explains all CloudGenie IDP (Internal Developer Platform) MCP server capabilities. Shows available blueprints (currently xgitrepo for Git repositories, with React app + K8s deployment coming soon), resource management, and workflows. Use when user asks "what can you do?" or "how can you help me?"`,
		Arguments:   []mcp.Argument{},
	})

	// 2. DEPLOY APPLICATION PROMPT - "Deploy react application in k8s"
	server.RegisterPrompt(mcp.Prompt{
		Name:        "deploy_application_k8s",
		Description: `Guide for deploying applications (React, Node.js, etc.) to Kubernetes using IDP blueprints. Currently explains that React app + K8s deployment blueprint is coming soon. When available, will use blueprints to deploy containerized apps. Use when user wants to deploy applications to Kubernetes.`,
		Arguments: []mcp.Argument{
			{
				Name:        "app_type",
				Description: "Type of application to deploy (e.g., react, nodejs, python, java)",
				Required:    false,
			},
			{
				Name:        "app_name",
				Description: "Name for the application deployment",
				Required:    false,
			},
		},
	})

	// 3. GIT REPOSITORY PROMPT - "Create a git repository"
	// 3. GIT REPOSITORY PROMPT - "Create a git repository"
	server.RegisterPrompt(mcp.Prompt{
		Name:        "create_git_repository",
		Description: `Complete guide for creating Git repositories using the xgitrepo blueprint in the IDP (Internal Developer Platform). Explains how to discover the xgitrepo blueprint, understand its parameters, and create Git repository resources. Use when user wants to create a Git repository.`,
		Arguments: []mcp.Argument{
			{
				Name:        "repo_name",
				Description: "Desired name for the Git repository",
				Required:    false,
			},
			{
				Name:        "description",
				Description: "Description of the repository purpose",
				Required:    false,
			},
		},
	})

	// 4. DEPLOY INFRASTRUCTURE PROMPT - General resource creation via blueprints
	server.RegisterPrompt(mcp.Prompt{
		Name:        "deploy_infrastructure",
		Description: `Step-by-step guide for creating resources using IDP blueprints. Explains the blueprint-based model: first discover available blueprints (like xgitrepo), understand their parameters, then create resource instances. Use when user wants to create any resource type.`,
		Arguments: []mcp.Argument{
			{
				Name:        "blueprint_id",
				Description: "ID of the blueprint to use (e.g., xgitrepo)",
				Required:    false,
			},
			{
				Name:        "resource_name",
				Description: "Name for the resource to be created",
				Required:    false,
			},
		},
	})

	// 5. TROUBLESHOOT RESOURCE PROMPT
	server.RegisterPrompt(mcp.Prompt{
		Name:        "troubleshoot_resource",
		Description: `Comprehensive troubleshooting guide for IDP resources. Covers checking status, diagnosing common issues (pending, failed, error states), backend health checks, and recovery actions. Use when user reports problems with resource creation or deployments.`,
		Arguments: []mcp.Argument{
			{
				Name:        "resource_id",
				Description: "ID of the resource experiencing issues",
				Required:    false,
			},
		},
	})

	// 6. LIST BLUEPRINTS PROMPT
	server.RegisterPrompt(mcp.Prompt{
		Name:        "list_blueprints",
		Description: `Guide for discovering available IDP blueprints. Explains that blueprints define what resources can be created (currently xgitrepo for Git repositories, with React+K8s blueprint coming). Shows how to list blueprints, examine parameters, and choose the right one. Use when user wants to see what can be created.`,
		Arguments:   []mcp.Argument{},
	})

	log.Println("[MAIN] Prompts registered successfully (6 prompts)")
}
