# CloudGenie IDP MCP Server - Prompts Guide

## Overview

Prompts in the Model Context Protocol (MCP) are templates that help AI models understand how to interact with your IDP (Internal Developer Platform) and guide users through creating resources via blueprints. This document explains all registered prompts and when they're used.

## IDP Architecture: Blueprints & Resources

**CloudGenie IDP** follows a blueprint-based model:

- **Blueprints**: Define what types of resources can be created (e.g., `xgitrepo` for Git repositories)
- **Resources**: Instances of blueprints (e.g., a specific Git repository created from `xgitrepo` blueprint)

### Current Capabilities

- ✅ **xgitrepo blueprint**: Create Git repositories
- 🚧 **Coming Soon**: React app + K8s deployment blueprint

### Workflow

1. Model calls `get_blueprints` to see what can be created
2. Model examines blueprint parameters to understand requirements
3. User specifies what they want to create
4. Model calls `create_resource` with blueprint ID and parameters
5. Model monitors resource status

## Understanding Prompts vs Tools

### Tools (`tools/call`)

- **Purpose**: Execute actions or retrieve dynamic data
- **When used**: Model needs to DO something or GET live data
- **Examples**: Create resources, list blueprints, check health
- **Side effects**: Yes (can modify system state)

### Prompts (`prompts/get`)

- **Purpose**: Provide guidance, templates, and instructions
- **When used**: Model needs CONTEXT on how to approach a task
- **Examples**: How to deploy apps, troubleshooting guides
- **Side effects**: No (read-only guidance)

## Registered Prompts

### 1. `cloudgenie_capabilities`

**When user asks**: "How can you help me?" or "What can you do?"

**Description**: Comprehensive overview of CloudGenie IDP capabilities, available blueprints, resource management, and common workflows.

**Arguments**: None

**Use case**: When users need to understand what the IDP can do and what resources can be created.

**What the model will learn**:

- IDP blueprint-based model (blueprints define resource types)
- Currently available: `xgitrepo` blueprint for creating Git repositories
- Coming soon: React app + K8s deployment blueprint
- How to discover blueprints and create resources
- Resource lifecycle management
- Common workflows and best practices

---

### 2. `deploy_application_k8s`

**When user asks**: "Deploy a React application to Kubernetes" or "How do I deploy my Node.js app?"

**Description**: Guide for deploying applications to Kubernetes. Currently explains that the React+K8s blueprint is coming soon. When available, will show how to use the blueprint to deploy containerized applications.

**Arguments**:

- `app_type` (optional): Type of application (e.g., react, nodejs, python, java)
- `app_name` (optional): Name for the application deployment

**Use case**: When users want to deploy applications to Kubernetes, but the blueprint isn't available yet.

**What the model will learn**:

- React+K8s deployment blueprint is coming soon
- In the meantime, can deploy K8s cluster infrastructure
- Once blueprint is available, will follow same pattern as xgitrepo
- Alternative approaches until blueprint is ready
- How to prepare for the upcoming capability

**Example usage**:

```text
User: "Deploy my React app to Kubernetes"
Model: Gets prompt, explains React+K8s blueprint is coming soon
Model: Offers to help with other IDP capabilities (create Git repo, etc.)
Model: Explains workflow for when blueprint becomes available
```

---

### 3. `create_git_repository`

**When user asks**: "Create a Git repository" or "I need a new repo"

**Description**: Complete guide for creating Git repositories using the `xgitrepo` blueprint. This is a **working feature** in the IDP.

**Arguments**:

- `repo_name` (optional): Desired name for the Git repository
- `description` (optional): Description of the repository purpose

**Use case**: When users want to create a new Git repository through the IDP.

**What the model will learn**:

- How to use `cloudgenie_get_blueprints` to find the `xgitrepo` blueprint
- How to examine the blueprint to understand required parameters
- How to call `cloudgenie_create_resource` with:
  - `blueprint_id`: "xgitrepo"
  - `name`: User's desired repo name
  - `type`: "git-repository" (or as defined by blueprint)
  - `properties`: Any additional parameters from the blueprint
- How to monitor resource creation status
- What to do if creation fails

**Example usage**:

```text
User: "Create a git repository called my-awesome-app"
Model: Gets prompt with repo_name="my-awesome-app"
Model: Calls cloudgenie_get_blueprints to find xgitrepo
Model: Calls cloudgenie_get_blueprint("xgitrepo") to see parameters
Model: Calls cloudgenie_create_resource with blueprint_id="xgitrepo", name="my-awesome-app"
Model: Monitors status and confirms creation
```

