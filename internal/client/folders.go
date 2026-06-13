package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type FolderForm struct {
	Name     string         `json:"name"`
	Data     map[string]any `json:"data,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
	ParentID *string        `json:"parent_id,omitempty"`
}

type FolderUpdateForm struct {
	Name *string        `json:"name,omitempty"`
	Data map[string]any `json:"data,omitempty"`
	Meta map[string]any `json:"meta,omitempty"`
}

type FolderParentIDForm struct {
	ParentID *string `json:"parent_id"`
}

type FolderIsExpandedForm struct {
	IsExpanded bool `json:"is_expanded"`
}

type FolderModel struct {
	ID         string         `json:"id"`
	ParentID   *string        `json:"parent_id"`
	UserID     string         `json:"user_id"`
	Name       string         `json:"name"`
	Items      map[string]any `json:"items"`
	Meta       map[string]any `json:"meta"`
	Data       map[string]any `json:"data"`
	IsExpanded bool           `json:"is_expanded"`
	CreatedAt  int64          `json:"created_at"`
	UpdatedAt  int64          `json:"updated_at"`
}

func (c *Client) CreateFolder(ctx context.Context, form FolderForm) (*FolderModel, error) {
	var resp FolderModel
	if err := c.do(ctx, http.MethodPost, "folders/", nil, form, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetFolder(ctx context.Context, id string) (*FolderModel, error) {
	var resp FolderModel
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("folders/%s", url.PathEscape(id)), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateFolder(ctx context.Context, id string, form FolderUpdateForm) (*FolderModel, error) {
	var resp FolderModel
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("folders/%s/update", url.PathEscape(id)), nil, form, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateFolderParent(ctx context.Context, id string, parentID *string) (*FolderModel, error) {
	var resp FolderModel
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("folders/%s/update/parent", url.PathEscape(id)), nil, FolderParentIDForm{ParentID: parentID}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateFolderExpanded(ctx context.Context, id string, expanded bool) (*FolderModel, error) {
	var resp FolderModel
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("folders/%s/update/expanded", url.PathEscape(id)), nil, FolderIsExpandedForm{IsExpanded: expanded}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteFolder(ctx context.Context, id string, deleteContents bool) error {
	query := url.Values{}
	query.Set("delete_contents", fmt.Sprintf("%t", deleteContents))
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("folders/%s", url.PathEscape(id)), query, nil, nil)
}
