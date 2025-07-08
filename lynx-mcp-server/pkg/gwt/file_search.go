package gwt

import (
	"fmt"
	"strings"
)

type FileSearchByPartyNameArgs struct {
	RemoteHost string
	PartyName  string
}

// BuildFileSearchByPartyNameGWTBody constructs the GWT-RPC file search body with the given party name, including quotations
func BuildFileSearchByPartyNameGWTBody(args *FileSearchByPartyNameArgs) string {
	return fmt.Sprintf("7|0|9|https://%s/lynx/lynx/|63A734E3E71C14883B20AFEC1238F6A7|com.lynxtraveltech.client.client.rpc.FileService|search|com.lynxtraveltech.client.shared.model.FileSearchCriteria/1867541444||%s|PARTY_NAME|DD MMM YYYY|1|2|3|4|1|5|5|6|6|1|1|1|7|6|50|8|6|0|9|0|0|6|", args.RemoteHost, args.PartyName)
}

type FileSearchByFileReferenceArgs struct {
	RemoteHost    string
	FileReference string
}

// BuildFileSearchByFileReferenceGWTBody constructs the GWT-RPC file search body with the given file reference, including quotations
func BuildFileSearchByFileReferenceGWTBody(args *FileSearchByFileReferenceArgs) string {
	return fmt.Sprintf("7|0|9|https://%s/lynx/lynx/|63A734E3E71C14883B20AFEC1238F6A7|com.lynxtraveltech.client.client.rpc.FileService|search|com.lynxtraveltech.client.shared.model.FileSearchCriteria/1867541444||%s|PARTY_NAME|DD MMM YYYY|1|2|3|4|1|5|5|6|7|1|1|1|6|6|50|8|6|0|9|0|0|6|", args.RemoteHost, args.FileReference)
}

type FileSearchResponseArray struct {
	Count   int                `json:"count"`
	Results []FileSearchResult `json:"results"`
}

type FileSearchResult struct {
	CompanyCode      string `json:"companyCode"`
	ClientIdentifier string `json:"clientIdentifier"`
	ClientReference  string `json:"clientReference"`
	Currency         string `json:"currency"`
	FileIdentifier   string `json:"fileIdentifier"`
	FileReference    string `json:"fileReference"`
	PartyName        string `json:"partyName"`
	Status           string `json:"status"`
	TravelDate       string `json:"traveDate"`
}

// ParseFileSearchResponseBody parses a GWT response in context of file search.
// Returns the parsed data as FileSearchResponseArray containing the array elements.
func ParseFileSearchResponseBody(responseBody string) (*FileSearchResponseArray, error) {
	// Ensure response starts with "//OK"
	if !strings.HasPrefix(responseBody, "//OK") {
		return nil, fmt.Errorf("response body missing //OK")
	}

	// Remove the "//OK" prefix if present
	body := strings.TrimPrefix(responseBody, "//OK")

	// Parse the main array structure
	parsedArray, err := parseGWTArray(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse main array: %w", err)
	}

	// check that parsedArray has at least 4 items
	if len(parsedArray) < 4 {
		return nil, fmt.Errorf("response body should contain at least 4 items, got %d", len(parsedArray))
	}

	// check that last element of parsedArray contains protocol version 7
	if parsedArray[len(parsedArray)-1] != 7 {
		return nil, fmt.Errorf("response body should contain protocol version 7, got %d", parsedArray[len(parsedArray)-1])
	}

	// check that last last last element of parsedArray is an array, extract it as dataArray
	dataArray, ok := parsedArray[len(parsedArray)-3].([]interface{})
	if !ok {
		return nil, fmt.Errorf("response body should contain a data array, got %T", parsedArray[len(parsedArray)-3])
	}

	// check that last last last element of parsedArray is our starting point, a GWT_TYPE_ARRAY
	oneBasedIndex, ok := parsedArray[len(parsedArray)-4].(int)
	if !ok {
		return nil, fmt.Errorf("expected one-based index, got %T", parsedArray[len(parsedArray)-4])
	}

	mappedFirstStringValue, ok := dataArray[oneBasedIndex-1].(string)
	if !ok {
		return nil, fmt.Errorf("expected string value, got %T", dataArray[oneBasedIndex-1])
	}

	if !strings.HasPrefix(mappedFirstStringValue, GWT_TYPE_ARRAY) {
		return nil, fmt.Errorf("first item should be an array, got %T", mappedFirstStringValue)
	}

	// check that our GWT_TYPE_ARRAY has a size
	arraySize, ok := parsedArray[len(parsedArray)-5].(int)
	if !ok {
		return nil, fmt.Errorf("first array item should have a size, got %T", parsedArray[len(parsedArray)-5])
	}

	fileSearchResponse := FileSearchResponseArray{
		Count:   arraySize,
		Results: make([]FileSearchResult, arraySize),
	}

	for i, resultIndex := len(parsedArray)-6, 0; i >= 0; i-- {
		if oneBasedIndex, ok := parsedArray[i].(int); ok {
			// Ignore 0 aka nil values
			if oneBasedIndex == 0 {
				continue
			}

			currentValue := dataArray[oneBasedIndex-1]

			if currentStringValue, ok := currentValue.(string); ok && strings.HasPrefix(currentStringValue, GWT_TYPE_FILE_SEARCH_RESULTS) {
				fileSearchResult := FileSearchResult{
					CompanyCode:      dataArray[parsedArray[i-1].(int)-1].(string),
					ClientIdentifier: parsedArray[i-2].(string),
					ClientReference:  dataArray[parsedArray[i-3].(int)-1].(string),
					Currency:         dataArray[parsedArray[i-4].(int)-1].(string),
					FileIdentifier:   parsedArray[i-5].(string),
					FileReference:    dataArray[parsedArray[i-6].(int)-1].(string),
					PartyName:        unescapeGWTString(dataArray[parsedArray[i-8].(int)-1].(string)),
					Status:           dataArray[parsedArray[i-9].(int)-1].(string),
					TravelDate:       dataArray[parsedArray[i-10].(int)-1].(string),
				}

				fileSearchResponse.Results[resultIndex] = fileSearchResult
				i = i - 10
				resultIndex++
			}
		}
	}

	return &fileSearchResponse, nil
}
