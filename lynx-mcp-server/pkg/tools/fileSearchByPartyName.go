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
	TOOL_FILE_SEARCH_BY_PARTY_NAME                            string = "file_search_by_party_name"
	TOOL_FILE_SEARCH_BY_PARTY_NAME_DESCRIPTION                string = "Retrieve file from party name"
	TOOL_FILE_SEARCH_BY_PARTY_NAME_ARG_PARTY_NAME             string = "partyName"
	TOOL_FILE_SEARCH_BY_PARTY_NAME_ARG_PARTY_NAME_DESCRIPTION string = "Party name"
)

const (
	FILE_SEARCH_URL = "/lynx/service/file.rpc"
)

// FileSearchResponse represents the structured response for file search
type FileSearchResponse struct {
	Count   int                `json:"count"`
	Results []FileSearchResult `json:"results"`
}

// FileSearchResult represents a single file search result
type FileSearchResult struct {
	CompanyCode     string `json:"companyCode"`
	ClientReference string `json:"clientReference"`
	Currency        string `json:"currency"`
	FileIdentifier  string `json:"fileIdentifier"`
	FileReference   string `json:"fileReference"`
	PartyName       string `json:"partyName"`
	Status          string `json:"status"`
	TravelDate      string `json:"traveDate"`
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
	responseBody, err := utils.ParseGWTResponseBody(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GWT response: %w", err)
	}

	// Convert parsed data to structured format
	fileSearchResponse, err := parseFileSearchResponse(responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file search response: %w", err)
	}

	// Convert to JSON for MCP response
	jsonData, err := json.Marshal(fileSearchResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// parseFileSearchResponse converts the parsed GWT data into structured FileSearchResponse object
func parseFileSearchResponse(responseBody any) (*FileSearchResponse, error) {
	if gwtArrayResult, ok := responseBody.(utils.GWTArrayResult); ok {
		fileSearchResponse := FileSearchResponse{
			Count:   gwtArrayResult.Size,
			Results: make([]FileSearchResult, gwtArrayResult.Size),
		}

		for i, item := range gwtArrayResult.Items {
			if gwtFileSearchResult, ok := item.(utils.GWTFileSearchResult); ok {
				fileSearchResult := FileSearchResult{
					CompanyCode:     gwtFileSearchResult.CompanyCode,
					ClientReference: gwtFileSearchResult.ClientReference,
					Currency:        gwtFileSearchResult.Currency,
					FileIdentifier:  gwtFileSearchResult.FileIdentifier,
					FileReference:   gwtFileSearchResult.FileReference,
					PartyName:       gwtFileSearchResult.PartyName,
					Status:          gwtFileSearchResult.Status,
					TravelDate:      gwtFileSearchResult.TravelDate,
				}
				fileSearchResponse.Results[i] = fileSearchResult
			} else {
				return nil, fmt.Errorf("invalid GWTFileSearchResult")
			}
		}

		return &fileSearchResponse, nil
	} else {
		return nil, fmt.Errorf("invalid GWTArrayResult")
	}
}
