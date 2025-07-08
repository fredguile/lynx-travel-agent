package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/gwt"
	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	TOOL_FILE_SEARCH_BY_FILE_REFERENCE             string = "file_search_by_file_reference"
	TOOL_FILE_SEARCH_BY_FILE_REFERENCE_DESCRIPTION string = "Retrieve file from file reference"
	TOOL_FILE_SEARCH_BY_FILE_REFERENCE_SCHEMA      string = `{
		"type": "object",
		"description": "Retrieve file from file reference",
		"properties": {
			"fileReference": {
				"type": "string",
				"description": "File reference"
			}
		},
		"required": ["fileReference"],
		"outputSchema": {
			"type": "object",
			"properties": {
				"count": {
					"type": "integer",
					"description": "Number of results found"
				},
				"results": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"companyCode": {
								"type": "string",
								"description": "Company code"
							},
							"clientReference": {
								"type": "string",
								"description": "Client reference"
							},
							"currency": {
								"type": "string",
								"description": "Currency code"
							},
							"fileIdentifier": {
								"type": "string",
								"description": "File identifier"
							},
							"fileReference": {
								"type": "string",
								"description": "File reference"
							},
							"partyName": {
								"type": "string",
								"description": "Party name"
							},
							"status": {
								"type": "string",
								"description": "File status"
							},
							"travelDate": {
								"type": "string",
								"description": "Travel date"
							}
						},
						"required": ["companyCode", "clientReference", "currency", "fileIdentifier", "fileReference", "partyName", "status", "travelDate"]
					}
				}
			},
			"required": ["count", "results"]
		}
	}`

	LYNX_FILE_SEARCH_BY_FILE_REFERENCE_URL string = "/lynx/service/file.rpc"
)

func HandleFileSearchByFileReference(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	session, _, err := utils.GetOrCreateSession(ctx, lynxConfig)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	arguments := request.GetArguments()

	fileReference, ok := arguments["fileReference"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid file reference argument: %v", arguments["fileReference"])
	}

	body := gwt.BuildFileSearchByFileReferenceGWTBody(&gwt.FileSearchByFileReferenceArgs{
		RemoteHost:    lynxConfig.RemoteHost,
		FileReference: fileReference,
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, LYNX_FILE_SEARCH_BY_FILE_REFERENCE_URL), strings.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("failed to create file search request: %w", err)
	}

	req.Header.Set("Content-Type", gwt.CONTENT_TYPE)
	req.AddCookie(utils.CreateAuthCookie(lynxConfig, session))

	// Use retry utility with exponential backoff
	resp, bodyStr, err := utils.RetryHTTPRequest(ctx, client, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute file search request after retries: %w", err)
	}
	defer resp.Body.Close()

	// Parse the GWT response body
	fileSearchResponseBody, err := gwt.ParseFileSearchResponseBody(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file search response: %w", err)
	}

	return utils.NewToolResultJSON(fileSearchResponseBody), nil
}

// GetFileSearchByFileReferenceSchema returns the complete JSON schema for the file search by file reference tool
func GetFileSearchByFileReferenceSchema() json.RawMessage {
	return json.RawMessage(TOOL_FILE_SEARCH_BY_FILE_REFERENCE_SCHEMA)
}
