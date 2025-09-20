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
	var resp PromptModel
	apiCommand := promptPathSegment(command)
	path := fmt.Sprintf("prompts/command/%s", url.PathEscape(apiCommand))
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err == nil {
		return &resp, nil
	} else {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Status == http.StatusMethodNotAllowed {
			return nil, err
		}

		if errors.Is(err, ErrNotFound) {
			return nil, err
		}

		if errors.As(err, &apiErr) && (apiErr.Status == http.StatusUnauthorized || apiErr.Status == http.StatusNotFound) {
			prompts, listErr := c.ListPrompts(ctx)
			if listErr != nil {
				return nil, err
			}

			normalized := normalizePromptCommandForAPI(command)
			for i := range prompts {
				candidate := normalizePromptCommandForAPI(prompts[i].Command)
				if prompts[i].Command == command || candidate == normalized {
					return &prompts[i], nil
				}
			}

			return nil, ErrNotFound
		}

		return nil, err
	}
}

// UpdatePrompt updates an existing prompt.
func (c *Client) UpdatePrompt(ctx context.Context, command string, form PromptForm) (*PromptModel, error) {
	var resp PromptModel
	apiCommand := promptPathSegment(command)
	form.Command = normalizePromptCommandForAPI(form.Command)
	path := fmt.Sprintf("prompts/command/%s/update", url.PathEscape(apiCommand))
	err := c.do(ctx, http.MethodPost, path, nil, form, &resp)
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.Status == http.StatusMethodNotAllowed {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeletePrompt removes a prompt by command identifier.
func (c *Client) DeletePrompt(ctx context.Context, command string) error {
	apiCommand := promptPathSegment(command)
	path := fmt.Sprintf("prompts/command/%s/delete", url.PathEscape(apiCommand))
	if err := c.do(ctx, http.MethodDelete, path, nil, nil, nil); err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}

		var apiErr *APIError
		if errors.As(err, &apiErr) && (apiErr.Status == http.StatusUnauthorized || apiErr.Status == http.StatusNotFound) {
			prompts, listErr := c.ListPrompts(ctx)
			if listErr != nil {
				return err
			}

			normalized := normalizePromptCommandForAPI(command)
			for _, p := range prompts {
				if p.Command == command || normalizePromptCommandForAPI(p.Command) == normalized {
					altPath := fmt.Sprintf("prompts/command/%s/delete", url.PathEscape(promptPathSegment(p.Command)))
					if delErr := c.do(ctx, http.MethodDelete, altPath, nil, nil, nil); delErr != nil {
						if errors.Is(delErr, ErrNotFound) {
							return ErrNotFound
						}

						var altAPIError *APIError
						if errors.As(delErr, &altAPIError) && (altAPIError.Status == http.StatusUnauthorized || altAPIError.Status == http.StatusNotFound) {
							return ErrNotFound
						}

						return delErr
					}

					return nil
				}
			}

			return ErrNotFound
		}

		return err
	}

	return nil
}
