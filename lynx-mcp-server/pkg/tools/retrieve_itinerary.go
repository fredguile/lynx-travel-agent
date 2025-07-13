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
	TOOL_RETRIEVE_ITINERARY             string = "retrieve_itinerary"
	TOOL_RETRIEVE_ITINERARY_DESCRIPTION string = "Retrieve file itinerary"
	TOOL_RETRIEVE_ITINERARY_SCHEMA      string = `{
		"type": "object",
		"description": "Retrieve file itinerary",
		"properties": {
			"fileIdentifier": {
				"type": "string",
				"description": "File identifier"
			}
		},
		"required": ["fileIdentifier"],
		"outputSchema": {
			"type": "object",
			"properties": {
				"type": {
					"type": "string",
					"description": "Itinerary type"
				},
				"partyName": {
					"type": "string",
					"description": "Party name"
				},
				"fileReference": {
					"type": "string",
					"description": "File reference"
				},
				"fileIdentifier": {
					"type": "string",
					"description": "File identifier"
				},
				"clientIdentifier": {
					"type": "string",
					"description": "Client identifier"
				},
				"agentReference": {
					"type": "string",
					"description": "Agent reference"
				},
				"itineraryCount": {
					"type": "integer",
					"description": "Number of itineraries found"
				},
				"itineraries": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"voucherIdentifier": {
								"type": "string",
								"description": "Voucher identifier"
							},
							"date": {
								"type": "string",
								"description": "Date"
							},
							"transactionIdentifier": {
								"type": "string",
								"description": "Transaction identifier"
							},
							"supplier": {
								"type": "string",
								"description": "Supplier name"
							},
							"status": {
								"type": "string",
								"description": "Status"
							},
							"confirmationNumber": {
								"type": "string",
								"description": "Confirmation number"
							},
							"location": {
								"type": "string",
								"description": "Location"
							}
						},
						"required": ["voucherIdentifier", "date", "transactionIdentifier", "supplier", "status", "location"]
					}
				}
			},
			"required": ["type", "partyName", "fileReference", "fileIdentifier", "clientIdentifier", "agentReference", "itineraryCount", "itineraries"]
		}
	}`

	LYNX_RETRIEVE_ITINERARY_URL string = "/lynx/service/file.rpc"
)

func HandleRetrieveItinerary(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	session, _, err := utils.GetOrCreateSession(ctx, lynxConfig)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	arguments := request.GetArguments()

	fileIdentifier, ok := arguments["fileIdentifier"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid file identifier argument: %v", arguments["fileIdentifier"])
	}

	body := gwt.BuildRetrieveItineraryGWTBody(&gwt.RetrieveItineraryArgs{
		RemoteHost:     lynxConfig.RemoteHost,
		FileIdentifier: fileIdentifier,
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, LYNX_RETRIEVE_ITINERARY_URL), strings.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("failed to create retrieve itinerary request: %w", err)
	}

	req.Header.Set("Content-Type", gwt.CONTENT_TYPE)
	req.AddCookie(utils.CreateAuthCookie(lynxConfig, session))

	// Use retry utility with exponential backoff
	resp, bodyStr, err := utils.RetryHTTPRequest(ctx, client, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute retrieve itinerary request after retries: %w", err)
	}
	defer resp.Body.Close()

	// Parse the GWT response body
	retrieveItineraryResponseBody, err := gwt.ParseRetrieveItineraryResponseBody(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse retrieve itinerary response: %w", err)
	}

	return utils.NewToolResultJSON(retrieveItineraryResponseBody), nil
}

// GetRetrieveItinerarySchema returns the complete JSON schema for the retrieve itinerary tool
func GetRetrieveItinerarySchema() json.RawMessage {
	return json.RawMessage(TOOL_RETRIEVE_ITINERARY_SCHEMA)
}
