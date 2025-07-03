package utils

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// newToolResultJSON creates a new CallToolResult with JSON content
// This is a helper function to properly return JSON data instead of plain text
func NewToolResultJSON(data interface{}) *mcp.CallToolResult {
	jsonData, err := json.Marshal(data)
	if err != nil {
		// If marshaling fails, return an error result
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response to JSON: %v", err))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(jsonData),
			},
		},
	}
}
