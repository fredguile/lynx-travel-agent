package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// formatJSONResult formats a result object as pretty JSON for logging
func FormatJSONResult(result interface{}) string {
	// First, marshal the result to JSON
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}

	// Parse the JSON back to a map to process nested JSON strings
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return string(jsonBytes) // Return original if parsing fails
	}

	// Recursively process nested JSON strings
	processNestedJSON(data)

	// Marshal back to pretty JSON
	finalJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return string(jsonBytes) // Return original if final marshaling fails
	}

	return string(finalJSON)
}

// processNestedJSON recursively processes a map to prettify nested JSON strings
func processNestedJSON(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			switch val := value.(type) {
			case string:
				// Check if this string looks like JSON
				if isJSONString(val) {
					// Try to parse the JSON string and replace it with the parsed object
					if parsed := parseJSONString(val); parsed != nil {
						v[key] = parsed
					}
				}
			case map[string]interface{}:
				processNestedJSON(val)
			case []interface{}:
				processNestedJSON(val)
			}
		}
	case []interface{}:
		for i, item := range v {
			switch val := item.(type) {
			case string:
				// Check if this string looks like JSON
				if isJSONString(val) {
					// Try to parse the JSON string and replace it with the parsed object
					if parsed := parseJSONString(val); parsed != nil {
						v[i] = parsed
					}
				}
			case map[string]interface{}:
				processNestedJSON(val)
			case []interface{}:
				processNestedJSON(val)
			}
		}
	}
}

// isJSONString checks if a string looks like JSON
func isJSONString(s string) bool {
	s = strings.TrimSpace(s)
	return (strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) ||
		(strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"))
}

// parseJSONString attempts to parse a JSON string and returns the parsed object
func parseJSONString(s string) interface{} {
	var data interface{}
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return nil // Return nil if parsing fails
	}
	return data
}
