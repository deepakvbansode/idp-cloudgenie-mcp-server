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
       mux := http.NewServeMux()
       mux.HandleFunc("/mcp", s.httpHandler)
       s.httpSrv = &http.Server{Addr: addr, Handler: mux}
       log.Printf("Starting MCP HTTP server on %s", addr)
       return s.httpSrv.ListenAndServe()
}

// httpHandler handles JSON-RPC 2.0 requests over HTTP POST
func (s *Server) httpHandler(w http.ResponseWriter, r *http.Request) {
       if r.Method != http.MethodPost {
	       http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	       return
       }
	   fmt.Println("Received HTTP MCP request")
       defer r.Body.Close()
       var msg JSONRPCMessage
       decoder := json.NewDecoder(r.Body)
       if err := decoder.Decode(&msg); err != nil {
	       http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
	       return
       }

       // Capture response
       var respData []byte
       respWriter := &responseBuffer{buf: &respData}
       origWriter := s.writer
       s.writer = respWriter
       s.handleRequest(&msg)
       s.writer = origWriter

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
	s.tools[tool.Name] = tool
	s.toolHandlers[tool.Name] = handler
}

// RegisterResource registers a resource
func (s *Server) RegisterResource(resource Resource) {
	s.resources[resource.URI] = resource
}

// RegisterPrompt registers a prompt
func (s *Server) RegisterPrompt(prompt Prompt) {
	s.prompts[prompt.Name] = prompt
}

// Start starts the MCP server
func (s *Server) Start() error {
	log.Printf("Starting MCP server: %s v%s", s.name, s.version)

	for {
		// Read a line from stdin
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		// Parse JSON-RPC message
		var msg JSONRPCMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			s.sendError(nil, -32700, "Parse error", nil)
			continue
		}

		// Handle the request
		s.handleRequest(&msg)
	}
}

func (s *Server) handleRequest(msg *JSONRPCMessage) {
	switch msg.Method {
	case "initialize":
		s.handleInitialize(msg)
	case "initialized":
		// Client notification that initialization is complete
		log.Println("Client initialized")
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
		s.sendError(msg.ID, -32601, "Method not found", nil)
	}
}

func (s *Server) handleInitialize(msg *JSONRPCMessage) {
	var params InitializeParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.sendError(msg.ID, -32602, "Invalid params", nil)
		return
	}

	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities:    s.capabilities,
		ServerInfo: ServerInfo{
			Name:    s.name,
			Version: s.version,
		},
	}

	s.sendResult(msg.ID, result)
}

func (s *Server) handleListTools(msg *JSONRPCMessage) {
	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}

	result := ListToolsResult{
		Tools: tools,
	}

	s.sendResult(msg.ID, result)
}

func (s *Server) handleCallTool(msg *JSONRPCMessage) {
	var params CallToolParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.sendError(msg.ID, -32602, "Invalid params", nil)
		return
	}

	handler, exists := s.toolHandlers[params.Name]
	if !exists {
		s.sendError(msg.ID, -32602, fmt.Sprintf("Tool not found: %s", params.Name), nil)
		return
	}

	content, err := handler(params.Arguments)
	if err != nil {
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

	result := CallToolResult{
		Content: content,
		IsError: false,
	}

	s.sendResult(msg.ID, result)
}

func (s *Server) handleListResources(msg *JSONRPCMessage) {
	resources := make([]Resource, 0, len(s.resources))
	for _, resource := range s.resources {
		resources = append(resources, resource)
	}

	result := ListResourcesResult{
		Resources: resources,
	}

	s.sendResult(msg.ID, result)
}

func (s *Server) handleReadResource(msg *JSONRPCMessage) {
	var params ReadResourceParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.sendError(msg.ID, -32602, "Invalid params", nil)
		return
	}

	resource, exists := s.resources[params.URI]
	if !exists {
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

	s.sendResult(msg.ID, result)
}

func (s *Server) handleListPrompts(msg *JSONRPCMessage) {
	prompts := make([]Prompt, 0, len(s.prompts))
	for _, prompt := range s.prompts {
		prompts = append(prompts, prompt)
	}

	result := ListPromptsResult{
		Prompts: prompts,
	}

	s.sendResult(msg.ID, result)
}

func (s *Server) handleGetPrompt(msg *JSONRPCMessage) {
	var params GetPromptParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.sendError(msg.ID, -32602, "Invalid params", nil)
		return
	}

	prompt, exists := s.prompts[params.Name]
	if !exists {
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

	s.sendResult(msg.ID, result)
}

func (s *Server) handlePing(msg *JSONRPCMessage) {
	s.sendResult(msg.ID, map[string]interface{}{})
}

func (s *Server) sendResult(id interface{}, result interface{}) {
	response := JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

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
