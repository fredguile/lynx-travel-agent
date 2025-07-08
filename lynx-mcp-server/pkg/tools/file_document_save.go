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
	TOOL_FILE_DOCUMENT_SAVE             string = "file_document_save"
	TOOL_FILE_DOCUMENT_SAVE_DESCRIPTION string = "Save file document details"
	TOOL_FILE_DOCUMENT_SAVE_SCHEMA      string = `{
		"type": "object",
		"description": "Save file document details",
		"properties": {
			"fileIdentifier": {
				"type": "string",
				"description": "File identifier"
			},
			"transactionIdentifier": {
				"type": "string",
				"description": "Transaction identifier"
			},
			"name": {
				"type": "string",
				"description": "Document name"
			},
			"content": {
				"type": "string",
				"description": "Document content (as plain text or HTML)"
			},
			"type": {
				"type": "string",
				"description": "Document type"
			},
			"attachmentUrl": {
				"type": "string",
				"description": "Attachment URL"
			}
		},
		"required": ["fileIdentifier", "transactionIdentifier", "name", "content", "type"],
		"outputSchema": {
			"type": "object",
			"properties": {}
		}
	}`

	LYNX_FILE_DOCUMENT_SAVE_DETAILS_URL string = "/lynx/service/file.rpc"
)

func HandleFileDocumentSave(
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

	name, ok := arguments["name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid name argument: %v", arguments["name"])
	}

	content, ok := arguments["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content argument: %v", arguments["content"])
	}

	documentType, ok := arguments["type"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type argument: %v", arguments["type"])
	}

	attachmentUrl, ok := arguments["attachmentUrl"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid attachmentUrl argument: %v", arguments["attachmentUrl"])
	}

	body := gwt.BuildFileDocumentSaveGWTBody(&gwt.FileDocumentSaveDetailsArgs{
		RemoteHost:            lynxConfig.RemoteHost,
		FileIdentifier:        fileIdentifier,
		TransactionIdentifier: transactionIdentifier,

		Name:          name,
		Content:       content,
		Type:          documentType,
		AttachmentURL: attachmentUrl,
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, LYNX_FILE_DOCUMENT_SAVE_DETAILS_URL), strings.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("failed to create document save details request: %w", err)
	}

	req.Header.Set("Content-Type", gwt.CONTENT_TYPE)
	req.AddCookie(utils.CreateAuthCookie(lynxConfig, session))

	// Use retry utility with exponential backoff
	resp, bodyStr, err := utils.RetryHTTPRequest(ctx, client, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute document save details request after retries: %w", err)
	}
	defer resp.Body.Close()

	// Parse the GWT response body
	err = gwt.ParseFileDocumentSaveResponseBody(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file document save details response: %w", err)
	}

	return utils.NewToolResultJSON(map[string]interface{}{}), nil
}

func GetFileDocumentSaveDetailsSchema() json.RawMessage {
	return json.RawMessage(TOOL_FILE_DOCUMENT_SAVE_SCHEMA)
}
