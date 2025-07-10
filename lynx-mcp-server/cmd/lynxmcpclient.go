//go:build client

package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"dodmcdund.cc/lynx-travel-agent/lynxmcpserver/pkg/config"
	"dodmcdund.cc/lynx-travel-agent/lynxmcpserver/pkg/utils"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

var clientConfig = config.NewClientConfig()

// readDemoFileAsBase64 reads the demo PDF file and returns its base64 encoded content
func readDemoFileAsBase64() (string, error) {
	filePath := "assets/dummy.pdf"

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open demo file %s: %w", filePath, err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read demo file %s: %w", filePath, err)
	}

	encodedData := base64.StdEncoding.EncodeToString(fileData)
	return encodedData, nil
}

func main() {
	// Set custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [flags] --command \"tool_name [arguments]\"\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nCommand arguments:\n")
		fmt.Fprintf(os.Stderr, "  --binary                    Include demo PDF file as base64-encoded binary argument\n")
		fmt.Fprintf(os.Stderr, "  --argName=value             Named arguments for the tool (inline value)\n")
		fmt.Fprintf(os.Stderr, "  --argName \"value with spaces\" Named arguments for the tool (separate value)\n")
		fmt.Fprintf(os.Stderr, "  \"quoted value\"              Positional arguments for the tool\n")
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --url http://localhost:9600/sse --command \"some_tool --binary --param1=value1 --param2 \\\"value with spaces\\\"\"\n", os.Args[0])
	}

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
		hasBinary := false

		for i, arg := range toolArgs {
			switch v := arg.(type) {
			case string:
				// Check if this is the --binary flag
				if v == "--binary" {
					hasBinary = true
					continue
				}
				args[fmt.Sprintf("arg%d", i)] = v
			case map[string]string:
				// Convert map[string]string to map[string]interface{}
				for k, val := range v {
					args[k] = val
				}
			}
		}

		// If --binary was specified, add the encoded file content
		if hasBinary {
			binaryBase64, err := readDemoFileAsBase64()
			if err != nil {
				log.Printf("Failed to read demo file: %v", err)
				os.Exit(1)
			}
			args["binary"] = binaryBase64
			fmt.Println("Added binary argument using demo PDF content")
		}

		// Create a copy of args for logging with obfuscated binary
		logArgs := make(map[string]interface{})
		for k, v := range args {
			if k == "binary" {
				logArgs[k] = "<BINARY>"
			} else {
				logArgs[k] = v
			}
		}
		fmt.Printf("Calling tool named \"%s\" with args: %+v\n", toolName, logArgs)

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

		fmt.Println("Tool execution result:")
		fmt.Println(utils.FormatJSONResult(result))
	}

	fmt.Println("Client initialized successfully. Shutting down...")
	c.Close()
}

func parseCommand(cmd string) []interface{} {
	var result []interface{}
	var current string
	var inQuote bool
	var quoteChar rune
	var expectingValue bool
	var currentArgName string

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
						// This is a flag without a value, expect the next token to be the value
						currentArgName = current[2:] // Remove the -- prefix
						expectingValue = true
						result = append(result, current)
					}
				} else if expectingValue {
					// This is the value for the previous argument
					result = append(result, map[string]string{currentArgName: current})
					expectingValue = false
					currentArgName = ""
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
				// This is a flag without a value at the end
				result = append(result, current)
			}
		} else if expectingValue {
			// This is the value for the previous argument
			result = append(result, map[string]string{currentArgName: current})
		} else {
			result = append(result, current)
		}
	}

	return result
}
