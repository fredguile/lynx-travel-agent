package config

import (
	"os"
	"time"
)

type LynxServerConfig struct {
	RemoteHost         string
	AuthCookieDuration time.Duration
	Username           string
	Password           string
	CompanyCode        string
}

func NewLynxServerConfig() LynxServerConfig {
	if os.Getenv("LYNX_USERNAME") == "" {
		panic("LYNX_USERNAME is not set")
	}

	if os.Getenv("LYNX_PASSWORD") == "" {
		panic("LYNX_PASSWORD is not set")
	}

	if os.Getenv("LYNX_COMPANY_CODE") == "" {
		panic("LYNX_COMPANY_CODE is not set")
	}

	return LynxServerConfig{
		RemoteHost:         "www.lynx-reservations.com",
		AuthCookieDuration: 15 * time.Minute,
		Username:           os.Getenv("LYNX_USERNAME"),
		Password:           os.Getenv("LYNX_PASSWORD"),
		CompanyCode:        os.Getenv("LYNX_COMPANY_CODE"),
	}
}
