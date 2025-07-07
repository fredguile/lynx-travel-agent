package gwt

import (
	"fmt"
	"strconv"
	"strings"
)

// parseGWTArray parses a GWT array format and returns the elements
func parseGWTArray(arrayStr string) ([]interface{}, error) {
	// Remove outer brackets
	arrayStr = strings.Trim(arrayStr, "[]")

	if arrayStr == "" {
		return []interface{}{}, nil
	}

	var result []interface{}
	var current strings.Builder
	var inString bool
	var quoteChar byte
	var depth int

	for i := 0; i < len(arrayStr); i++ {
		char := arrayStr[i]

		switch char {
		case '\'', '"':
			if !inString {
				inString = true
				quoteChar = char
			} else if char == quoteChar {
				// Check if it's an escaped quote
				if i+1 < len(arrayStr) && arrayStr[i+1] == char {
					current.WriteByte(char)
					i++ // Skip the next quote
				} else {
					inString = false
				}
			}
			current.WriteByte(char)
		case '[':
			if !inString {
				depth++
			}
			current.WriteByte(char)
		case ']':
			if !inString {
				depth--
			}
			current.WriteByte(char)
		case ',':
			if !inString && depth == 0 {
				// End of current element
				element := strings.TrimSpace(current.String())
				if element != "" {
					parsed, err := parseGWTElement(element)
					if err != nil {
						return nil, fmt.Errorf("failed to parse element '%s': %w", element, err)
					}
					result = append(result, parsed)
				}
				current.Reset()
			} else {
				current.WriteByte(char)
			}
		default:
			current.WriteByte(char)
		}
	}

	// Don't forget the last element
	element := strings.TrimSpace(current.String())
	if element != "" {
		parsed, err := parseGWTElement(element)
		if err != nil {
			return nil, fmt.Errorf("failed to parse last element '%s': %w", element, err)
		}
		result = append(result, parsed)
	}

	return result, nil
}

// parseGWTElement parses a single GWT element (string, number, or nested array)
func parseGWTElement(element string) (interface{}, error) {
	element = strings.TrimSpace(element)

	// Check if it's a quoted string (single or double quotes)
	if (strings.HasPrefix(element, "'") && strings.HasSuffix(element, "'")) ||
		(strings.HasPrefix(element, "\"") && strings.HasSuffix(element, "\"")) {
		// Remove quotes and unescape
		content := element[1 : len(element)-1]
		if strings.HasPrefix(element, "'") {
			content = strings.ReplaceAll(content, "''", "'")
		} else {
			content = strings.ReplaceAll(content, "\"\"", "\"")
		}
		return content, nil
	}

	// Check if it's a number
	if num, err := strconv.Atoi(element); err == nil {
		return num, nil
	}

	// Check if it's a float
	if num, err := strconv.ParseFloat(element, 64); err == nil {
		return num, nil
	}

	// Check if it's a nested array
	if strings.HasPrefix(element, "[") && strings.HasSuffix(element, "]") {
		return parseGWTArray(element)
	}

	// Return as string if nothing else matches
	return element, nil
}

// unescapeGWTString removes surrounding double quotes and converts unicode escape sequences
func unescapeGWTString(s string) string {
	// Remove surrounding double quotes if present
	s = strings.Trim(s, "\"")

	// Replace unicode escape sequences like \x27 with actual characters
	// This is a simple implementation - in a more robust version you might want to use regex
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if i+3 < len(s) && s[i] == '\\' && s[i+1] == 'x' {
			// Found \x sequence, try to parse the hex value
			hexStr := s[i+2 : i+4]
			if val, err := strconv.ParseUint(hexStr, 16, 8); err == nil {
				result.WriteByte(byte(val))
				i += 3 // Skip the \x and the two hex digits
			} else {
				// If parsing fails, keep the original sequence
				result.WriteByte(s[i])
			}
		} else {
			result.WriteByte(s[i])
		}
	}

	return result.String()
}

// convertHexEscapes converts hex escape sequences like \x26 to their actual characters
func convertHexEscapes(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if i+3 < len(s) && s[i] == '\\' && s[i+1] == 'x' {
			// Found \x sequence, try to parse the hex value
			hexStr := s[i+2 : i+4]
			if val, err := strconv.ParseUint(hexStr, 16, 8); err == nil {
				result.WriteByte(byte(val))
				i += 3 // Skip the \x and the two hex digits
			} else {
				// If parsing fails, keep the original sequence
				result.WriteByte(s[i])
			}
		} else {
			result.WriteByte(s[i])
		}
	}
	return result.String()
}

// ParseResponseError parses a GWT error response body and extracts the error message.
// Returns the parsed error message as a string.
func ParseResponseError(responseBody string) (string, error) {
	// Remove the "//EX" prefix if present
	body := strings.TrimPrefix(responseBody, "//EX")

	// Parse the main array structure
	parsedArray, err := parseGWTArray(body)
	if err != nil {
		return "", fmt.Errorf("failed to parse error array: %w", err)
	}

	// Check if we have enough items
	if len(parsedArray) < 3 {
		return "", fmt.Errorf("error response array too short, expected at least 3 items, got %d", len(parsedArray))
	}

	// The third item (index 2) contains the error details array
	errorItem := parsedArray[2]

	// Check if the third item is an array
	errorArray, ok := errorItem.([]interface{})
	if !ok {
		return "", fmt.Errorf("third item is not an array, got %T", errorItem)
	}

	// Check if we have enough items in the error array
	if len(errorArray) < 2 {
		return "", fmt.Errorf("error details array too short, expected at least 2 items, got %d", len(errorArray))
	}

	// The error message might be split across multiple elements due to commas in the message
	// Start from index 1 (after the exception class name) and concatenate all string elements
	var errorMessageParts []string
	for i := 1; i < len(errorArray); i++ {
		if msgPart, ok := errorArray[i].(string); ok {
			errorMessageParts = append(errorMessageParts, msgPart)
		}
	}

	if len(errorMessageParts) == 0 {
		return "", fmt.Errorf("no error message found in error array")
	}

	// Join the message parts and unescape
	errorMessage := strings.Join(errorMessageParts, ", ")
	return unescapeGWTString(errorMessage), nil
}
