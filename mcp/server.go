package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Server represents an MCP server
type Server struct {
	name         string
	version      string
	capabilities ServerCapabilities
	tools        map[string]Tool
	toolHandlers map[string]ToolHandler
	resources    map[string]Resource
	prompts      map[string]Prompt
	reader       *bufio.Reader
	writer       io.Writer
	httpSrv      *http.Server
}

// ToolHandler is a function that executes a tool
type ToolHandler func(arguments map[string]interface{}) ([]Content, error)

// NewServer creates a new MCP server
func NewServer(name, version string) *Server {
       log.Printf("[MCP] Creating new MCP server: %s v%s", name, version)
       return &Server{
	       name:    name,
	       version: version,
	       capabilities: ServerCapabilities{
		       Tools:     &ToolsCapability{},
		       Resources: &ResourcesCapability{},
		       Prompts:   &PromptsCapability{},
	       },
	       tools:        make(map[string]Tool),
	       toolHandlers: make(map[string]ToolHandler),
	       resources:    make(map[string]Resource),
	       prompts:      make(map[string]Prompt),
	       reader:       bufio.NewReader(os.Stdin),
	       writer:       os.Stdout,
       }
}

// StartHTTP starts the MCP server as an HTTP server on the given address (e.g., ":8080")
func (s *Server) StartHTTP(addr string) error {
       log.Printf("[MCP] Initializing HTTP server on %s", addr)
       mux := http.NewServeMux()
       mux.HandleFunc("/mcp", s.httpHandler)
       s.httpSrv = &http.Server{Addr: addr, Handler: mux}
       log.Printf("[MCP] Starting MCP HTTP server on %s", addr)
       log.Printf("[MCP] Registered %d tools, %d resources, %d prompts", len(s.tools), len(s.resources), len(s.prompts))
       return s.httpSrv.ListenAndServe()
}

// httpHandler handles JSON-RPC 2.0 requests over HTTP POST
func (s *Server) httpHandler(w http.ResponseWriter, r *http.Request) {
       log.Printf("[MCP HTTP] Received request from %s: %s %s", r.RemoteAddr, r.Method, r.URL.Path)
       
       if r.Method != http.MethodPost {
	       log.Printf("[MCP HTTP] Rejected non-POST request")
	       http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	       return
       }
       
       defer r.Body.Close()
       var msg JSONRPCMessage
       decoder := json.NewDecoder(r.Body)
       if err := decoder.Decode(&msg); err != nil {
	       log.Printf("[MCP HTTP] Failed to decode JSON-RPC message: %v", err)
	       http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
	       return
       }

       log.Printf("[MCP HTTP] Processing JSON-RPC method: %s (ID: %v)", msg.Method, msg.ID)

       // Capture response
       var respData []byte
       respWriter := &responseBuffer{buf: &respData}
       origWriter := s.writer
       s.writer = respWriter
       s.handleRequest(&msg)
       s.writer = origWriter

       log.Printf("[MCP HTTP] Sending response for method: %s (size: %d bytes)", msg.Method, len(respData))
       w.Header().Set("Content-Type", "application/json")
       w.Write(respData)
}

// responseBuffer is a helper to capture JSON-RPC responses for HTTP
type responseBuffer struct {
       buf *[]byte
}

func (rb *responseBuffer) Write(p []byte) (int, error) {
       *rb.buf = append(*rb.buf, p...)
       return len(p), nil
}

// RegisterTool registers a tool with its handler
func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
	log.Printf("[MCP] Registering tool: %s", tool.Name)
	s.tools[tool.Name] = tool
	s.toolHandlers[tool.Name] = handler
}

// RegisterResource registers a resource
func (s *Server) RegisterResource(resource Resource) {
	log.Printf("[MCP] Registering resource: %s (%s)", resource.Name, resource.URI)
	s.resources[resource.URI] = resource
}

// RegisterPrompt registers a prompt
func (s *Server) RegisterPrompt(prompt Prompt) {
	log.Printf("[MCP] Registering prompt: %s", prompt.Name)
	s.prompts[prompt.Name] = prompt
}

