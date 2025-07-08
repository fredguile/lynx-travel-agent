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
	TOOL_FILE_DOCUMENT_UPLOAD                                 string = "file_document_upload"
	TOOL_FILE_DOCUMENT_UPLOAD_DESCRIPTION                     string = "Upload a file document (to use for transaction)"
	TOOL_FILE_DOCUMENT_UPLOAD_ARG_FILE_BINARY64               string = "fileBinary"
	TOOL_FILE_DOCUMENT_UPLOAD_ARG_FILE_BINARY64_DESCRIPTION   string = "Base64-encoded file binary"
	TOOL_FILE_DOCUMENT_UPLOAD_ARG_FILE_IDENTIFIER             string = "fileIdentifier"
	TOOL_FILE_DOCUMENT_UPLOAD_ARG_FILE_IDENTIFIER_DESCRIPTION string = "File identifier"
	TOOL_FILE_DOCUMENT_UPLOAD_ARG_FILE_NAME                   string = "fileName"
	TOOL_FILE_DOCUMENT_UPLOAD_ARG_FILE_NAME_DESCRIPTION       string = "File name"

	TOOL_FILE_DOCUMENT_UPLOAD_SCHEMA = `{
		"type": "object",
		"description": "Upload a file document (to use for transaction)",
		"properties": {
			"fileBinary": {
				"type": "string",
				"description": "Base64-encoded file binary"
			},
			"fileIdentifier": {
				"type": "string",
				"description": "File identifier"
			},
			"fileName": {
				"type": "string",
				"description": "File name"
			}
		},
		"required": ["fileBinary", "fileIdentifier", "fileName"],
		"outputSchema": {
			"type": "object",
			"properties": {
				"fileUrl": {
					"type": "string",
					"description": "File URL on lynx server"
				}
			},
			"required": ["fileUrl"]
		}
	}`

	LYNX_FILE_DOCUMENT_UPLOAD_URL string = "/lynx/fileDocumentUpload"
)

func HandleFileDocumentUpload(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	session, _, err := utils.GetOrCreateSession(ctx, lynxConfig)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	arguments := request.GetArguments()
	fileBinaryAsBase64, ok := arguments["fileBinary"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid file binary argument")
	}
	fileIdentifier, ok := arguments["fileIdentifier"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid file identifier argument")
	}
	fileName, ok := arguments["fileName"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid file name argument")
	}

	// Create a buffer to write multipart data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add fileId part
	fileIdField, err := writer.CreateFormField("fileId")
	if err != nil {
		return nil, fmt.Errorf("failed to create fileId field: %w", err)
	}
	fileIdField.Write([]byte(fileIdentifier))

	// Decode base64 data
	fileData, err := base64.StdEncoding.DecodeString(fileBinaryAsBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 file data: %w", err)
	}

	// Add file part
	fileField, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create file field: %w", err)
	}
	fileField.Write(fileData)

	// Close the multipart writer
	writer.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, LYNX_FILE_DOCUMENT_UPLOAD_URL), &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create file document upload request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(utils.CreateAuthCookie(lynxConfig, session))

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute file document upload request: %w", err)
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
		return nil, fmt.Errorf("file document upload failed with status %d: %s", resp.StatusCode, bodyStr)
	}

	fileUrl, err := parseResponseBody(bodyStr)

	if err != nil {
		return nil, fmt.Errorf("Invalid file upload response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(`{"fileUrl": "%s"}`, fileUrl),
			},
		},
	}, nil
}

// Parse response body to extract the file URL
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
		return "", fmt.Errorf("invalid file URL in response: %s", responseBody)
	}

	return urlPart, nil
}

// GetFileSearchByPartyNameSchema returns the complete JSON schema for the file search tool
func GetFileDocumentUploadSchema() json.RawMessage {
	return json.RawMessage(TOOL_FILE_DOCUMENT_UPLOAD_SCHEMA)
}
