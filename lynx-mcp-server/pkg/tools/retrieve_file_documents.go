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
	TOOL_RETRIEVE_FILE_DOCUMENTS             string = "retrieve_file_documents"
	TOOL_RETRIEVE_FILE_DOCUMENTS_DESCRIPTION string = "Retrieve file documents from transaction reference"
	TOOL_RETRIEVE_FILE_DOCUMENTS_SCHEMA      string = `{
		"type": "object",
		"description": "Retrieve file documents from transaction reference",
		"properties": {
			"fileIdentifier": {
				"type": "string",
				"description": "File identifier"
			},
			"transactionIdentifier": {
				"type": "string",
				"description": "Transaction identifier"
			}
		},
		"required": ["fileIdentifier", "transactionIdentifier"],
		"outputSchema": {
			"type": "object",
			"properties": {
				"count": {
					"type": "integer",
					"description": "Number of documents found"
				},
				"results": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"fileIdentifier": {
								"type": "string",
								"description": "File identifier"
							},
							"transactionIdentifier": {
								"type": "string",
								"description": "Transaction identifier"
							},
							"documentIdentifier": {
								"type": "string",
								"description": "Document identifier"
							},
							"documentName": {
								"type": "string",
								"description": "Document name"
							},
							"documentType": {
								"type": "string",
								"description": "Document Type"
							},
							"content": {
								"type": "string",
								"description": "Content"
							},
							"attachmentUrl": {
								"type": "string",
								"description": "Attached file"
							}
						},
						"required": ["fileIdentifier", "transactionIdentifier", "documentIdentifier", "documentName", "documentType", "content", "attachedFile"]
					}
				}
			},
			"required": ["count", "results"]
		}
	}`

	LYNX_FILE_DOCUMENTS_BY_TRANSACTION_REFERENCE_URL string = "/lynx/service/file.rpc"
)

func HandleRetrieveFileDocuments(
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

	transactionIdentifier, ok := arguments["transactionIdentifier"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid transaction identifier argument: %v", arguments["transactionIdentifier"])
	}

	body := gwt.BuildFileDocumentsByTransactionReferenceGWTBody(&gwt.FileDocumentsByTransactionReferenceArgs{
		RemoteHost:            lynxConfig.RemoteHost,
		FileIdentifier:        fileIdentifier,
		TransactionIdentifier: transactionIdentifier,
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, LYNX_FILE_DOCUMENTS_BY_TRANSACTION_REFERENCE_URL), strings.NewReader(body))

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
	fielDocumentsListResponseBody, err := gwt.ParseFileDocumentsListResponseBody(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file documents response: %w", err)
	}

	return utils.NewToolResultJSON(fielDocumentsListResponseBody), nil
}

// GetRetrieveFileDocumentsSchema returns the complete JSON schema for the retrieve file documents tool
func GetRetrieveFileDocumentsSchema() json.RawMessage {
	return json.RawMessage(TOOL_RETRIEVE_FILE_DOCUMENTS_SCHEMA)
}