---

### 4. `deploy_infrastructure`

**When user asks**: "Create a resource using a blueprint" or "Deploy infrastructure"

**Description**: Step-by-step guide for creating resources using IDP blueprints. Explains the general workflow applicable to any blueprint.

**Arguments**:

- `blueprint_id` (optional): ID of the blueprint to use (e.g., xgitrepo)
- `resource_name` (optional): Name for the resource to create

**Use case**: When users want to create any type of resource through blueprints.

**What the model will learn**:

- IDP follows a blueprint → resource model
- First, discover available blueprints with `cloudgenie_get_blueprints`
- Examine specific blueprint with `cloudgenie_get_blueprint(id)`
- Understand required parameters from blueprint specification
- Create resource instance with `cloudgenie_create_resource`
- Monitor resource creation status
- Troubleshooting if creation fails

**Example usage**:

```text
User: "Use the xgitrepo blueprint to create a repository"
Model: Gets prompt with blueprint_id="xgitrepo"
Model: Calls cloudgenie_get_blueprint("xgitrepo") to understand parameters
Model: Guides user through parameter collection
Model: Calls cloudgenie_create_resource to create instance
```

---

### 5. `troubleshoot_resource`

**When user asks**: "My resource creation failed" or "Why is my Git repo stuck?"

**Description**: Comprehensive troubleshooting guide for IDP resources. Covers diagnostics and recovery.

**Arguments**:

- `resource_id` (optional): ID of the resource experiencing issues

**Use case**: When resource creation fails or gets stuck.

**What the model will learn**:

- How to check resource status with `cloudgenie_get_resource`
- Common status meanings:
  - `pending`: Still creating, normal for a few minutes
  - `failed`: Creation failed, check error message
  - `error`: Parameter or dependency issue
  - `ready`: Successfully created
- Backend health verification
- Recovery actions (retry vs recreate)
- When to escalate

**Example usage**:

```text
User: "My Git repository is stuck in pending"
Model: Gets prompt with resource_id
Model: Calls cloudgenie_get_resource to check status
Model: Analyzes status, checks backend health
Model: Suggests appropriate action
```

---

### 6. `list_blueprints`

**When user asks**: "What can I create?" or "Show me available blueprints"

**Description**: Guide for discovering IDP blueprints. Currently shows `xgitrepo` blueprint, with more coming.

**Arguments**: None

**Use case**: When users want to see what types of resources can be created.

**What the model will learn**:

- Blueprints define what resource types can be created
- Currently available: `xgitrepo` (Git repositories)
- Coming soon: React+K8s deployment blueprint
- How to list all blueprints with `cloudgenie_get_blueprints`
- How to examine blueprint details with `cloudgenie_get_blueprint(id)`
- Understanding blueprint parameters and requirements
- Choosing the right blueprint for their needs

**Example usage**:

```text
User: "What can I create with the IDP?"
Model: Gets prompt
Model: Calls cloudgenie_get_blueprints
Model: Shows xgitrepo blueprint (Git repos available)
Model: Mentions React+K8s coming soon
Model: Helps user choose and use appropriate blueprint
```

---

## How AI Models Use These Prompts

### Typical Workflow

1. **User makes a request**

   - Example: "Deploy a React app to Kubernetes"

2. **Model recognizes the need for guidance**

   - Determines which prompt is relevant: `deploy_application_k8s`

3. **Model retrieves the prompt**

   - Calls `prompts/get` with name="deploy_application_k8s"
   - Passes arguments: app_type="react", app_name (if mentioned)

4. **Model receives comprehensive guidance**

   - Gets step-by-step instructions
   - Understands prerequisites
   - Knows which tools to call in which order

5. **Model follows the guidance**

   - Checks for K8s cluster (calls `cloudgenie_get_resources`)
   - Lists blueprints (calls `cloudgenie_get_blueprints`)
   - Creates resource (calls `cloudgenie_create_resource`)
   - Monitors status (calls `cloudgenie_get_resource`)

6. **Model provides helpful response to user**
   - With context from prompt, explains what's happening
   - Guides user through process
   - Handles issues using troubleshooting guidance

## Prompt Design Principles

### Current Implementation

The prompts are **registered as metadata only**. The actual detailed guidance is described in the `Description` field, which helps the model understand:

