# CloudGenie MCP Server

A Model Context Protocol (MCP) server implementation in Go for the CloudGenie Internal Developer Platform (IDP). This server provides AI assistants with direct access to CloudGenie's infrastructure management APIs.

## Overview

This MCP server acts as a bridge between AI assistants (like Claude) and the CloudGenie backend API, enabling natural language interactions for infrastructure provisioning, resource management, and blueprint operations.

## Features

- ✅ Full MCP protocol support (version 2024-11-05)
- ✅ JSON-RPC 2.0 over stdio transport
- ✅ Complete CloudGenie API integration
- ✅ Blueprint discovery and management
- ✅ Resource CRUD operations
- ✅ Health monitoring
- ✅ Infrastructure provisioning workflows

## Architecture

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────────┐
│             │  MCP    │                  │  HTTP   │                 │
│ AI Assistant│◄───────►│ MCP Server (Go)  │◄────────│ CloudGenie API  │
│  (Claude)   │ stdio   │                  │         │    Backend      │
└─────────────┘         └──────────────────┘         └─────────────────┘
```

## Installation

### Prerequisites

- Go 1.21 or higher
- CloudGenie backend API running (or access to a hosted instance)

### Build from source

```bash
# Clone the repository
git clone https://github.com/deepakvbansode/idp-cloudgenie-mcp-server.git
cd idp-cloudgenie-mcp-server

# Download dependencies
go mod download

# Build the server
go build -o cloudgenie-mcp-server .
```

## Configuration

### Environment Variables

- `CLOUDGENIE_BACKEND_URL`: URL of the CloudGenie backend API (default: `http://localhost:8080`)

Example:

```bash
export CLOUDGENIE_BACKEND_URL="https://api.cloudgenie.example.com"
./cloudgenie-mcp-server
```

## Usage

### Running the Server

The server runs in stdio mode, communicating through standard input/output:

```bash
./cloudgenie-mcp-server
```

### MCP Client Configuration

#### Claude Desktop

Add to your configuration file (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

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

#### Other MCP Clients

Follow your client's specific configuration format for stdio-based MCP servers.

## Available Tools

**For Claude Desktop** (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "cloudgenie": {
      "command": "/path/to/cloudgenie-mcp-server"
    }
  }
}
```

**For other MCP clients**, follow their specific configuration format.

### 1. **cloudgenie_health_check**

Checks the health status of the CloudGenie backend API.

**Parameters:** None

**Example:**

```json
{
  "name": "cloudgenie_health_check",
  "arguments": {}
}
```

### 2. **cloudgenie_get_blueprints**

Retrieves all available blueprints from CloudGenie.

**Parameters:** None

**Returns:** List of all blueprints with their details (ID, name, description, version, category, tags)

**Example:**

```json
{
  "name": "cloudgenie_get_blueprints",
  "arguments": {}
}
```

### 3. **cloudgenie_get_blueprint**

Retrieves detailed information about a specific blueprint.

**Parameters:**

- `id` (string, required): The unique identifier of the blueprint

**Example:**

```json
{
  "name": "cloudgenie_get_blueprint",
  "arguments": {
    "id": "blueprint-123"
  }
}
```

### 4. **cloudgenie_get_resources**

Retrieves all resources from CloudGenie.

**Parameters:** None

**Returns:** List of all resources with their details (ID, name, type, status, blueprint ID, timestamps)

**Example:**

```json
{
  "name": "cloudgenie_get_resources",
  "arguments": {}
}
```

### 5. **cloudgenie_get_resource**

Retrieves detailed information about a specific resource.

**Parameters:**

- `id` (string, required): The unique identifier of the resource

**Example:**

```json
{
  "name": "cloudgenie_get_resource",
  "arguments": {
    "id": "resource-456"
  }
}
```

### 6. **cloudgenie_create_resource**

Creates a new resource in CloudGenie.

**Parameters:**

- `name` (string, required): Name of the resource
- `type` (string, required): Type of the resource (e.g., 'vm', 'database', 'storage')
- `blueprint_id` (string, required): ID of the blueprint to use
- `properties` (object, optional): Additional properties for the resource

**Example:**

```json
{
  "name": "cloudgenie_create_resource",
  "arguments": {
    "name": "my-database",
    "type": "database",
    "blueprint_id": "blueprint-123",
    "properties": {
      "size": "large",
      "region": "us-east-1"
    }
  }
}
```

### 7. **cloudgenie_delete_resource**

Deletes a resource from CloudGenie.

**Parameters:**

- `id` (string, required): The unique identifier of the resource to delete

**Example:**

```json
{
  "name": "cloudgenie_delete_resource",
  "arguments": {
    "id": "resource-456"
  }
}
```

### 8. **cloudgenie_update_resource_status**

Updates the status of a resource.

**Parameters:**

- `id` (string, required): The unique identifier of the resource
- `status` (string, required): New status (e.g., 'running', 'stopped', 'error')

**Example:**

```json
{
  "name": "cloudgenie_update_resource_status",
  "arguments": {
    "id": "resource-456",
    "status": "running"
  }
}
```

## CloudGenie API Endpoints

The MCP server integrates with the following CloudGenie backend endpoints:

| Endpoint                    | Method | Description            |
| --------------------------- | ------ | ---------------------- |
| `/v1/healthcheck`           | GET    | Health check endpoint  |
| `/v1/blueprints`            | GET    | List all blueprints    |
| `/v1/blueprints/{id}`       | GET    | Get blueprint by ID    |
| `/v1/resources`             | GET    | List all resources     |
| `/v1/resources`             | POST   | Create a new resource  |
| `/v1/resources/{id}`        | GET    | Get resource by ID     |
| `/v1/resources/{id}`        | DELETE | Delete a resource      |
| `/v1/resources/{id}/status` | PATCH  | Update resource status |

## Available Resources

The server exposes the following resources:

- `cloudgenie://server/info` - Server information
- `cloudgenie://server/capabilities` - Server capabilities
- `cloudgenie://docs/api` - API documentation

