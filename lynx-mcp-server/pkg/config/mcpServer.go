package config

import "os"

type MCPServerConfig struct {
	Name        string
	Version     string
	Port        string
	BearerToken string
}

func NewMCPServerConfig() MCPServerConfig {
	if os.Getenv("BEARER_TOKEN") == "" {
		panic("BEARER_TOKEN is not set")
	}

	return MCPServerConfig{
		Name:        "lynx-mcp-server",
		Version:     "1.0.0",
		Port:        "9600",
		BearerToken: os.Getenv("BEARER_TOKEN"),
	}
}
