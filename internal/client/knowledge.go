package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// KnowledgeForm models the payload for creating or updating knowledge records.
type KnowledgeForm struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	AccessControl map[string]any `json:"access_control,omitempty"`
	Data          map[string]any `json:"data,omitempty"`
	Meta          map[string]any `json:"meta,omitempty"`
}

// FileModel represents a file associated with a knowledge base entry.
type FileModel struct {
	ID        string         `json:"id"`
	UserID    string         `json:"user_id"`
	Filename  string         `json:"filename"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	Hash      *string        `json:"hash,omitempty"`
	Path      *string        `json:"path,omitempty"`
	Meta      map[string]any `json:"meta,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
}

// KnowledgeResponse captures the core knowledge object returned by create and update endpoints.
type KnowledgeResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
	AccessControl map[string]any `json:"access_control,omitempty"`
	Data          map[string]any `json:"data,omitempty"`
	Meta          map[string]any `json:"meta,omitempty"`
	Files         []FileModel    `json:"files,omitempty"`
}

// KnowledgeFilesResponse is returned by the knowledge detail endpoint and includes file metadata.
type KnowledgeFilesResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
	AccessControl map[string]any `json:"access_control,omitempty"`
	Data          map[string]any `json:"data,omitempty"`
	Meta          map[string]any `json:"meta,omitempty"`
	Files         []FileModel    `json:"files"`
}

// CreateKnowledge provisions a new knowledge base entry.
func (c *Client) CreateKnowledge(ctx context.Context, form KnowledgeForm) (*KnowledgeResponse, error) {
	var resp KnowledgeResponse
	if err := c.do(ctx, http.MethodPost, "knowledge/create", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetKnowledge retrieves a knowledge record by identifier.
func (c *Client) GetKnowledge(ctx context.Context, id string) (*KnowledgeFilesResponse, error) {
	var resp KnowledgeFilesResponse
	path := fmt.Sprintf("knowledge/%s", url.PathEscape(id))
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateKnowledge mutates an existing knowledge record.
func (c *Client) UpdateKnowledge(ctx context.Context, id string, form KnowledgeForm) (*KnowledgeResponse, error) {
	var resp KnowledgeResponse
	path := fmt.Sprintf("knowledge/%s/update", url.PathEscape(id))
	if err := c.do(ctx, http.MethodPost, path, nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeleteKnowledge removes a knowledge record.
func (c *Client) DeleteKnowledge(ctx context.Context, id string) error {
	path := fmt.Sprintf("knowledge/%s/delete", url.PathEscape(id))
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
