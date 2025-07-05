package tools

import (
	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/config"
)

var lynxConfig = config.NewLynxServerConfig()

type ToolName string

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
