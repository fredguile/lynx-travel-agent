package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestUnescapeGWTErrorString(t *testing.T) {
	// Test the example provided by the user
	testInput := `"Type \x27com.lynxtraveltech.client.shared.model.FileSearchCriteria\x27 was not assignable to \x27com.google.gwt.user.client.rpc.IsSerializable\x27 and did not have a custom field serializer. For security purposes"`

	expected := `Type 'com.lynxtraveltech.client.shared.model.FileSearchCriteria' was not assignable to 'com.google.gwt.user.client.rpc.IsSerializable' and did not have a custom field serializer. For security purposes`

	result := unescapeGWTErrorString(testInput)

	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
		fmt.Printf("Expected length: %d, Got length: %d\n", len(expected), len(result))
	}
}

func TestParseResponseError(t *testing.T) {
	// Test case with the example error provided
	testError := `//EX[2,1,["com.google.gwt.user.client.rpc.IncompatibleRemoteServiceException/3936916533","Type \x27com.lynxtraveltech.client.shared.model.FileSearchCriteria\x27 was not assignable to \x27com.google.gwt.user.client.rpc.IsSerializable\x27 and did not have a custom field serializer. For security purposes, this type will not be deserialized."],0,7]`

	expectedError := `Type 'com.lynxtraveltech.client.shared.model.FileSearchCriteria' was not assignable to 'com.google.gwt.user.client.rpc.IsSerializable' and did not have a custom field serializer. For security purposes, this type will not be deserialized.`

	result, err := ParseResponseError(testError)
	if err != nil {
		t.Fatalf("ParseResponseError failed: %v", err)
	}

	actual := strings.TrimSpace(result)
	expected := strings.TrimSpace(expectedError)

	if actual != expected {
		t.Errorf("Expected error message: %s\nGot: %s", expected, actual)
		fmt.Printf("Expected length: %d, Got length: %d\n", len(expected), len(actual))
	}
}
