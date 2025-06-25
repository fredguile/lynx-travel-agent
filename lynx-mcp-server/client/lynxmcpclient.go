package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/config"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

var clientConfig = config.NewClientConfig()

func main() {
	sseURL := flag.String("url", "http://127.0.0.1:9600/sse", "URL for SSE transport (e.g. 'http://127.0.0.1:9600/sse')")
	toolCommand := flag.String("command", "", "Command of a remote tool to execute")
	flag.Parse()

	if *sseURL == "" {
		fmt.Println("Error: You must specify SSE URL using --url")
		flag.Usage()
		os.Exit(1)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Initializing SSE client with authentication...")

	// Prepare authentication headers
	options := []transport.ClientOption{}
	if clientConfig.BearerToken != "" {
		headers := map[string]string{
			"Authorization": "Bearer " + clientConfig.BearerToken,
		}
		options = append(options, transport.WithHeaders(headers))
	}

	// Create SSE transport with authentication headers
	c, err := client.NewSSEMCPClient(*sseURL, options...)
	if err != nil {
		log.Fatalf("Failed to create SSE client: %v", err)
	}

	// Start the client
	if err := c.Start(ctx); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Go SSE Client Example",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Display server information
	fmt.Printf("Connected to server: %s (version %s)\n",
		serverInfo.ServerInfo.Name,
		serverInfo.ServerInfo.Version)
	fmt.Printf("Server capabilities: %+v\n", serverInfo.Capabilities)

	// List available tools if the server supports them
	if serverInfo.Capabilities.Tools != nil {
		fmt.Println("Fetching available tools...")
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := c.ListTools(ctx, toolsRequest)
		if err != nil {
			log.Printf("Failed to list tools: %v", err)
		} else {
			fmt.Printf("Server has %d tools available\n", len(toolsResult.Tools))
			for i, tool := range toolsResult.Tools {
				fmt.Printf("  %d. %s - %s\n", i+1, tool.Name, tool.Description)
			}
		}
	}

	// List available resources if the server supports them
	if serverInfo.Capabilities.Resources != nil {
		fmt.Println("Fetching available resources...")
		resourcesRequest := mcp.ListResourcesRequest{}
		resourcesResult, err := c.ListResources(ctx, resourcesRequest)
		if err != nil {
			log.Printf("Failed to list resources: %v", err)
		} else {
			fmt.Printf("Server has %d resources available\n", len(resourcesResult.Resources))
			for i, resource := range resourcesResult.Resources {
				fmt.Printf("  %d. %s - %s\n", i+1, resource.URI, resource.Name)
			}
		}
	}

	if *toolCommand != "" {
		fmt.Println("Executing tool command...")
		toolArgs := parseCommand(*toolCommand)

		if len(toolArgs) == 0 {
			fmt.Println("Error: Invalid command")
			os.Exit(1)
		}

		// First argument should be the tool name
		toolName, ok := toolArgs[0].(string)
		if !ok {
			fmt.Println("Error: First argument must be the tool name")
			os.Exit(1)
		}
		toolArgs = toolArgs[1:]

		// Convert arguments to a map
		args := make(map[string]interface{})
		for i, arg := range toolArgs {
			switch v := arg.(type) {
			case string:
				args[fmt.Sprintf("arg%d", i)] = v
			case map[string]string:
				// Convert map[string]string to map[string]interface{}
				for k, val := range v {
					args[k] = val
				}
			}
		}

		fmt.Printf("Calling tool named \"%s\" with args: %+v\n", toolName, args)

		toolRequest := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name:      toolName,
				Arguments: args,
			},
		}

		result, err := c.CallTool(ctx, toolRequest)
		if err != nil {
			log.Printf("Failed to execute tool: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Tool execution result: %+v\n", result)
	}

	fmt.Println("Client initialized successfully. Shutting down...")
	c.Close()
}

func parseCommand(cmd string) []interface{} {
	var result []interface{}
	var current string
	var inQuote bool
	var quoteChar rune

	for i := 0; i < len(cmd); i++ {
		r := rune(cmd[i])
		switch {
		case r == ' ' && !inQuote:
			if current != "" {
				// Check if the current string is a named argument
				if len(current) > 2 && current[:2] == "--" {
					// Split on first '=' if it exists
					if idx := strings.Index(current, "="); idx != -1 {
						argName := current[2:idx]
						argValue := current[idx+1:]
						result = append(result, map[string]string{argName: argValue})
					} else {
						result = append(result, current)
					}
				} else {
					result = append(result, current)
				}
				current = ""
			}
		case (r == '"' || r == '\''):
			if inQuote && r == quoteChar {
				inQuote = false
				quoteChar = 0
			} else if !inQuote {
				inQuote = true
				quoteChar = r
			} else {
				current += string(r)
			}
		default:
			current += string(r)
		}
	}

	if current != "" {
		// Check if the last current string is a named argument
		if len(current) > 2 && current[:2] == "--" {
			if idx := strings.Index(current, "="); idx != -1 {
				argName := current[2:idx]
				argValue := current[idx+1:]
				result = append(result, map[string]string{argName: argValue})
			} else {
				result = append(result, current)
			}
		} else {
			result = append(result, current)
		}
	}

	return result
}
