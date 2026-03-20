package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/certimate-go/cli/internal/config"
)

// Client is the API client for Certimate
type Client struct {
	httpClient *http.Client
	server     string
	token      string
	insecure   bool
}

// NewClient creates a new API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		server: cfg.Server,
		token:  cfg.Token,
	}
}

// SetInsecure skips TLS verification
func (c *Client) SetInsecure(insecure bool) {
	c.insecure = insecure
}

// Do performs an HTTP request
func (c *Client) Do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	url := c.server + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &NetworkError{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		bodyBytes, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
			apiErr = APIError{
				Message: string(bodyBytes),
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// ValidateConnection tests if the connection is valid
func (c *Client) ValidateConnection(ctx context.Context) error {
	// Try to fetch workflows as a connection test
	_, err := c.ListWorkflows(ctx, "", "", 1, 1)
	if err != nil {
		return fmt.Errorf("connection validation failed: %w", err)
	}
	return nil
}