// Start starts the MCP server
func (s *Server) Start() error {
	log.Printf("[MCP] Starting MCP server (stdio mode): %s v%s", s.name, s.version)
	log.Printf("[MCP] Registered %d tools, %d resources, %d prompts", len(s.tools), len(s.resources), len(s.prompts))

	for {
		// Read a line from stdin
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("[MCP] EOF received, shutting down")
				return nil
			}
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		log.Printf("[MCP STDIO] Received message: %s", line)

		// Parse JSON-RPC message
		var msg JSONRPCMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			log.Printf("[MCP STDIO] Parse error: %v", err)
			s.sendError(nil, -32700, "Parse error", nil)
			continue
		}

		log.Printf("[MCP STDIO] Processing method: %s (ID: %v)", msg.Method, msg.ID)
		// Handle the request
		s.handleRequest(&msg)
	}
}

func (s *Server) handleRequest(msg *JSONRPCMessage) {
	log.Printf("[MCP] Handling request - Method: %s, ID: %v", msg.Method, msg.ID)
	
	switch msg.Method {
	case "initialize":
		s.handleInitialize(msg)
	case "initialized":
		// Client notification that initialization is complete
		log.Println("[MCP] Client initialized successfully")
	case "tools/list":
		s.handleListTools(msg)
	case "tools/call":
		s.handleCallTool(msg)
	case "resources/list":
		s.handleListResources(msg)
	case "resources/read":
		s.handleReadResource(msg)
	case "prompts/list":
		s.handleListPrompts(msg)
	case "prompts/get":
		s.handleGetPrompt(msg)
	case "ping":
		s.handlePing(msg)
	default:
		log.Printf("[MCP] Unknown method: %s", msg.Method)
		s.sendError(msg.ID, -32601, "Method not found", nil)
	}
}

func (s *Server) handleInitialize(msg *JSONRPCMessage) {
	log.Printf("[MCP] Processing initialize request")
	var params InitializeParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		log.Printf("[MCP] Initialize failed - invalid params: %v", err)
		s.sendError(msg.ID, -32602, "Invalid params", nil)
		return
	}

	log.Printf("[MCP] Client info: %s v%s, Protocol: %s", params.ClientInfo.Name, params.ClientInfo.Version, params.ProtocolVersion)

	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities:    s.capabilities,
		ServerInfo: ServerInfo{
			Name:    s.name,
			Version: s.version,
		},
	}

	log.Printf("[MCP] Initialization successful - Server: %s v%s", s.name, s.version)
	s.sendResult(msg.ID, result)
}

func (s *Server) handleListTools(msg *JSONRPCMessage) {
	log.Printf("[MCP] Listing tools - Total: %d", len(s.tools))
	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}

	result := ListToolsResult{
		Tools: tools,
	}

	log.Printf("[MCP] Returning %d tools", len(tools))
	s.sendResult(msg.ID, result)
}

func (s *Server) handleCallTool(msg *JSONRPCMessage) {
	log.Printf("[MCP] Processing tool call request")
	var params CallToolParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		log.Printf("[MCP] Tool call failed - invalid params: %v", err)
		s.sendError(msg.ID, -32602, "Invalid params", nil)
		return
	}

	log.Printf("[MCP] Calling tool: %s with arguments: %v", params.Name, params.Arguments)

	handler, exists := s.toolHandlers[params.Name]
	if !exists {
		log.Printf("[MCP] Tool not found: %s", params.Name)
		s.sendError(msg.ID, -32602, fmt.Sprintf("Tool not found: %s", params.Name), nil)
		return
	}

	log.Printf("[MCP] Executing tool handler: %s", params.Name)
	content, err := handler(params.Arguments)
	if err != nil {
		log.Printf("[MCP] Tool execution failed: %s - Error: %v", params.Name, err)
		result := CallToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Error: %s", err.Error()),
			}},
			IsError: true,
		}
		s.sendResult(msg.ID, result)
		return
	}

	log.Printf("[MCP] Tool execution successful: %s", params.Name)
	result := CallToolResult{
		Content: content,
		IsError: false,
	}

	s.sendResult(msg.ID, result)
}

