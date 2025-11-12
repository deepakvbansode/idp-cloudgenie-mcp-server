# CloudGenie IDP MCP Server

A Model Context Protocol (MCP) server that enables AI assistants (Claude, Gemini, etc.) to interact with the **CloudGenie Internal Developer Platform (IDP)**. The IDP follows a blueprint-based model where blueprints define resource types that can be created.

## 🎯 What is the IDP?

CloudGenie IDP is an **Internal Developer Platform** that enables developers to create resources through standardized blueprints:

- **Blueprints**: Template definitions that specify what types of resources can be created
- **Resources**: Actual instances created from blueprints

Think of blueprints as "classes" and resources as "objects" in OOP terms.

### Current Capabilities

✅ **xgitrepo Blueprint**: Create Git Repositories

- Create Git repositories through the IDP
- Configure repository settings (visibility, default branch, etc.)
- Manage repository lifecycle

🚧 **Coming Soon**: React App + Kubernetes Deployment Blueprint

- Deploy React applications to Kubernetes clusters
- Automated build and deployment pipeline
- Full application lifecycle management

## Architecture

```
┌─────────────────┐
│  AI Model       │  User: "Create a git repo called my-app"
│ (Claude/Gemini) │
└────────┬────────┘
         │ MCP Protocol (JSON-RPC)
         ↓
┌────────────────────────┐
│ CloudGenie MCP Server  │  ← This Project
│  - Protocol Handler    │
│  - Tool Registry       │
│  - Prompt Templates    │
└────────┬───────────────┘
         │ HTTP/REST
         ↓
┌───────────────────────┐
│ CloudGenie IDP Backend│
│  - Blueprint Registry │
│  - Resource Manager   │
│  - xgitrepo handler   │
└────────┬──────────────┘
         │
         ↓
┌────────────────────┐
│ Infrastructure     │
│  - Git Repositories│
│  - K8s Apps (soon) │
└────────────────────┘
```

## IDP Workflow

### 1. Discovery Phase

```bash
User: "What can I create?"
Model → get_blueprints tool
Result: [xgitrepo, ...]  # List of available blueprints
```

### 2. Inspection Phase

```bash
User: "How do I create a Git repository?"
Model → get_blueprint("xgitrepo") tool
Result: {
  id: "xgitrepo",
  name: "Git Repository",
  parameters: {
    name: {required: true, type: "string"},
    description: {required: false, type: "string"},
    visibility: {required: false, enum: ["public", "private"]}
  }
}
```

### 3. Creation Phase

```bash
User: "Create a git repository called my-awesome-app"
Model → create_resource tool
Request: {
  blueprint_id: "xgitrepo",
  name: "my-awesome-app",
  type: "git-repository",
  properties: {
    description: "My awesome application",
    visibility: "private"
  }
}
Result: {
  id: "res-abc-123",
  status: "pending"
}
```

### 4. Monitoring Phase

```bash
Model → get_resource("res-abc-123") tool
Result: {
  id: "res-abc-123",
  name: "my-awesome-app",
  status: "ready",  # pending → creating → ready
  url: "https://git.example.com/my-awesome-app"
}
```

## Features

### MCP Protocol Support

- ✅ Full MCP protocol compliance (version 2024-11-05)
- ✅ JSON-RPC 2.0 over stdio (for Claude Desktop)
- ✅ JSON-RPC 2.0 over HTTP (for Kubernetes deployment)
- ✅ Dual-mode operation

### Blueprint Management

- ✅ List all available blueprints
- ✅ Get detailed blueprint specifications
- ✅ Understand required/optional parameters
- ✅ Currently supports: xgitrepo blueprint

### Resource Lifecycle

- ✅ Create resources from blueprints
- ✅ List all resources
- ✅ Get resource details and status
- ✅ Update resource status
- ✅ Delete resources

### AI Guidance (Prompts)

- ✅ IDP capabilities overview
- ✅ Git repository creation guide (using xgitrepo)
- ✅ React+K8s deployment guide (coming soon)
- ✅ General blueprint usage workflows
- ✅ Troubleshooting guides

### Observability

- ✅ Comprehensive logging at all stages
- ✅ Request/response tracking
- ✅ API call monitoring
- ✅ Error diagnostics

## Installation

### Prerequisites

- Go 1.21 or higher
- Access to CloudGenie IDP Backend
- (Optional) Claude Desktop for local testing

