package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"dodmcdund.cc/lynx-travel-agent/lynxmcpserver/pkg/gwt"
)

// RetryConfig holds configuration for retry behavior
type RetryConfig struct {
	MaxAttempts       int
	InitialDelay      time.Duration
	BackoffMultiplier float64
	MaxDelay          time.Duration
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:       5,
		InitialDelay:      0, // Immediate first retry
		BackoffMultiplier: 2.0,
		MaxDelay:          30 * time.Second,
	}
}

// RetryHTTPRequest executes an HTTP request with exponential backoff retry logic
// Success condition: response status is OK, has body, and body starts with "//OK"
func RetryHTTPRequest(ctx context.Context, client *http.Client, req *http.Request, config *RetryConfig) (*http.Response, string, error) {
	if config == nil {
		config = DefaultRetryConfig()
	}

	// Read the original request body once to preserve it for retries
	var originalBody []byte
	var err error
	if req.Body != nil {
		originalBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read request body: %w", err)
		}
		req.Body.Close()
	}

	var lastErr error
	var lastResp *http.Response
	var lastBodyStr string

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// Create a new request for each attempt to ensure the body is preserved
		var retryReq *http.Request
		if len(originalBody) > 0 {
			retryReq, err = http.NewRequest(req.Method, req.URL.String(), strings.NewReader(string(originalBody)))
		} else {
			retryReq, err = http.NewRequest(req.Method, req.URL.String(), nil)
		}
		if err != nil {
			return nil, "", fmt.Errorf("failed to create retry request: %w", err)
		}

		// Copy headers from original request
		for key, values := range req.Header {
			for _, value := range values {
				retryReq.Header.Add(key, value)
			}
		}

		// Execute the request
		resp, err := client.Do(retryReq)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to execute request: %w", attempt+1, err)

			if shouldReturn(attempt, config.MaxAttempts) {
				return nil, "", lastErr
			}

			if err := waitForRetry(ctx, attempt, config); err != nil {
				return nil, "", err
			}
			continue
		}

		// Read response body
		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Always close the body

		if err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to read response body: %w", attempt+1, err)
			lastResp = resp

			if shouldReturn(attempt, config.MaxAttempts) {
				return resp, "", lastErr
			}

			if err := waitForRetry(ctx, attempt, config); err != nil {
				return resp, "", err
			}
			continue
		}

		bodyStr := string(bodyBytes)

		// Check success conditions: status OK, has body, and body starts with "//OK"
		if resp.StatusCode == http.StatusOK && len(bodyStr) > 0 && strings.HasPrefix(bodyStr, "//OK") {
			// Success! Return the response and body
			log.Printf("RetryHTTPRequest Success! (attempt=%d)\n", attempt)
			return resp, bodyStr, nil
		}

		// Check for GWT error responses that start with "//EX"
		if resp.StatusCode == http.StatusOK && len(bodyStr) > 0 && strings.HasPrefix(bodyStr, "//EX") {
			// Parse the GWT error response to extract the error message
			errorMessage, err := gwt.ParseResponseError(bodyStr)
			if err != nil {
				// If we can't parse the error, return the raw body as error
				lastErr = fmt.Errorf("RetryHTTPRequest attempt %d: GWT error response (unparseable): %s", attempt+1, bodyStr)
			} else {
				// Return the parsed error message
				lastErr = fmt.Errorf("RetryHTTPRequest attempt %d: GWT error: %s", attempt+1, errorMessage)
			}
			lastResp = resp
			lastBodyStr = bodyStr

			// GWT errors are not retryable, so return immediately
			return resp, "", lastErr
		}

		// If we reach here, the request was successful but didn't meet our success criteria
		lastErr = fmt.Errorf("RetryHTTPRequest attempt %d: unexpected response (status: %d, body: %s)", attempt+1, resp.StatusCode, bodyStr)
		lastResp = resp
		lastBodyStr = bodyStr

		if shouldReturn(attempt, config.MaxAttempts) {
			return resp, bodyStr, lastErr
		}

		if err := waitForRetry(ctx, attempt, config); err != nil {
			return resp, bodyStr, err
		}
	}

	// This should never be reached, but just in case
	return lastResp, lastBodyStr, lastErr
}

// shouldReturn checks if this is the last attempt and we should return instead of retry
func shouldReturn(attempt, maxAttempts int) bool {
	return attempt == maxAttempts-1
}

// waitForRetry waits for the calculated delay before the next retry attempt
func waitForRetry(ctx context.Context, attempt int, config *RetryConfig) error {
	delay := calculateDelay(attempt, config)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		return nil
	}
}

// calculateDelay calculates the delay for the next retry attempt
func calculateDelay(attempt int, config *RetryConfig) time.Duration {
	if attempt == 0 {
		return config.InitialDelay
	}

	// Calculate exponential backoff: 0s, 5s, 10s, 30s, 30s
	delays := []time.Duration{
		0 * time.Second,  // Immediate first retry
		5 * time.Second,  // 5 seconds
		10 * time.Second, // 10 seconds
		30 * time.Second, // 30 seconds
		30 * time.Second, // Max delay (capped)
	}

	if attempt < len(delays) {
		return delays[attempt]
	}

	return config.MaxDelay
}
