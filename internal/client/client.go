package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrNotFound indicates that the requested resource could not be located.
var ErrNotFound = errors.New("openwebui: resource not found")

// APIError represents a non-2xx HTTP response from the Open WebUI API.
type APIError struct {
	Status int
	Body   string
}

func (e *APIError) Error() string {
	if strings.TrimSpace(e.Body) == "" {
		return fmt.Sprintf("openwebui: unexpected status code %d", e.Status)
	}

	return fmt.Sprintf("openwebui: status %d: %s", e.Status, e.Body)
}

// Client wraps HTTP access to the Open WebUI API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient constructs a new API client instance.
func NewClient(endpoint, token string) (*Client, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint must be provided")
	}

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %w", err)
	}

	parsed.RawQuery = ""
	parsed.Fragment = ""

	base := strings.TrimRight(parsed.String(), "/")

	return &Client{
		baseURL: base,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// do performs an HTTP request against the API.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, payload any, out any) error {
	var body io.Reader

	if payload != nil {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(payload); err != nil {
			return fmt.Errorf("encode request body: %w", err)
		}
		body = buf
	}

	fullURL := c.baseURL
	trimmedPath := strings.TrimLeft(path, "/")
	if trimmedPath != "" {
		fullURL = fmt.Sprintf("%s/%s", fullURL, trimmedPath)
	}

	if query != nil {
		encoded := query.Encode()
		if encoded != "" {
			fullURL = fullURL + "?" + encoded
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		io.Copy(io.Discard, resp.Body)
		return ErrNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return &APIError{Status: resp.StatusCode, Body: strings.TrimSpace(string(respBody))}
	}

	if out == nil {
		io.Copy(io.Discard, resp.Body)
		return nil
	}

	// Read the full body to gracefully handle empty or null responses.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		return nil
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}
