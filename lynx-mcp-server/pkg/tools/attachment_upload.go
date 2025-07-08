package tools

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	ATTACHMENT_UPLOAD             string = "attachment_upload"
	ATTACHMENT_UPLOAD_DESCRIPTION string = "Upload attachment for using with file document"
	ATTACHMENT_UPLOAD_SCHEMA      string = `{
		"type": "object",
		"description": "Upload an attachment for using with file document",
		"properties": {
			"binary": {
				"type": "string",
				"description": "Base64-encoded binary"
			},
			"identifer": {
				"type": "string",
				"description": "Unique identifier"
			},
			"fileName": {
				"type": "string",
				"description": "File name"
			}
		},
		"required": ["binary", "identifer", "fileName"],
		"outputSchema": {
			"type": "object",
			"properties": {
				"attachmentUrl": {
					"type": "string",
					"description": "Attachment URL on server"
				}
			},
			"required": ["attachmentUrl"]
		}
	}`

	LYNX_ATTACHMENT_UPLOAD_URL string = "/lynx/fileDocumentUpload"
)

func HandleAttachmentUpload(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	session, _, err := utils.GetOrCreateSession(ctx, lynxConfig)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	arguments := request.GetArguments()

	binaryAsBase64, ok := arguments["binary"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid binary argument: %v", arguments["binary"])
	}

	identifer, ok := arguments["identifer"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid identifier argument: %v", arguments["identifer"])
	}

	fileName, ok := arguments["fileName"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid file name argument: %v", arguments["fileName"])
	}

	// Create a buffer to write multipart data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add fileId part
	fileIdField, err := writer.CreateFormField("fileId")
	if err != nil {
		return nil, fmt.Errorf("failed to create fileId field: %w", err)
	}
	fileIdField.Write([]byte(identifer))

	// Decode base64 data
	fileData, err := base64.StdEncoding.DecodeString(binaryAsBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Add file part
	fileField, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create file field: %w", err)
	}
	fileField.Write(fileData)

	// Close the multipart writer
	writer.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, LYNX_ATTACHMENT_UPLOAD_URL), &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create attachment upload request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(utils.CreateAuthCookie(lynxConfig, session))

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute attachment upload request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	bodyStr := string(bodyBytes)

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Attachment upload failed with status %d: %s", resp.StatusCode, bodyStr)
	}

	attachmentUrl, err := parseResponseBody(bodyStr)

	if err != nil {
		return nil, fmt.Errorf("Invalid attachment upload response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(`{"attachmentUrl": "%s"}`, attachmentUrl),
			},
		},
	}, nil
}

// Parse response body to extract the attachment URL
// Expected format: "SUCCESS:/documents/file/f16476987/d20250708231038.pdf:\n"
func parseResponseBody(responseBody string) (string, error) {
	if !strings.HasPrefix(responseBody, "SUCCESS:") {
		return "", fmt.Errorf("unexpected response format: %s", responseBody)
	}

	// Remove "SUCCESS:" prefix
	urlPart := strings.TrimPrefix(responseBody, "SUCCESS:")

	// Trim whitespace and line breaks
	urlPart = strings.TrimSpace(urlPart)

	// Ensure response ends with ": " and strip it out
	if !strings.HasSuffix(urlPart, ":") {
		return "", fmt.Errorf("response does not end with ':': %s", responseBody)
	}
	urlPart = strings.TrimSuffix(urlPart, ":")

	// Validate that we have a URL path
	if urlPart == "" || !strings.HasPrefix(urlPart, "/") {
		return "", fmt.Errorf("invalid attachment URL in response: %s", responseBody)
	}

	return urlPart, nil
}

// GetAttachmentUploadSchema returns the complete JSON schema for the file search tool
func GetAttachmentUploadSchema() json.RawMessage {
	return json.RawMessage(ATTACHMENT_UPLOAD_SCHEMA)
}