The server provides the following prompt templates:

- **select_blueprint** - Help users select the appropriate blueprint for their infrastructure needs
- **configure_resource** - Guide users through resource configuration
- **troubleshoot_resource** - Help troubleshoot resource issues
- **design_infrastructure** - Help design infrastructure architecture using CloudGenie

## Use Cases

### Example: Creating Infrastructure with AI

```
User: "I need to create a PostgreSQL database for my production environment"

AI Assistant (using MCP tools):
1. Calls cloudgenie_get_blueprints() to find database blueprints
2. Identifies the PostgreSQL blueprint
3. Calls cloudgenie_create_resource() with appropriate parameters
4. Returns the created resource details to the user
```

### Example: Monitoring Resource Status

```
User: "Show me the status of all my resources"

AI Assistant (using MCP tools):
1. Calls cloudgenie_get_resources() to fetch all resources
2. Formats and presents the resource list with their statuses
3. Can drill down into specific resources using cloudgenie_get_resource()
```

## Development

### Project Structure

```text
idp-cloudgenie-mcp-server/
├── main.go              # Main entry point with tool registration
├── go.mod               # Go module dependencies
├── cloudgenie/
│   └── client.go        # CloudGenie API client
├── mcp/
│   ├── server.go        # MCP server implementation
│   └── types.go         # MCP protocol type definitions
├── test.sh              # Test script
└── README.md            # This file
```

### Building

```bash
go build -o cloudgenie-mcp-server .
```

### Testing

Use the included test script:

```bash
./test.sh
```

Or test manually with JSON-RPC messages:

```bash
# Start the server
./cloudgenie-mcp-server

# Send test requests
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | ./cloudgenie-mcp-server
```

## Testing

### Manual Testing with JSON-RPC

You can test the server manually by sending JSON-RPC messages via stdin:

```bash
# Start the server
./cloudgenie-mcp-server

# Send initialization request
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}

# List available tools
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}

# Call a tool
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"echo","arguments":{"text":"Hello!"}}}
```

## Protocol Reference

This server implements the Model Context Protocol (MCP) specification. Key methods supported:

- `initialize` - Initialize the connection
- `tools/list` - List available tools
- `tools/call` - Execute a tool
- `resources/list` - List available resources
- `resources/read` - Read a resource
- `prompts/list` - List available prompts
- `prompts/get` - Get a prompt template
- `ping` - Health check

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on GitHub.
