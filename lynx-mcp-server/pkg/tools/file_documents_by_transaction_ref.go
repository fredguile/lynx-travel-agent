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
	TOOL_FILE_DOCUMENTS_BY_TRANSACTION_REFERENCE                                       string = "file_documents_by_transaction_reference"
	TOOL_FILE_DOCUMENTS_BY_TRANSACTION_REFERENCE_DESCRIPTION                           string = "Retrieve file documents from transaction reference"
	TOOL_FILE_DOCUMENTS_BY_TRANSACTION_REFERENCE_ARG_TRANSACTION_REFERENCE             string = "transactionReference"
	TOOL_FILE_DOCUMENTS_BY_TRANSACTION_REFERENCE_ARG_TRANSACTION_REFERENCE_DESCRIPTION string = "Transaction reference"

	TOOL_FILE__DOCUMENTS_BY_TRANSACTION_REFERENCE_SCHEMA = `{
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
							"attachedFile": {
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

func HandleFileDocumentsByTransactionReference(
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
		return nil, fmt.Errorf("invalid file identifier argument")
	}

	transactionIdentifier, ok := arguments["transactionIdentifier"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid transaction identifier argument")
	}

	body := gwt.BuildFileDocumentsByTransactionReferenceGWTBody(&gwt.FileDocumentsByTransactionReferenceArgs{
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

func GetFileDocumentsByTransactioReferenceSchema() json.RawMessage {
	return json.RawMessage(TOOL_FILE__DOCUMENTS_BY_TRANSACTION_REFERENCE_SCHEMA)
}
