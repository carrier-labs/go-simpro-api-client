// Package client provides the core HTTP client, authentication, and configuration for SimPro v3 REST API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/carrier-labs/go-simpro-api-client/debug"
)

const (
	DefaultBaseURL = "https://simpro4.wirelesslogic.com/api/v3"
)

// Config holds configuration for the SimPro API client.
type Config struct {
	BaseURL   string        // Optional; if empty, DefaultBaseURL is used
	APIKey    string        // SimPro API Key (required)
	APIClient string        // SimPro API Client identifier (required)
	Timeout   time.Duration // Optional; if zero, 10s is used
}

// Client is the main struct for interacting with the SimPro API.
type Client struct {
	baseURL    string
	apiKey     string
	apiClient  string
	httpClient *http.Client
	mu         sync.Mutex
}

// New creates a new SimPro API client using the provided Config.
func New(cfg Config) *Client {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		baseURL:    baseURL,
		apiKey:     cfg.APIKey,
		apiClient:  cfg.APIClient,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// SetAPIKey allows updating the API key at runtime.
func (c *Client) SetAPIKey(apiKey string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.apiKey = apiKey
	debug.Debug("SetAPIKey called", "api_key_set", apiKey != "")
}

// SetAPIClient allows updating the API client identifier at runtime.
func (c *Client) SetAPIClient(apiClient string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.apiClient = apiClient
	debug.Debug("SetAPIClient called", "api_client_set", apiClient != "")
}

// DoRequest performs an HTTP request with authentication.
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	c.mu.Lock()
	apiKey := c.apiKey
	apiClient := c.apiClient
	c.mu.Unlock()

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			debug.Debug("Failed to marshal request body", "error", err)
			return nil, err
		}
	}
	url := c.baseURL + path
	debug.Debug("Preparing HTTP request", "method", method, "url", url)
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqBody))
	if err != nil {
		debug.Debug("Failed to create HTTP request", "error", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
		debug.Debug("Using API key for authentication")
	}
	if apiClient != "" {
		req.Header.Set("x-api-client", apiClient)
		debug.Debug("Using API client identifier", "x-api-client", apiClient)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		debug.Debug("HTTP request failed", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		debug.Debug("Failed to read response body", "error", err)
		return nil, err
	}
	if resp.StatusCode >= 400 {
		debug.Debug("SimPro API error response", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("SimPro API error: %s", respBody)
	}
	debug.Debug("HTTP request successful", "status", resp.StatusCode, "url", url)
	return respBody, nil
}