- **When to use the prompt**: What user queries trigger it
- **What information it provides**: What guidance the user will receive
- **What the model should do**: How to apply the guidance

### Future Enhancement

Currently, prompts return simple static content. For more dynamic responses, you could enhance the server to:

1. Add prompt handlers that generate custom content based on arguments
2. Query backend for dynamic information (available blueprints, current resources)
3. Personalize guidance based on user's current infrastructure state

Example enhancement in `mcp/server.go`:

```go
type PromptHandler func(args map[string]interface{}) ([]PromptMessage, error)

func (s *Server) RegisterPromptWithHandler(prompt Prompt, handler PromptHandler) {
    s.prompts[prompt.Name] = prompt
    s.promptHandlers[prompt.Name] = handler
}
```

## Testing Prompts

### Using the Test Script

```bash
# Test capabilities prompt
echo '{"jsonrpc":"2.0","id":1,"method":"prompts/get","params":{"name":"cloudgenie_capabilities"}}' | ./cloudgenie-mcp-server

# Test deployment guide
echo '{"jsonrpc":"2.0","id":1,"method":"prompts/get","params":{"name":"deploy_application_k8s","arguments":{"app_type":"react","app_name":"my-app"}}}' | ./cloudgenie-mcp-server

# List all prompts
echo '{"jsonrpc":"2.0","id":1,"method":"prompts/list","params":{}}' | ./cloudgenie-mcp-server
```

### From Claude Desktop

Add to your Claude Desktop config:

```json
{
  "mcpServers": {
    "cloudgenie": {
      "command": "/path/to/cloudgenie-mcp-server",
      "env": {
        "CLOUDGENIE_BACKEND_URL": "http://your-backend:50051/cloud-genie"
      }
    }
  }
}
```

Then ask Claude:

- "How can CloudGenie help me?"
- "Deploy my React app to Kubernetes"
- "Create a Git repository for my infrastructure"

## Summary

| Prompt Name               | Triggers On            | Provides Guidance For                       |
| ------------------------- | ---------------------- | ------------------------------------------- |
| `cloudgenie_capabilities` | "What can you do?"     | Overview of IDP capabilities and blueprints |
| `deploy_application_k8s`  | "Deploy [app] to K8s"  | React+K8s blueprint (coming soon)           |
| `create_git_repository`   | "Create Git repo"      | Using xgitrepo blueprint to create repos    |
| `deploy_infrastructure`   | "Create resource"      | General blueprint → resource workflow       |
| `troubleshoot_resource`   | "Resource not working" | Diagnostics and recovery                    |
| `list_blueprints`         | "What can I create?"   | Blueprint discovery (xgitrepo available)    |

## IDP Workflow Example: Creating a Git Repository

Here's how the model uses blueprints and tools together:

**User**: "Create a Git repository called my-awesome-app"

**Model's Actions**:

1. **Gets prompt guidance**:

   ```
   Model calls: prompts/get with name="create_git_repository"
   Arguments: {repo_name: "my-awesome-app"}
   Result: Comprehensive guide on using xgitrepo blueprint
   ```

2. **Discovers available blueprints**:

   ```
   Model calls: cloudgenie_get_blueprints tool
   Result: List including xgitrepo blueprint
   ```

3. **Examines xgitrepo blueprint**:

   ```
   Model calls: cloudgenie_get_blueprint tool with id="xgitrepo"
   Result: Blueprint specification with required parameters:
     - name (required): Repository name
     - description (optional): Repository description
     - visibility (optional): public/private
     - default_branch (optional): main/master
   ```

4. **Creates the Git repository resource**:

   ```
   Model calls: cloudgenie_create_resource tool
   Arguments: {
     blueprint_id: "xgitrepo",
     name: "my-awesome-app",
     type: "git-repository",
     properties: {
       description: "My awesome application",
       visibility: "private",
       default_branch: "main"
     }
   }
   Result: Resource created with ID, status: "pending"
   ```

5. **Monitors resource status**:

   ```
   Model calls: cloudgenie_get_resource tool with resource_id
   Result: Status transitions: pending → creating → ready
   ```

6. **Confirms to user**:
   ```
   "✅ Git repository 'my-awesome-app' created successfully!
   - Repository ID: abc-123-def
   - Status: Ready
   - URL: https://git.example.com/my-awesome-app
   - Default branch: main"
   ```

This same workflow applies to any future blueprint (React+K8s, databases, etc.)!
