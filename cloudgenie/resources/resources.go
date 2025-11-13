package resources

import (
	"context"

	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterResources registers MCP resources using the SDK
func RegisterResources(server *mcp.Server, client *client.CGClient) {
	// Example: Blueprint resource
	server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: "cloudgenie://blueprint/{id}",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:  req.Params.URI,
				Text: "Blueprint: xGitRepo" + "\nDescription: create repository ",
			}},
		}, nil
	})
}