func (s *Server) handleListResources(msg *JSONRPCMessage) {
	log.Printf("[MCP] Listing resources - Total: %d", len(s.resources))
	resources := make([]Resource, 0, len(s.resources))
	for _, resource := range s.resources {
		resources = append(resources, resource)
	}

	result := ListResourcesResult{
		Resources: resources,
	}

	log.Printf("[MCP] Returning %d resources", len(resources))
	s.sendResult(msg.ID, result)
}

func (s *Server) handleReadResource(msg *JSONRPCMessage) {
	log.Printf("[MCP] Processing read resource request")
	var params ReadResourceParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		log.Printf("[MCP] Read resource failed - invalid params: %v", err)
		s.sendError(msg.ID, -32602, "Invalid params", nil)
		return
	}

	log.Printf("[MCP] Reading resource: %s", params.URI)

	resource, exists := s.resources[params.URI]
	if !exists {
		log.Printf("[MCP] Resource not found: %s", params.URI)
		s.sendError(msg.ID, -32602, fmt.Sprintf("Resource not found: %s", params.URI), nil)
		return
	}

	// For this example, we'll return the resource description as content
	result := ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      resource.URI,
			MimeType: resource.MimeType,
			Text:     resource.Description,
		}},
	}

	log.Printf("[MCP] Resource read successful: %s", params.URI)
	s.sendResult(msg.ID, result)
}

func (s *Server) handleListPrompts(msg *JSONRPCMessage) {
	log.Printf("[MCP] Listing prompts - Total: %d", len(s.prompts))
	prompts := make([]Prompt, 0, len(s.prompts))
	for _, prompt := range s.prompts {
		prompts = append(prompts, prompt)
	}

	result := ListPromptsResult{
		Prompts: prompts,
	}

	log.Printf("[MCP] Returning %d prompts", len(prompts))
	s.sendResult(msg.ID, result)
}

func (s *Server) handleGetPrompt(msg *JSONRPCMessage) {
	log.Printf("[MCP] Processing get prompt request")
	var params GetPromptParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		log.Printf("[MCP] Get prompt failed - invalid params: %v", err)
		s.sendError(msg.ID, -32602, "Invalid params", nil)
		return
	}

	log.Printf("[MCP] Getting prompt: %s", params.Name)

	prompt, exists := s.prompts[params.Name]
	if !exists {
		log.Printf("[MCP] Prompt not found: %s", params.Name)
		s.sendError(msg.ID, -32602, fmt.Sprintf("Prompt not found: %s", params.Name), nil)
		return
	}

	result := GetPromptResult{
		Description: prompt.Description,
		Messages: []PromptMessage{{
			Role: "user",
			Content: PromptContent{
				Type: "text",
				Text: fmt.Sprintf("Prompt: %s", prompt.Name),
			},
		}},
	}

	log.Printf("[MCP] Prompt retrieved successfully: %s", params.Name)
	s.sendResult(msg.ID, result)
}

func (s *Server) handlePing(msg *JSONRPCMessage) {
	log.Printf("[MCP] Ping received")
	s.sendResult(msg.ID, map[string]interface{}{})
}

func (s *Server) sendResult(id interface{}, result interface{}) {
	response := JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	log.Printf("[MCP] Sending result for request ID: %v", id)
	s.sendMessage(response)
}

func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	response := JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	log.Printf("[MCP] Sending error for request ID: %v - Code: %d, Message: %s", id, code, message)
	s.sendMessage(response)
}

func (s *Server) sendMessage(msg JSONRPCMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	_, err = fmt.Fprintf(s.writer, "%s\n", data)
	if err != nil {
		log.Printf("Error writing message: %v", err)
	}
}
