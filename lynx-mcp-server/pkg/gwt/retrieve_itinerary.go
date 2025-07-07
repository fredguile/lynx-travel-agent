package gwt

import (
	"fmt"
	"strings"
)

type RetrieveItineraryArgs struct {
	RemoteHost     string
	FileIdentifier string
}

// BuildRetrieveItineraryGWTBody constructs the GWT-RPC retrieve itinerary body with the given file identifier
func BuildRetrieveItineraryGWTBody(args *RetrieveItineraryArgs) string {
	return fmt.Sprintf("7|0|6|https://%s/lynx/lynx/|63A734E3E71C14883B20AFEC1238F6A7|com.lynxtraveltech.client.client.rpc.FileService|retrieveItinerary|J|Z|1|2|3|4|4|5|6|6|6|%s|0|0|0|", args.RemoteHost, args.FileIdentifier)
}

type RetrieveItineraryResponseArray struct {
	Type             string                        `json:"type"`
	PartyName        string                        `json:"partyName"`
	FileReference    string                        `json:"fileReference"`
	FileIdentifier   string                        `json:"fileIdentifier"`
	ClientIdentifier string                        `json:"clientIdentifier"`
	AgentReference   string                        `json:"agentReference"`
	ItineraryCount   int                           `json:"itineraryCount"`
	Itineraries      []ItineraryTransactionSummary `json:"itineraries"`
}

type ItineraryTransactionSummary struct {
	VoucherIdentifier     string `json:"voucherIdentifier"`
	Date                  string `json:"date"`
	TransactionIdentifier string `json:"transactionIdentifier"`
	Supplier              string `json:"supplier"`
	Status                string `json:"status"`
	ConfirmationNumber    string `json:"confirmationNumber"`
	Location              string `json:"location"`
}

// ParseGWTFileSearchResponseBody parses a GWT response in context of retrieve itinerary.
// Returns the parsed data as RetrieveItineraryResponseArray containing the array elements.
func ParseRetrieveItineraryResponseBody(responseBody string) (*RetrieveItineraryResponseArray, error) {
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

	retrieveItineraryResponse := RetrieveItineraryResponseArray{
		Type:             dataArray[parsedArray[1].(int)-1].(string),
		PartyName:        unescapeGWTString(dataArray[parsedArray[2].(int)-1].(string)),
		FileReference:    dataArray[parsedArray[4].(int)-1].(string),
		FileIdentifier:   parsedArray[5].(string),
		AgentReference:   dataArray[parsedArray[7].(int)-1].(string),
		ClientIdentifier: parsedArray[8].(string),
		ItineraryCount:   0,
		Itineraries:      make([]ItineraryTransactionSummary, 0),
	}

	currentItinerary := ItineraryTransactionSummary{}

	for i, rel := 11, 0; i < len(parsedArray); i, rel = i+1, rel+1 {
		if stringValue, ok := parsedArray[i].(string); ok {
			switch rel {
			case 3:
				currentItinerary.TransactionIdentifier = stringValue
			}

			// FIXME: little hacky
			if currentItinerary.TransactionIdentifier == "" &&
				len(stringValue) > 3 &&
				strings.HasPrefix(stringValue, "B") &&
				len(stringValue) > 2 &&
				stringValue[1:2] >= "a" &&
				stringValue[1:2] <= "z" {
				currentItinerary.TransactionIdentifier = stringValue
			}

		}

		if oneBasedIndex, ok := parsedArray[i].(int); ok {
			// Ignore values that don't match a oneBasedIndex
			if oneBasedIndex <= 0 || oneBasedIndex >= len(dataArray) {
				continue
			}

			currentValue := dataArray[oneBasedIndex-1]

			if currentStringValue, ok := currentValue.(string); ok &&
				(strings.HasPrefix(currentStringValue, GWT_TYPE_BIGDECIMAL) ||
					strings.HasPrefix(currentStringValue, GWT_TYPE_SQL_DATE) ||
					strings.HasPrefix(currentStringValue, GWT_TYPE_DOUBLE) ||
					strings.HasPrefix(currentStringValue, GWT_TYPE_LONG) ||
					strings.HasPrefix(currentStringValue, GWT_TYPE_STRING)) {
				// skip next value
				i++
				continue
			}

			switch rel {
			case 0:
				currentItinerary.VoucherIdentifier = currentValue.(string)
			case 2:
				currentItinerary.Date = currentValue.(string)
			case 5:
				currentItinerary.Supplier = currentValue.(string)
			case 9:
				currentItinerary.Status = currentValue.(string)
			}

			if currentStringValue, ok := currentValue.(string); ok && strings.HasPrefix(currentStringValue, GWT_TYPE_TRANSACTION_SUMMARY) {
				// retrieve confirmation number
				oneBasedIndexCfrnNum, ok := parsedArray[i-9].(int)
				if ok && oneBasedIndexCfrnNum > 0 && oneBasedIndexCfrnNum < len(dataArray) {
					currentItinerary.ConfirmationNumber = dataArray[oneBasedIndexCfrnNum-1].(string)
				}

				// retrieve location
				oneBasedIndexLoc, ok := parsedArray[i-14].(int)
				if ok && oneBasedIndexLoc > 0 && oneBasedIndexLoc < len(dataArray) {
					currentItinerary.Location = dataArray[oneBasedIndexLoc-1].(string)
				}

				retrieveItineraryResponse.ItineraryCount += 1
				retrieveItineraryResponse.Itineraries = append(retrieveItineraryResponse.Itineraries, currentItinerary)

				currentItinerary = ItineraryTransactionSummary{}

				// reset j on next iteration
				rel = -1
			}
		}
	}

	// reverse order of itineries
	for i, j := 0, len(retrieveItineraryResponse.Itineraries)-1; i < j; i, j = i+1, j-1 {
		retrieveItineraryResponse.Itineraries[i], retrieveItineraryResponse.Itineraries[j] = retrieveItineraryResponse.Itineraries[j], retrieveItineraryResponse.Itineraries[i]
	}

	return &retrieveItineraryResponse, nil
}
