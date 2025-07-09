package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/config"
	"dodmcdund.cc/panpac-helper/lynxmcpserver/pkg/gwt"
)

const (
	JSESSIONID           = "JSESSIONID"
	AUTH_URL             = "/lynx/service/security.rpc"
	COOKIE_PATH          = "/lynx"
	AUTHORIZATION_HEADER = "Authorization"
)

type SessionContext struct {
	jsessionID string
	expiresAt  time.Time
}

func (s *SessionContext) JSESSIONID() string {
	return s.jsessionID
}

type contextKey string

const sessionContextKey = contextKey("session")

// GetSessionFromContext retrieves the session from context
func GetSessionFromContext(ctx context.Context) (*SessionContext, bool) {
	session, ok := ctx.Value(sessionContextKey).(*SessionContext)
	return session, ok
}

// GetOrCreateSession retrieves an existing session from context or creates a new one if needed
func GetOrCreateSession(ctx context.Context, lynxConfig config.LynxServerConfig) (*SessionContext, context.Context, error) {
	// Try to get existing session
	session, ok := GetSessionFromContext(ctx)
	if ok && session.JSESSIONID() != "" && session.expiresAt.After(time.Now()) {
		return session, ctx, nil
	}

	// No valid session found, create new one
	session, err := makeAuthRequest(lynxConfig)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to login: %w", err)
	}

	// Store new session in context
	ctx = context.WithValue(ctx, sessionContextKey, session)
	return session, ctx, nil
}

// makeAuthRequest performs authentication and returns JSESSIONID
func makeAuthRequest(lynxConfig config.LynxServerConfig) (*SessionContext, error) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	args := &gwt.GWTLoginArgs{
		RemoteHost:  lynxConfig.RemoteHost,
		CompanyCode: lynxConfig.CompanyCode,
		Username:    lynxConfig.Username,
		Password:    lynxConfig.Password,
	}

	body := gwt.BuildGWTLoginBody(args)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s%s", lynxConfig.RemoteHost, AUTH_URL), strings.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", gwt.CONTENT_TYPE)

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to perform auth request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth request failed with status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth response: %w", err)
	}
	bodyStr := string(bodyBytes)

	if !strings.HasPrefix(bodyStr, "//OK") {
		return nil, fmt.Errorf("unexpected response: %s", bodyStr)
	}

	// Extract JSESSIONID from cookies
	for _, cookie := range resp.Cookies() {
		if cookie.Name == JSESSIONID {
			log.Printf("Using JSESSIONID: %s\n", cookie.Value)
			return &SessionContext{
				jsessionID: cookie.Value,
				expiresAt:  time.Now().Add(lynxConfig.AuthCookieDuration),
			}, nil
		}
	}

	return nil, fmt.Errorf("JSESSIONID not found in response cookies")
}

func CreateAuthCookie(lynxConfig config.LynxServerConfig, session *SessionContext) *http.Cookie {
	return &http.Cookie{
		Name:     JSESSIONID,
		Value:    session.JSESSIONID(),
		Domain:   lynxConfig.RemoteHost,
		Path:     COOKIE_PATH,
		Expires:  time.Time{},
		HttpOnly: true,
	}
}

// BearerAuthMiddleware creates middleware that checks for Authorization header
func BearerAuthMiddleware(expectedToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(AUTHORIZATION_HEADER)
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			expectedTokenTrimmed := strings.TrimSpace(expectedToken)
			if token != expectedTokenTrimmed {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
