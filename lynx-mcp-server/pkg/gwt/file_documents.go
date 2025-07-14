package gwt

import (
	"fmt"
	"strings"
)

type BuildFileDocumentsByFileReferenceArgs struct {
	RemoteHost     string
	FileIdentifier string
}

func BuildFileDocumentsByFileReferenceGWTBody(args *BuildFileDocumentsByFileReferenceArgs) string {
	return fmt.Sprintf("7|0|8|https://%s/lynx/lynx/|63A734E3E71C14883B20AFEC1238F6A7|com.lynxtraveltech.client.client.rpc.FileService|getFileDocumentsAsList|J|java.lang.Long/4227064769|I|java.lang.String/2004016611|1|2|3|4|4|5|6|7|8|%s|0|1|0|", args.RemoteHost, args.FileIdentifier)
}

type FileDocumentsByTransactionReferenceArgs struct {
	RemoteHost            string
	FileIdentifier        string
	TransactionIdentifier string
}

func BuildFileDocumentsByTransactionReferenceGWTBody(args *FileDocumentsByTransactionReferenceArgs) string {
	return fmt.Sprintf("7|0|8|https://%s/lynx/lynx/|63A734E3E71C14883B20AFEC1238F6A7|com.lynxtraveltech.client.client.rpc.FileService|getFileDocumentsAsList|J|java.lang.Long/4227064769|I|java.lang.String/2004016611|1|2|3|4|4|5|6|7|8|%s|6|%s|1|0|", args.RemoteHost, args.FileIdentifier, args.TransactionIdentifier)
}

type TransactionDocumentSaveDetailsArgs struct {
	RemoteHost            string
	FileIdentifier        string
	TransactionIdentifier string
	Name                  string
	Content               string
	Type                  string
	AttachmentURL         string
}

func BuildTransactionDocumentSaveGWTBody(args *TransactionDocumentSaveDetailsArgs) string {
	return fmt.Sprintf("7|0|10|https://%s/lynx/lynx/|63A734E3E71C14883B20AFEC1238F6A7|com.lynxtraveltech.client.client.rpc.FileService|saveFileDocumentsDetails|com.lynxtraveltech.common.gui.shared.model.DocumentDetails/2779362264|java.lang.Long/4227064769|%s|%s|%s|%s|1|2|3|4|1|5|5|6|%s|1|A|0|0|0|P__________|7|8|0|%s|9|10|0|", args.RemoteHost, args.Content, args.Type, args.Name, args.AttachmentURL, args.TransactionIdentifier, args.FileIdentifier)
}

type FileDocumentSaveDetailsArgs struct {
	RemoteHost     string
	FileIdentifier string
	Name           string
	Content        string
	Type           string
	AttachmentURL  string
}

func BuildFileDocumentSaveGWTBody(args *FileDocumentSaveDetailsArgs) string {
	return fmt.Sprintf("7|0|9|https://%s/lynx/lynx/|63A734E3E71C14883B20AFEC1238F6A7|com.lynxtraveltech.client.client.rpc.FileService|saveFileDocumentsDetails|com.lynxtraveltech.common.gui.shared.model.DocumentDetails/2779362264|%s|%s|%s|%s|1|2|3|4|1|5|5|0|1|A|0|0|0|P__________|6|7|0|%s|8|9|0|", args.RemoteHost, args.Content, args.Type, args.Name, args.AttachmentURL, args.FileIdentifier)
}

type FileDocumentsResponseArray struct {
	Count   int            `json:"count"`
	Results []FileDocument `json:"results"`
}

type FileDocument struct {
	FileIdentifier        string `json:"fileIdentifier"`
	TransactionIdentifier string `json:"transactionIdentifier"`
	DocumentIdentifier    string `json:"documentIdentifier"`
	DocumentName          string `json:"documentName"`
	DocumentType          string `json:"documentType"`
	Content               string `json:"content"`
	AttachmentUrl         string `json:"attachmentUrl"`
}

// ParseFileDocumentsListResponseBody parses a GWT response in context of listing file documents
// Returns the parsed data as FileDocumentsResponseArray containing the array elements.
func ParseFileDocumentsListResponseBody(responseBody string) (*FileDocumentsResponseArray, error) {
	// Ensure response starts with "//OK"
	if !strings.HasPrefix(responseBody, "//OK") {
		return nil, fmt.Errorf("response body missing //OK")
	}

	// Remove the "//OK" prefix
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

	fileDocumentsResponse := FileDocumentsResponseArray{
		Count:   arraySize,
		Results: make([]FileDocument, arraySize),
	}

	for i, resultIndex := len(parsedArray)-6, 0; i >= 0; i-- {
		if oneBasedIndex, ok := parsedArray[i].(int); ok {
			// Ignore values that don't match a oneBasedIndex
			if oneBasedIndex <= 0 || oneBasedIndex >= len(dataArray) {
				continue
			}

			currentValue := dataArray[oneBasedIndex-1]

			if currentStringValue, ok := currentValue.(string); ok && strings.HasPrefix(currentStringValue, GWT_TYPE_DOCUMENT_DETAILS) {
				fileDocument := FileDocument{
					TransactionIdentifier: parsedArray[i-2].(string),
					Content:               unescapeGWTString(dataArray[parsedArray[i-11].(int)-1].(string)),
					DocumentType:          dataArray[parsedArray[i-13].(int)-1].(string),
					FileIdentifier:        parsedArray[i-14].(string),
					DocumentName:          dataArray[parsedArray[i-15].(int)-1].(string),
					DocumentIdentifier:    dataArray[parsedArray[i-17].(int)-1].(string),
				}

				if oneBasedIndex, ok := parsedArray[i-16].(int); ok && oneBasedIndex > 0 {
					fileDocument.AttachmentUrl = dataArray[parsedArray[i-16].(int)-1].(string)
				}

				fileDocumentsResponse.Results[resultIndex] = fileDocument
				i -= 17
				resultIndex++
			}
		}
	}

	return &fileDocumentsResponse, nil
}

// ParseDocumentSaveResponseBody parsed a GWT response in content of saving document details
// Returns nil if successful
func ParseDocumentSaveResponseBody(responseBody string) error {
	// Ensure response starts with "//OK"
	if !strings.HasPrefix(responseBody, "//OK") {
		return fmt.Errorf("response body missing //OK")
	}

	return nil
}
