package prompts

import (
	"context"

	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterPrompts registers MCP prompts using the SDK
func RegisterPrompts(server *mcp.Server, client *client.CGClient) {
	// Example: Blueprint info prompt
	prompt := &mcp.Prompt{
		Name: "blueprint_info",
		Arguments: []*mcp.PromptArgument{{
			Name:        "blueprint_id",
			Description: "ID of the blueprint",
			Required:    true,
		}},
		Description: "Get information about a CloudGenie blueprint",
	}
	server.AddPrompt(prompt, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		
		return &mcp.GetPromptResult{
			Description: "Creates repository",
			Messages: []*mcp.PromptMessage{{
				Role:    "system",
				Content: &mcp.TextContent{Text: "Blueprint: xGitRepo"  + "\nDescription: create repository " },
			}},
		}, nil
	})
}