### Build from Source

```bash
# Clone the repository
git clone https://github.com/deepakvbansode/idp-cloudgenie-mcp-server.git
cd idp-cloudgenie-mcp-server

# Build the server
go build -o cloudgenie-mcp-server .

# Run in stdio mode (for Claude Desktop)
export CLOUDGENIE_BACKEND_URL="http://your-backend:50051/cloud-genie"
./cloudgenie-mcp-server

# Or run in HTTP mode (for Kubernetes)
export CLOUDGENIE_BACKEND_URL="http://your-backend:50051/cloud-genie"
export MCP_HTTP_PORT=5100
./cloudgenie-mcp-server
```

## Configuration

### Environment Variables

| Variable                 | Required | Default                  | Description                                |
| ------------------------ | -------- | ------------------------ | ------------------------------------------ |
| `CLOUDGENIE_BACKEND_URL` | Yes      | `http://localhost:50051` | CloudGenie IDP backend API URL             |
| `MCP_HTTP_PORT`          | No       | -                        | Enable HTTP mode on this port (e.g., 5100) |
| `MCP_DEBUG`              | No       | `false`                  | Enable verbose debug logging               |

### Claude Desktop Configuration

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "cloudgenie-idp": {
      "command": "/path/to/cloudgenie-mcp-server",
      "env": {
        "CLOUDGENIE_BACKEND_URL": "http://your-backend:50051/cloud-genie"
      }
    }
  }
}
```

## Available Tools

### 1. `cloudgenie_health`

Check CloudGenie backend health status.

**Usage**: "Is CloudGenie backend healthy?"

### 2. `cloudgenie_get_blueprints`

List all available blueprints (currently xgitrepo, more coming).

**Usage**: "What can I create?" or "Show me available blueprints"

### 3. `cloudgenie_get_blueprint`

Get detailed specifications for a specific blueprint.

**Parameters**:

- `id` (required): Blueprint ID (e.g., "xgitrepo")

**Usage**: "Show me details of the xgitrepo blueprint"

### 4. `cloudgenie_create_resource`

Create a new resource from a blueprint.

**Parameters**:

- `blueprint_id` (required): Blueprint to use (e.g., "xgitrepo")
- `name` (required): Resource name
- `type` (required): Resource type (e.g., "git-repository")
- `properties` (optional): Blueprint-specific parameters

**Usage**: "Create a git repository called my-app"

### 5. `cloudgenie_get_resources`

List all created resources.

**Usage**: "Show me all my resources" or "List my Git repositories"

### 6. `cloudgenie_get_resource`

Get detailed information about a specific resource.

**Parameters**:

- `id` (required): Resource ID

**Usage**: "Show me details of resource abc-123"

### 7. `cloudgenie_delete_resource`

Delete a resource.

**Parameters**:

- `id` (required): Resource ID

**Usage**: "Delete the repository abc-123"

### 8. `cloudgenie_update_resource_status`

Update resource status (for recovery/retry).

**Parameters**:

- `id` (required): Resource ID
- `status` (required): New status

**Usage**: "Retry the failed resource abc-123"

## Available Prompts

Prompts provide AI models with comprehensive guidance for complex workflows.

### 1. `cloudgenie_capabilities`

Overview of IDP capabilities, blueprints, and workflows.

**When to use**: User asks "What can you do?" or "How can you help?"

### 2. `create_git_repository`

Complete guide for creating Git repositories using xgitrepo blueprint.

**When to use**: User wants to create a Git repository

**Arguments**:

- `repo_name` (optional): Desired repository name
- `description` (optional): Repository description

### 3. `deploy_application_k8s`

Guide for React+K8s deployment (coming soon).

**When to use**: User wants to deploy an application to Kubernetes

**Arguments**:

- `app_type` (optional): Application type (react, nodejs, etc.)
- `app_name` (optional): Application name

### 4. `deploy_infrastructure`

General guide for using any blueprint to create resources.

**When to use**: User wants to create resources but needs workflow guidance

**Arguments**:

- `blueprint_id` (optional): Blueprint to use
- `resource_name` (optional): Resource name

### 5. `troubleshoot_resource`

Comprehensive troubleshooting guide.

**When to use**: Resource creation fails or gets stuck

**Arguments**:

- `resource_id` (optional): Resource experiencing issues

### 6. `list_blueprints`

Guide for discovering and choosing blueprints.

**When to use**: User wants to explore available options

## Example Conversations

### Creating a Git Repository

```
User: "I need to create a Git repository for my new project"

