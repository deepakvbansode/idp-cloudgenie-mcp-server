# CloudGenie MCP Server - Implementation Summary

## Overview

This MCP (Model Context Protocol) server provides AI assistants like Claude with direct access to the CloudGenie IDP (Internal Developer Platform) backend API. It enables natural language interactions for infrastructure provisioning, resource management, and blueprint operations.

## Project Structure

```
idp-cloudgenie-mcp-server/
├── main.go                      # Main entry point with tool registration
├── go.mod                       # Go module dependencies
├── cloudgenie/
│   └── client.go               # CloudGenie API HTTP client
├── mcp/
│   ├── server.go               # MCP server implementation (stdio transport)
│   └── types.go                # MCP protocol type definitions
├── test.sh                     # Test script for manual testing
├── claude_desktop_config.json  # Example configuration for Claude Desktop
├── .gitignore                  # Git ignore patterns
└── README.md                   # Comprehensive documentation

```

## Key Components

### 1. MCP Server (`mcp/server.go` & `mcp/types.go`)

- Full MCP protocol implementation (version 2024-11-05)
- JSON-RPC 2.0 message handling
- Stdio-based communication
- Tool, Resource, and Prompt management
- Request routing and error handling

### 2. CloudGenie API Client (`cloudgenie/client.go`)

- HTTP client for CloudGenie backend API
- Complete CRUD operations for resources
- Blueprint discovery and retrieval
- Health check monitoring
- Type-safe request/response models

### 3. Main Application (`main.go`)

- Server initialization
- Tool registration with CloudGenie client integration
- Resource and prompt registration
- Environment-based configuration

## Implemented Tools

The server provides 8 CloudGenie-specific tools:

1. **cloudgenie_health_check** - Check backend API health
2. **cloudgenie_get_blueprints** - List all available blueprints
3. **cloudgenie_get_blueprint** - Get specific blueprint details
4. **cloudgenie_get_resources** - List all resources
5. **cloudgenie_get_resource** - Get specific resource details
6. **cloudgenie_create_resource** - Create a new resource
7. **cloudgenie_delete_resource** - Delete a resource
8. **cloudgenie_update_resource_status** - Update resource status

## CloudGenie API Integration

The server integrates with these CloudGenie backend endpoints:

| Endpoint                    | Method | MCP Tool                          |
| --------------------------- | ------ | --------------------------------- |
| `/v1/healthcheck`           | GET    | cloudgenie_health_check           |
| `/v1/blueprints`            | GET    | cloudgenie_get_blueprints         |
| `/v1/blueprints/{id}`       | GET    | cloudgenie_get_blueprint          |
| `/v1/resources`             | GET    | cloudgenie_get_resources          |
| `/v1/resources`             | POST   | cloudgenie_create_resource        |
| `/v1/resources/{id}`        | GET    | cloudgenie_get_resource           |
| `/v1/resources/{id}`        | DELETE | cloudgenie_delete_resource        |
| `/v1/resources/{id}/status` | PATCH  | cloudgenie_update_resource_status |

## Configuration

### Environment Variables

- `CLOUDGENIE_BACKEND_URL`: CloudGenie backend API URL (default: `http://localhost:8080`)

### Claude Desktop Setup

1. Build the server:

   ```bash
   go build -o cloudgenie-mcp-server .
   ```

2. Add to Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

   ```json
   {
     "mcpServers": {
       "cloudgenie": {
         "command": "/path/to/cloudgenie-mcp-server",
         "env": {
           "CLOUDGENIE_BACKEND_URL": "http://localhost:8080"
         }
       }
     }
   }
   ```

3. Restart Claude Desktop

## Usage Examples

### Example 1: List Available Blueprints

```
User: "What blueprints are available in CloudGenie?"

Claude uses: cloudgenie_get_blueprints()
Returns: List of all blueprints with names, descriptions, versions, etc.
```

### Example 2: Create a Resource

```
User: "Create a PostgreSQL database named 'prod-db' using blueprint 'postgres-ha'"

Claude uses:
1. cloudgenie_get_blueprints() to find the blueprint ID
2. cloudgenie_create_resource() with:
   - name: "prod-db"
   - type: "database"
   - blueprint_id: "postgres-ha"
```

### Example 3: Monitor Resources

```
User: "Show me all my running resources"

Claude uses:
1. cloudgenie_get_resources() to fetch all resources
2. Filters and displays resources with status "running"
```

## Testing

### Manual Testing with Test Script

```bash
./test.sh
```

This sends a series of JSON-RPC messages to test:

- Server initialization
- Tool listing
- Health check
- Blueprint retrieval
- Resource operations
- Ping/pong

### Testing with CloudGenie Backend

1. Start CloudGenie backend on `http://localhost:8080`
2. Run the MCP server:
   ```bash
   export CLOUDGENIE_BACKEND_URL="http://localhost:8080"
   ./cloudgenie-mcp-server
   ```
3. Use Claude Desktop or another MCP client to interact

## Architecture Flow

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────────┐
│             │  stdin/ │                  │  HTTP   │                 │
│ AI Assistant│  stdout │ MCP Server (Go)  │  REST   │ CloudGenie API  │
│  (Claude)   │◄───────►│                  │◄────────│    Backend      │
│             │ JSON-RPC│  - Tool handlers │         │                 │
│             │         │  - API client    │         │  - Blueprints   │
│             │         │  - Protocol impl │         │  - Resources    │
└─────────────┘         └──────────────────┘         └─────────────────┘
```

## Development Notes

### Adding New Tools

To add a new CloudGenie API tool:

1. Add the API method to `cloudgenie/client.go`:

   ```go
   func (c *Client) NewOperation() (*Response, error) {
       // Implementation
   }
   ```

2. Register the tool in `main.go`:
   ```go
   newTool := mcp.Tool{
       Name: "cloudgenie_new_operation",
       Description: "Description",
       InputSchema: mcp.InputSchema{...},
   }
   server.RegisterTool(newTool, func(args map[string]interface{}) ([]mcp.Content, error) {
       // Call client.NewOperation()
       // Format and return results
   })
   ```

### Error Handling

- API errors are propagated to the AI assistant
- HTTP errors include status codes and response bodies
- Tool execution errors return structured error messages

## Features

✅ **Implemented:**

- Full MCP protocol support (2024-11-05)
- Complete CloudGenie API integration
- All CRUD operations for resources
- Blueprint discovery
- Health monitoring
- Environment-based configuration
- Comprehensive error handling
- Test script
- Documentation

🔮 **Future Enhancements:**

- Authentication/authorization support
- Caching for frequently accessed data
- Webhooks for real-time updates
- Batch operations
- Resource filtering and search
- Metrics and logging
- Configuration file support

## Dependencies

- Go 1.21+
- Standard library packages (net/http, encoding/json, etc.)
- No external dependencies (pure Go implementation)

## License

MIT License

## Support

For issues and questions:

- GitHub Issues: https://github.com/deepakvbansode/idp-cloudgenie-mcp-server/issues
- Email: (your email)

---

**Built with ❤️ for the CloudGenie IDP Platform**
