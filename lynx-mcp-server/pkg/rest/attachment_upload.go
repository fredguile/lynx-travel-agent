package rest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"dodmcdund.cc/lynx-travel-agent/lynxmcpserver/pkg/config"
	"dodmcdund.cc/lynx-travel-agent/lynxmcpserver/pkg/tools"
	"dodmcdund.cc/lynx-travel-agent/lynxmcpserver/pkg/utils"
)

const (
	MAX_UPLOAD_SIZE = 32 << 20 // 32MB
)

// HandleAttachmentUpload handles the REST endpoint for attachment upload
func HandleAttachmentUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get fileId from form
	fileId := r.FormValue("fileId")
	if fileId == "" {
		http.Error(w, "fileId is required", http.StatusBadRequest)
		return
	}

	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode file data to base64
	binaryAsBase64 := base64.StdEncoding.EncodeToString(fileData)

	// Create a buffer to write multipart data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add fileId part
	fileIdField, err := writer.CreateFormField("fileId")
	if err != nil {
		http.Error(w, "Failed to create fileId field: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fileIdField.Write([]byte(fileId))

	// Decode base64 data
	decodedFileData, err := base64.StdEncoding.DecodeString(binaryAsBase64)
	if err != nil {
		http.Error(w, "Failed to decode base64 data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add file part
	fileField, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		http.Error(w, "Failed to create file field: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fileField.Write(decodedFileData)

	// Close the multipart writer
	writer.Close()

	// Get session
	lynxConfig := config.NewLynxServerConfig()
	session, _, err := utils.GetOrCreateSession(r.Context(), lynxConfig)
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, tools.LYNX_ATTACHMENT_UPLOAD_URL), &requestBody)
	if err != nil {
		http.Error(w, "Failed to create attachment upload request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(utils.CreateAuthCookie(lynxConfig, session))

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to execute attachment upload request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	bodyStr := string(bodyBytes)

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Attachment upload failed with status %d: %s", resp.StatusCode, bodyStr), http.StatusInternalServerError)
		return
	}

	// Parse response to extract attachment URL
	attachmentUrl, err := parseResponseBody(bodyStr)
	if err != nil {
		http.Error(w, "Invalid attachment upload response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"attachmentUrl": "%s"}`, attachmentUrl)))
}

// parseResponseBody parses the response body to extract the attachment URL
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
