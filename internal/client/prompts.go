package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func normalizePromptCommandForAPI(command string) string {
	if command == "" {
		return command
	}

	trimmed := strings.TrimPrefix(command, "/")
	return "/" + trimmed
}

func promptPathSegment(command string) string {
	return strings.TrimPrefix(command, "/")
}

// PromptForm represents the payload for managing prompt definitions.
type PromptForm struct {
	Command       string         `json:"command"`
	Title         string         `json:"name"`
	Content       string         `json:"content"`
	Data          map[string]any `json:"data,omitempty"`
	Meta          map[string]any `json:"meta,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	CommitMessage *string        `json:"commit_message,omitempty"`
	AccessControl map[string]any `json:"access_control,omitempty"`
}

// PromptModel is returned by the prompt endpoints.
type PromptModel struct {
	ID            *string        `json:"id"`
	Command       string         `json:"command"`
	Title         string         `json:"name"`
	Content       string         `json:"content"`
	Data          map[string]any `json:"data,omitempty"`
	Meta          map[string]any `json:"meta,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	IsActive      *bool          `json:"is_active,omitempty"`
	Timestamp     int64          `json:"timestamp"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
	UserID        string         `json:"user_id"`
	AccessControl map[string]any `json:"access_control,omitempty"`
	AccessGrants  []any          `json:"access_grants,omitempty"`
}

// CreatePrompt registers a new prompt.
func (c *Client) CreatePrompt(ctx context.Context, form PromptForm) (*PromptModel, error) {
	var resp PromptModel
	form.Command = normalizePromptCommandForAPI(form.Command)
	if err := c.do(ctx, http.MethodPost, "prompts/create", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListPrompts returns all prompts.
func (c *Client) ListPrompts(ctx context.Context) ([]PromptModel, error) {
	var resp []PromptModel
	if err := c.do(ctx, http.MethodGet, "prompts/", nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetPrompt fetches a prompt by its command identifier.
func (c *Client) GetPrompt(ctx context.Context, command string) (*PromptModel, error) {
	prompts, err := c.ListPrompts(ctx)
	if err != nil {
		return nil, err
	}

	normalized := normalizePromptCommandForAPI(command)
	for i := range prompts {
		candidate := normalizePromptCommandForAPI(prompts[i].Command)
		if prompts[i].Command == command || candidate == normalized {
			if prompts[i].ID != nil && *prompts[i].ID != "" {
				return c.GetPromptByID(ctx, *prompts[i].ID)
			}
			return &prompts[i], nil
		}
	}

	return nil, ErrNotFound
}

// GetPromptByID fetches a prompt by Open WebUI prompt id.
func (c *Client) GetPromptByID(ctx context.Context, id string) (*PromptModel, error) {
	var resp PromptModel
	path := fmt.Sprintf("prompts/id/%s", url.PathEscape(id))
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdatePrompt updates an existing prompt.
func (c *Client) UpdatePrompt(ctx context.Context, command string, form PromptForm) (*PromptModel, error) {
	current, err := c.GetPrompt(ctx, command)
	if err != nil {
		return nil, err
	}
	if current.ID == nil || *current.ID == "" {
		return nil, fmt.Errorf("prompt %q does not include an id", command)
	}

	var resp PromptModel
	form.Command = normalizePromptCommandForAPI(form.Command)
	path := fmt.Sprintf("prompts/id/%s/update", url.PathEscape(*current.ID))
	if err := c.do(ctx, http.MethodPost, path, nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeletePrompt removes a prompt by command identifier.
func (c *Client) DeletePrompt(ctx context.Context, command string) error {
	current, err := c.GetPrompt(ctx, command)
	if err != nil {
		return err
	}
	if current.ID == nil || *current.ID == "" {
		return fmt.Errorf("prompt %q does not include an id", command)
	}

	path := fmt.Sprintf("prompts/id/%s/delete", url.PathEscape(*current.ID))
	if err := c.do(ctx, http.MethodDelete, path, nil, nil, nil); err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// ListPromptsRaw returns the raw prompt list payload.
func (c *Client) ListPromptsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "prompts/list", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ListPromptTagsRaw returns the raw prompt tags payload.
func (c *Client) ListPromptTagsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "prompts/tags", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPromptHistoryRaw returns raw prompt history for a prompt command.
func (c *Client) GetPromptHistoryRaw(ctx context.Context, command string, page int) (any, error) {
	current, err := c.GetPrompt(ctx, command)
	if err != nil {
		return nil, err
	}
	if current.ID == nil || *current.ID == "" {
		return nil, fmt.Errorf("prompt %q does not include an id", command)
	}
	query := url.Values{}
	if page > 0 {
		query.Set("page", fmt.Sprintf("%d", page))
	}
	var resp any
	path := fmt.Sprintf("prompts/id/%s/history", url.PathEscape(*current.ID))
	if err := c.do(ctx, http.MethodGet, path, query, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TogglePrompt toggles active state by prompt command.
func (c *Client) TogglePrompt(ctx context.Context, command string) (*PromptModel, error) {
	current, err := c.GetPrompt(ctx, command)
	if err != nil {
		return nil, err
	}
	if current.ID == nil || *current.ID == "" {
		return nil, fmt.Errorf("prompt %q does not include an id", command)
	}
	var resp PromptModel
	path := fmt.Sprintf("prompts/id/%s/toggle", url.PathEscape(*current.ID))
	if err := c.do(ctx, http.MethodPost, path, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPromptHistoryEntryRaw returns a raw prompt history entry by history id.
func (c *Client) GetPromptHistoryEntryRaw(ctx context.Context, command string, historyID string) (any, error) {
	current, err := c.GetPrompt(ctx, command)
	if err != nil {
		return nil, err
	}
	if current.ID == nil || *current.ID == "" {
		return nil, fmt.Errorf("prompt %q does not include an id", command)
	}
	var resp any
	path := fmt.Sprintf("prompts/id/%s/history/%s", url.PathEscape(*current.ID), url.PathEscape(historyID))
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPromptHistoryDiffRaw returns a raw diff between two prompt history entries.
func (c *Client) GetPromptHistoryDiffRaw(ctx context.Context, command string, fromID string, toID string) (any, error) {
	current, err := c.GetPrompt(ctx, command)
	if err != nil {
		return nil, err
	}
	if current.ID == nil || *current.ID == "" {
		return nil, fmt.Errorf("prompt %q does not include an id", command)
	}
	query := url.Values{}
	query.Set("from_id", fromID)
	query.Set("to_id", toID)
	var resp any
	path := fmt.Sprintf("prompts/id/%s/history/diff", url.PathEscape(*current.ID))
	if err := c.do(ctx, http.MethodGet, path, query, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
