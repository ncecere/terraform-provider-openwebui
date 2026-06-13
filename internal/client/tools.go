package client

import (
	"context"
	"net/http"
	"net/url"
)

// ToolMeta captures descriptive metadata for a tool.
type ToolMeta struct {
	Description *string        `json:"description,omitempty"`
	Manifest    map[string]any `json:"manifest,omitempty"`
}

// ToolForm represents the payload for creating or updating tools.
type ToolForm struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Content       string         `json:"content"`
	Meta          ToolMeta       `json:"meta"`
	AccessControl map[string]any `json:"access_control,omitempty"`
}

// ToolResponse captures basic tool details.
type ToolResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	Meta          ToolMeta       `json:"meta"`
	AccessControl map[string]any `json:"access_control,omitempty"`
	UpdatedAt     int64          `json:"updated_at"`
	CreatedAt     int64          `json:"created_at"`
}

// ToolModel includes full tool content and specifications.
type ToolModel struct {
	ID            string           `json:"id"`
	UserID        string           `json:"user_id"`
	Name          string           `json:"name"`
	Content       string           `json:"content"`
	Specs         []map[string]any `json:"specs"`
	Meta          ToolMeta         `json:"meta"`
	AccessControl map[string]any   `json:"access_control,omitempty"`
	UpdatedAt     int64            `json:"updated_at"`
	CreatedAt     int64            `json:"created_at"`
}

// ToolAccessResponse captures tool details with access metadata.
type ToolAccessResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	Meta          ToolMeta       `json:"meta"`
	AccessControl map[string]any `json:"access_control,omitempty"`
	UpdatedAt     int64          `json:"updated_at"`
	CreatedAt     int64          `json:"created_at"`
	User          *User          `json:"user,omitempty"`
	WriteAccess   *bool          `json:"write_access,omitempty"`
}

// ToolUserResponse captures tool details including user metadata.
type ToolUserResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	Meta          ToolMeta       `json:"meta"`
	AccessControl map[string]any `json:"access_control,omitempty"`
	UpdatedAt     int64          `json:"updated_at"`
	CreatedAt     int64          `json:"created_at"`
	User          *User          `json:"user,omitempty"`
}

// CreateTool provisions a new tool.
func (c *Client) CreateTool(ctx context.Context, form ToolForm) (*ToolResponse, error) {
	var resp ToolResponse
	if err := c.do(ctx, http.MethodPost, "tools/create", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetTool retrieves a tool by ID.
func (c *Client) GetTool(ctx context.Context, id string) (*ToolAccessResponse, error) {
	var resp ToolAccessResponse
	path := "tools/id/" + url.PathEscape(id)
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateTool updates an existing tool by ID.
func (c *Client) UpdateTool(ctx context.Context, id string, form ToolForm) (*ToolModel, error) {
	var resp ToolModel
	path := "tools/id/" + url.PathEscape(id) + "/update"
	if err := c.do(ctx, http.MethodPost, path, nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeleteTool removes a tool by ID.
func (c *Client) DeleteTool(ctx context.Context, id string) error {
	path := "tools/id/" + url.PathEscape(id) + "/delete"
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListTools returns tool summaries.
func (c *Client) ListTools(ctx context.Context) ([]ToolUserResponse, error) {
	var resp []ToolUserResponse
	if err := c.do(ctx, http.MethodGet, "tools/", nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ListToolAccess returns tool access summaries.
func (c *Client) ListToolAccess(ctx context.Context) ([]ToolAccessResponse, error) {
	var resp []ToolAccessResponse
	if err := c.do(ctx, http.MethodGet, "tools/list", nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ExportTools returns full tool payloads including content.
func (c *Client) ExportTools(ctx context.Context) ([]ToolModel, error) {
	var resp []ToolModel
	if err := c.do(ctx, http.MethodGet, "tools/export", nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetToolValves retrieves valve settings for a tool.
func (c *Client) GetToolValves(ctx context.Context, id string) (map[string]any, error) {
	var resp map[string]any
	path := "tools/id/" + url.PathEscape(id) + "/valves"
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetToolValvesSpec retrieves the valve specification for a tool.
func (c *Client) GetToolValvesSpec(ctx context.Context, id string) (map[string]any, error) {
	var resp map[string]any
	path := "tools/id/" + url.PathEscape(id) + "/valves/spec"
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateToolValves updates valve settings for a tool.
func (c *Client) UpdateToolValves(ctx context.Context, id string, valves map[string]any) (map[string]any, error) {
	var resp map[string]any
	path := "tools/id/" + url.PathEscape(id) + "/valves/update"
	if err := c.do(ctx, http.MethodPost, path, nil, valves, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ListToolsRaw returns the raw tool list payload.
func (c *Client) ListToolsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "tools/list", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ExportToolsRaw returns the raw tools export payload.
func (c *Client) ExportToolsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "tools/export", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
