package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/prompts"
	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/resources"
	"github.com/deepakvbansode/idp-cloudgenie-mcp-server/cloudgenie/tools"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("==========================================")
	log.Println("CloudGenie MCP Server (SDK) Starting...")
	log.Println("==========================================")

	backendURL := os.Getenv("CLOUDGENIE_BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:50051/cloud-genie"
	}
	log.Printf("[MAIN] CloudGenie backend URL: %s", backendURL)

	client := client.NewCGClient(backendURL)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "idp-cloudgenie-mcp-server",
		Version: "1.0.0",
	}, nil)

	prompts.RegisterPrompts(server, client)
	resources.RegisterResources(server, client)
	tools.RegisterTools(server, client)

	log.Println("[MAIN] CloudGenie MCP Server initialized successfully (SDK)")
	log.Printf("[MAIN] Backend URL: %s", backendURL)

	// Get HTTP port from environment or use default
	port := os.Getenv("MCP_HTTP_PORT")
	if port == "" {
		port = "5100"
	}
	
	// Start server with HTTP transport using StreamableHTTPHandler
	log.Printf("[MAIN] Starting HTTP server on port %s", port)
	log.Println("==========================================")
	
	// Create HTTP handler that serves the MCP server
	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server
		},
		nil, // Use default options
	)
	
	// Set up HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
	
	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		
		log.Println("\n[MAIN] Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("[MAIN] Server shutdown error: %v", err)
		}
	}()
	
	log.Printf("[MAIN] MCP server listening on http://0.0.0.0:%s", port)
	log.Println("[MAIN] Ready to accept connections")
	
	// Start the HTTP server
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[MAIN] Server error: %v", err)
	}
	
	log.Println("[MAIN] Server stopped")
}





