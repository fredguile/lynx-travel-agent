package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RequestToCurl converts an http.Request to a curl command string for debugging
func RequestToCurl(req *http.Request) (string, error) {
	var curlCmd strings.Builder

	// Start with curl command
	curlCmd.WriteString("curl")

	// Add method if not GET
	if req.Method != "GET" {
		curlCmd.WriteString(fmt.Sprintf(" -X %s", req.Method))
	}

	// Add URL
	curlCmd.WriteString(fmt.Sprintf(" '%s'", req.URL.String()))

	// Add headers
	for name, values := range req.Header {
		for _, value := range values {
			curlCmd.WriteString(fmt.Sprintf(" -H '%s: %s'", name, value))
		}
	}

	// Add cookies
	if req.Header.Get("Cookie") != "" {
		curlCmd.WriteString(fmt.Sprintf(" -H 'Cookie: %s'", req.Header.Get("Cookie")))
	}

	// Add body if present
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read request body: %w", err)
		}

		// Restore the body for the original request
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) > 0 {
			// Escape single quotes in the body
			bodyStr := strings.ReplaceAll(string(bodyBytes), "'", "'\"'\"'")
			curlCmd.WriteString(fmt.Sprintf(" -d '%s'", bodyStr))
		}
	}

	// Add verbose flag for debugging
	curlCmd.WriteString(" -v")

	return curlCmd.String(), nil
}

// DebugRequest logs the request details and returns a curl command
func DebugRequest(req *http.Request) (string, error) {
	// Log request body if present
	var bodyStr string
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read request body: %w", err)
		}
		// Restore the body for the original request
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyStr = string(bodyBytes)
	}

	curlCmd, err := RequestToCurl(req)
	if err != nil {
		return "", err
	}

	fmt.Printf("=== HTTP Request Debug ===\n")
	fmt.Printf("Method: %s\n", req.Method)
	fmt.Printf("URL: %s\n", req.URL.String())
	fmt.Printf("Headers:\n")
	for name, values := range req.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}
	fmt.Printf("Cookies:\n")
	for _, cookie := range req.Cookies() {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
	if bodyStr != "" {
		fmt.Printf("Body:\n%s\n", bodyStr)
	}
	fmt.Printf("Curl command:\n%s\n", curlCmd)
	fmt.Printf("========================\n")

	return curlCmd, nil
}
