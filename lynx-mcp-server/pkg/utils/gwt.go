package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	GWT_CONTENT_TYPE = "text/x-gwt-rpc; charset=utf-8"
)

type GWTLoginArgs struct {
	RemoteHost  string
	CompanyCode string
	Username    string
	Password    string
}

// BuildGWTLoginBody constructs the GWT-RPC login body with the given company code.
func BuildGWTLoginBody(args *GWTLoginArgs) string {
	return fmt.Sprintf("7|0|9|https://%s/lynx/lynx/|4775EB021C85EC0B04470837F40FC64A|com.lynxtraveltech.common.gui.client.rpc.SecurityService|login|java.lang.String/2004016611|Z|%s|%s|%s|1|2|3|4|4|5|5|5|6|7|8|9|0|", args.RemoteHost, args.CompanyCode, args.Username, args.Password)
}

type GWTFileSearchArgs struct {
	RemoteHost string
	PartyName  string
}

func BuildGWTFileSearchBody(args *GWTFileSearchArgs) string {
	return fmt.Sprintf("7|0|9|https://%s/lynx/lynx/|63A734E3E71C14883B20AFEC1238F6A7|com.lynxtraveltech.client.client.rpc.FileService|search|com.lynxtraveltech.client.shared.model.FileSearchCriteria/1867541444||%s|PARTY_NAME|DD MMM YYYY|1|2|3|4|1|5|5|6|6|1|0|1|7|6|50|8|6|0|9|0|0|6|", args.RemoteHost, args.PartyName)
}

// ParseResponseBody parses a GWT response body and extracts the data array.
// Returns the parsed data as a slice of interface{} containing the array elements.
func ParseResponseBody(responseBody string) ([]interface{}, error) {
	// Remove the "//OK" prefix if present
	body := strings.TrimPrefix(responseBody, "//OK")

	// Parse the main array structure
	parsedArray, err := parseGWTArray(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse main array: %w", err)
	}

	// Reverse the order of items
	for i, j := 0, len(parsedArray)-1; i < j; i, j = i+1, j-1 {
		parsedArray[i], parsedArray[j] = parsedArray[j], parsedArray[i]
	}

	// Check if we have enough items
	if len(parsedArray) < 3 {
		return nil, fmt.Errorf("response array too short, expected at least 3 items, got %d", len(parsedArray))
	}

	// The third item (index 2) contains the actual data array
	dataItem := parsedArray[2]

	// Check if the third item is an array
	dataArray, ok := dataItem.([]interface{})
	if !ok {
		return nil, fmt.Errorf("third item is not an array, got %T", dataItem)
	}

	return dataArray, nil
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
	return unescapeGWTErrorString(errorMessage), nil
}

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
	var depth int

	for i := 0; i < len(arrayStr); i++ {
		char := arrayStr[i]

		switch char {
		case '\'':
			if !inString {
				inString = true
			} else {
				// Check if it's an escaped quote
				if i+1 < len(arrayStr) && arrayStr[i+1] == '\'' {
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

	// Check if it's a quoted string
	if strings.HasPrefix(element, "'") && strings.HasSuffix(element, "'") {
		// Remove quotes and unescape
		content := element[1 : len(element)-1]
		content = strings.ReplaceAll(content, "''", "'")
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

// unescapeGWTErrorString removes surrounding double quotes and converts unicode escape sequences
func unescapeGWTErrorString(s string) string {
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
