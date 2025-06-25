package config

import "os"

type ClientConfig struct {
	BearerToken string
}

func NewClientConfig() ClientConfig {
	return ClientConfig{
		BearerToken: os.Getenv("BEARER_TOKEN"),
	}
}
