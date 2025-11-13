package tools

import (
	"context"
	
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/client"
)

// RegisterTools registers tool-related handlers to the MCP server
func RegisterTools(server *mcp.Server, client *client.CGClient) {
	log.Println("[TOOLS] Registering tool handlers...")

	// Create the resource tool with its own HTTP client
	resourceTool := NewResourceTool(client)
	blueprintTool := NewBlueprintTool(client)

	//-------blueprint tool  registration start-------
	// Register the get_blueprints tool
	getBlueprintsHandler := func(ctx context.Context, req *mcp.CallToolRequest, in GetBlueprintsInput) (*mcp.CallToolResult, GetBlueprintsOutput, error) {
		out, err := blueprintTool.GetBlueprints(ctx)
		if err != nil {
			return nil, GetBlueprintsOutput{}, err
		}
		return nil, GetBlueprintsOutput{Blueprints: out}, nil
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_blueprints",
		Description: "Get information about a CloudGenie blueprints",
	}, getBlueprintsHandler)
	
	//-------blueprint tool  registration end-------

	//--- resource tool registration start ---

	// Register the create_resource tool
	createResourceHandler := func(ctx context.Context, req *mcp.CallToolRequest, in CreateResourceInput) (*mcp.CallToolResult, CreateResourceOutput, error) {
		out, err := resourceTool.CreateResource(ctx, in)
		if err != nil {
			return nil, CreateResourceOutput{}, err
		}
		return nil, *out, nil
	}
	
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_resource",
		Description: "Create a new CloudGenie resource from a blueprint",
	}, createResourceHandler)
	
	
	getResourceHandler := func(ctx context.Context, req *mcp.CallToolRequest, in GetResourcesInput) (*mcp.CallToolResult, GetResourcesOutput, error) {
		out, err := resourceTool.GetResources()
		if err != nil {
			return nil, GetResourcesOutput{}, err
		}
		return nil, GetResourcesOutput{Resources: out}, nil
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_resources",
		Description: "Get all CloudGenie resources",
	}, getResourceHandler)

	getResourceByNameHandler := func(ctx context.Context, req *mcp.CallToolRequest, in GetResourceByNameInput) (*mcp.CallToolResult, GetResourceByNameOutput, error) {
		out, err := resourceTool.GetResourceByName(in.Name)
		if err != nil {
			return nil, GetResourceByNameOutput{}, err
		}
		return nil, GetResourceByNameOutput{Resource: out}, nil
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_resource_by_name",
		Description: "Get a CloudGenie resource by name",
	}, getResourceByNameHandler)
	//--- resource tool registration end ---
	log.Println("[TOOLS] Tool handlers registered successfully")
}