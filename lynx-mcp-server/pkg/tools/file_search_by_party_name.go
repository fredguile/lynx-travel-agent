package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"dodmcdund.cc/lynx-travel-agent/lynxmcpserver/pkg/gwt"
	"dodmcdund.cc/lynx-travel-agent/lynxmcpserver/pkg/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	TOOL_FILE_SEARCH_BY_PARTY_NAME             string = "file_search_by_party_name"
	TOOL_FILE_SEARCH_BY_PARTY_NAME_DESCRIPTION string = "Retrieve file from party name"
	TOOL_FILE_SEARCH_BY_PARTY_NAME_SCHEMA      string = `{
		"type": "object",
		"description": "Retrieve file from party name",
		"properties": {
			"partyName": {
				"type": "string",
				"description": "Party name"
			}
		},
		"required": ["partyName"],
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
							"clientIdentifier": {
								"type": "string",
								"description": "Client identifier"
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
						"required": ["companyCode", "clientIdentifier", "clientReference", "currency", "fileIdentifier", "fileReference", "partyName", "status", "travelDate"]
					}
				}
			},
			"required": ["count", "results"]
		}
	}`

	LYNX_FILE_SEARCH_BY_PARTY_NAME_URL string = "/lynx/service/file.rpc"
)

func HandleFileSearchByPartyName(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	session, _, err := utils.GetOrCreateSession(ctx, lynxConfig)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	arguments := request.GetArguments()

	partyName, ok := arguments["partyName"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid partyName argument: %v", arguments["partyName"])
	}

	body := gwt.BuildFileSearchByPartyNameGWTBody(&gwt.FileSearchByPartyNameArgs{
		RemoteHost: lynxConfig.RemoteHost,
		PartyName:  partyName,
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, LYNX_FILE_SEARCH_BY_PARTY_NAME_URL), strings.NewReader(body))

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
		return nil, fmt.Errorf("failed to parse  file search response: %w", err)
	}

	return utils.NewToolResultJSON(fileSearchResponseBody), nil
}

// GetFileSearchByPartyNameSchema returns the complete JSON schema for the file search tool
func GetFileSearchByPartyNameSchema() json.RawMessage {
	return json.RawMessage(TOOL_FILE_SEARCH_BY_PARTY_NAME_SCHEMA)
}