Claude: "I can help you create a Git repository using the CloudGenie IDP!
Let me check what's available..."

[Claude calls: cloudgenie_get_blueprints]

Claude: "Great! The 'xgitrepo' blueprint is available for creating Git repositories.
What would you like to name your repository?"

User: "Call it awesome-web-app"

Claude: "Perfect! Let me create that for you..."

[Claude calls: cloudgenie_create_resource with blueprint_id="xgitrepo",
 name="awesome-web-app", type="git-repository"]

Claude: "✅ Repository created successfully!
- Name: awesome-web-app
- Status: Ready
- Resource ID: res-xyz-789
- You can now push code to your new repository!"
```

### Discovering Capabilities

```
User: "What can you help me with?"

Claude: [Gets cloudgenie_capabilities prompt]

Claude: "I can help you with the CloudGenie Internal Developer Platform!

Currently available:
✅ Create Git Repositories (using xgitrepo blueprint)
   - Set repository name, visibility, default branch
   - Manage repository lifecycle

Coming soon:
🚧 Deploy React Applications to Kubernetes
   - Automated build and deployment
   - K8s resource management

I can help you:
1. Discover available blueprints
2. Create resources from blueprints
3. Manage resource lifecycle
4. Troubleshoot issues

What would you like to do?"
```

## Kubernetes Deployment

For production deployment in Kubernetes:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloudgenie-mcp-server
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: mcp-server
          image: cloudgenie-mcp-server:latest
          env:
            - name: CLOUDGENIE_BACKEND_URL
              value: "http://cloudgenie-backend:50051/cloud-genie"
            - name: MCP_HTTP_PORT
              value: "5100"
          ports:
            - containerPort: 5100
---
apiVersion: v1
kind: Service
metadata:
  name: cloudgenie-mcp-server
spec:
  selector:
    app: cloudgenie-mcp-server
  ports:
    - port: 5100
      targetPort: 5100
```

AI model pods can then call: `http://cloudgenie-mcp-server:5100/mcp`

## Development

### Project Structure

```
.
├── main.go                 # Entry point, tool/prompt registration
├── go.mod                  # Go module definition
├── mcp/
│   ├── server.go          # MCP protocol server (stdio & HTTP)
│   └── types.go           # MCP protocol type definitions
├── cloudgenie/
│   └── client.go          # CloudGenie API client
├── README.md              # Original documentation
├── README_IDP.md          # This file (IDP-focused)
├── PROMPTS.md             # Detailed prompt documentation
└── IMPLEMENTATION.md      # Technical implementation details
```

### Adding a New Blueprint Support

When a new blueprint (e.g., React+K8s) is added to the backend:

1. **No code changes needed!** The MCP server is blueprint-agnostic
2. Update prompt descriptions in `main.go` to mention the new capability
3. The model will automatically discover it via `get_blueprints`
4. The model will understand it via `get_blueprint(id)`
5. The model will use it via `create_resource`

This is the power of the blueprint-based model! 🎉

### Debugging

```bash
# Enable debug logging
export MCP_DEBUG=true
./cloudgenie-mcp-server

# Test with curl (HTTP mode)
curl -X POST http://localhost:5100/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'

# Test blueprint discovery
curl -X POST http://localhost:5100/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"cloudgenie_get_blueprints","arguments":{}}}'
```

## Roadmap

- [x] xgitrepo blueprint support (Git repositories)
- [ ] React+K8s blueprint support (Application deployment)
- [ ] Database blueprint support
- [ ] Network blueprint support
- [ ] Monitoring and alerting blueprints
- [ ] Blueprint composition (multi-resource blueprints)
- [ ] Resource dependencies and orchestration

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license here]

## Support

For issues, questions, or feature requests:

- GitHub Issues: [Create an issue](https://github.com/deepakvbansode/idp-cloudgenie-mcp-server/issues)
- Documentation: See PROMPTS.md and IMPLEMENTATION.md

## Acknowledgments

- Built with the [Model Context Protocol](https://modelcontextprotocol.io/)
- Designed for [Anthropic Claude](https://www.anthropic.com/claude) and other AI assistants
- Part of the CloudGenie Internal Developer Platform ecosystem
