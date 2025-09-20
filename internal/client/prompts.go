package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PromptForm represents the payload for managing prompt definitions.
type PromptForm struct {
	Command       string         `json:"command"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	AccessControl map[string]any `json:"access_control,omitempty"`
}

// PromptModel is returned by the prompt endpoints.
type PromptModel struct {
	Command       string         `json:"command"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	Timestamp     int64          `json:"timestamp"`
	UserID        string         `json:"user_id"`
	AccessControl map[string]any `json:"access_control,omitempty"`
}

// CreatePrompt registers a new prompt.
func (c *Client) CreatePrompt(ctx context.Context, form PromptForm) (*PromptModel, error) {
	var resp PromptModel
	if err := c.do(ctx, http.MethodPost, "prompts/create", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetPrompt fetches a prompt by its command identifier.
func (c *Client) GetPrompt(ctx context.Context, command string) (*PromptModel, error) {
	var resp PromptModel
	path := fmt.Sprintf("prompts/command/%s", url.PathEscape(command))
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdatePrompt updates an existing prompt.
func (c *Client) UpdatePrompt(ctx context.Context, command string, form PromptForm) (*PromptModel, error) {
	var resp PromptModel
	path := fmt.Sprintf("prompts/command/%s/update", url.PathEscape(command))
	if err := c.do(ctx, http.MethodPost, path, nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeletePrompt removes a prompt by command identifier.
func (c *Client) DeletePrompt(ctx context.Context, command string) error {
	path := fmt.Sprintf("prompts/command/%s/delete", url.PathEscape(command))
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
