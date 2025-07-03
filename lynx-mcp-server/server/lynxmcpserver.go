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

	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/config"
	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/tools"
	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/utils"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	PORT = ":9600"
)

func main() {
	// Load configuration
	serverConfig := config.NewMCPServerConfig()

	log.Printf("Starting SSE server on %s/see", serverConfig.Port)

	mcpServer := NewMCPServer()
	sse := server.NewSSEServer(mcpServer)

	// Create custom HTTP server with BearerAuthMiddleware
	httpServer := &http.Server{
		Handler: utils.BearerAuthMiddleware(serverConfig.BearerToken)(sse),
	}

	// Use WithHTTPServer to inject our custom server
	sse = server.NewSSEServer(mcpServer, server.WithHTTPServer(httpServer))

	// Create a channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		serverErrors <- sse.Start(":" + serverConfig.Port)
	}()

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for either an error or a signal
	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v, shutting down server...", sig)
		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := sse.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
		log.Println("Server shutdown complete")
	}
}

func NewMCPServer() *server.MCPServer {
	hooks := &server.Hooks{}

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		fmt.Printf("beforeAny: %s, %v, %v\n", method, id, message)
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		fmt.Printf("onSuccess: %s, %v, %v, %v\n", method, id, message, result)
	})
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		fmt.Printf("onError: %s, %v, %v, %v\n", method, id, message, err)
	})
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		fmt.Printf("beforeInitialize: %v, %v\n", id, message)
	})
	hooks.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
		fmt.Printf("AddOnRequestInitialization: %v\n", id)
		return nil
	})
	hooks.AddAfterInitialize(func(_ context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
		fmt.Printf("afterInitialize: %v, %v, %v\n", id, message, result)
	})
	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		fmt.Printf("afterCallTool: %v, %v, %v\n", id, message, result)
	})
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		fmt.Printf("beforeCallTool: %v, %v\n", id, message)
	})

	mcpServer := server.NewMCPServer(
		"lynx-mcp-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(hooks),
	)

	mcpServer.AddTool(mcp.NewToolWithRawSchema(
		string(tools.TOOL_FILE_SEARCH_BY_PARTY_NAME),
		tools.TOOL_FILE_SEARCH_BY_PARTY_NAME_DESCRIPTION,
		tools.GetFileSearchSchema(),
	), tools.HandleFileSearchByPartyName)

	return mcpServer
}
