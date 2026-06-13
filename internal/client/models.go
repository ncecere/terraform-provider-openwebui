package client

import (
	"context"
	"net/http"
	"net/url"
)

// ModelForm represents the payload for creating or updating models.
type ModelForm struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Meta          map[string]any `json:"meta"`
	Params        map[string]any `json:"params"`
	BaseModelID   *string        `json:"base_model_id,omitempty"`
	IsActive      *bool          `json:"is_active,omitempty"`
	AccessControl map[string]any `json:"access_control,omitempty"`
}

// ModelResponse captures details returned by the model endpoints.
type ModelResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	Meta          map[string]any `json:"meta"`
	Params        map[string]any `json:"params"`
	BaseModelID   *string        `json:"base_model_id,omitempty"`
	IsActive      bool           `json:"is_active"`
	AccessControl map[string]any `json:"access_control,omitempty"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
}

// CreateModel registers a new model.
func (c *Client) CreateModel(ctx context.Context, form ModelForm) (*ModelResponse, error) {
	var resp ModelResponse
	if err := c.do(ctx, http.MethodPost, "models/create", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetModel obtains model details by identifier.
func (c *Client) GetModel(ctx context.Context, id string) (*ModelResponse, error) {
	var resp ModelResponse
	query := url.Values{"id": []string{id}}
	if err := c.do(ctx, http.MethodGet, "models/model", query, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateModel updates a model by identifier.
func (c *Client) UpdateModel(ctx context.Context, id string, form ModelForm) (*ModelResponse, error) {
	var resp ModelResponse
	query := url.Values{"id": []string{id}}
	if err := c.do(ctx, http.MethodPost, "models/model/update", query, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ModelIDForm identifies a model for endpoints that expect a JSON body.
type ModelIDForm struct {
	ID string `json:"id"`
}

// DeleteModel removes a model by identifier.
func (c *Client) DeleteModel(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, "models/model/delete", nil, ModelIDForm{ID: id}, nil)
}

// ListModelsRaw returns the raw models list payload.
func (c *Client) ListModelsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "models", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ListBaseModelsRaw returns the raw base models payload.
func (c *Client) ListBaseModelsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "models/base", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ListModelTagsRaw returns the raw model tags payload.
func (c *Client) ListModelTagsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "models/tags", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ExportModelsRaw returns the raw model export payload.
func (c *Client) ExportModelsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "models/export", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
