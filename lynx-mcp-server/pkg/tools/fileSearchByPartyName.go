package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/config"
	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

var lynxConfig = config.NewLynxServerConfig()

const (
	TOOL_FILE_SEARCH_BY_PARTY_NAME                                   = "file_search_by_party_name"
	TOOL_FILE_SEARCH_BY_PARTY_NAME_DESCRIPTION                string = "Retrieve file from party name"
	TOOL_FILE_SEARCH_BY_PARTY_NAME_ARG_PARTY_NAME             string = "partyName"
	TOOL_FILE_SEARCH_BY_PARTY_NAME_ARG_PARTY_NAME_DESCRIPTION string = "Party name"
)

const (
	FILE_SEARCH_URL = "/lynx/service/file.rpc"
)

// FileSearchResult represents a single file search result
type FileSearchResult struct {
	CompanyCode string `json:"companyCode"`
	Currency    string `json:"currency"`
	FileNumber  string `json:"fileNumber"`
	PartyName   string `json:"partyName"`
	Status      string `json:"status"`
	Date        string `json:"date"`
}

// FileSearchResponse represents the structured response for file search
type FileSearchResponse struct {
	Results []FileSearchResult `json:"results"`
	Count   int                `json:"count"`
}

func HandleFileSearchByPartyName(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	session, _, err := utils.GetOrCreateSession(ctx, lynxConfig)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	arguments := request.GetArguments()
	partyName, ok := arguments["partyName"].(string)

	if !ok {
		return nil, fmt.Errorf("invalid number arguments")
	}

	body := utils.BuildGWTFileSearchBody(&utils.GWTFileSearchArgs{
		PartyName: partyName,
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, FILE_SEARCH_URL), strings.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("failed to create file search request: %w", err)
	}

	req.Header.Set("Content-Type", utils.GWT_CONTENT_TYPE)
	req.AddCookie(utils.CreateAuthCookie(lynxConfig, session))

	// Use retry utility with exponential backoff
	resp, bodyStr, err := utils.RetryHTTPRequest(ctx, client, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute file search request after retries: %w", err)
	}
	defer resp.Body.Close()

	// Parse the GWT response body
	parsedData, err := utils.ParseResponseBody(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GWT response: %w", err)
	}

	// Convert parsed data to structured format
	fileResult := parseFileSearchResult(parsedData)

	// Create structured response
	var results []FileSearchResult
	if fileResult != nil {
		results = []FileSearchResult{*fileResult}
	}

	response := FileSearchResponse{
		Results: results,
		Count:   len(results),
	}

	// Convert to JSON for MCP response
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// parseFileSearchResults converts the parsed GWT data into structured FileSearchResult object
func parseFileSearchResult(parsedData []interface{}) *FileSearchResult {
	// Ensure parsedData is a slice/array
	if len(parsedData) >= 9 {
		return &FileSearchResult{
			CompanyCode: toString(parsedData[2]),
			Currency:    toString(parsedData[4]),
			FileNumber:  toString(parsedData[5]),
			PartyName:   toString(parsedData[6]),
			Status:      toString(parsedData[7]),
			Date:        toString(parsedData[8]),
		}
	}

	return nil
}

// toString safely converts an interface{} to string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
