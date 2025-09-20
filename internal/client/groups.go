package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// GroupForm represents the payload for creating a group.
type GroupForm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GroupUpdateForm represents the payload for updating an existing group.
type GroupUpdateForm struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Permissions map[string]any `json:"permissions,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
	Data        map[string]any `json:"data,omitempty"`
}

// GroupResponse captures group details returned by the API.
type GroupResponse struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
	UserIDs     []string       `json:"user_ids"`
	AdminIDs    []string       `json:"admin_ids,omitempty"`
	Permissions map[string]any `json:"permissions,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
	Data        map[string]any `json:"data,omitempty"`
}

// CreateGroup provisions a new group.
func (c *Client) CreateGroup(ctx context.Context, form GroupForm) (*GroupResponse, error) {
	var resp GroupResponse
	if err := c.do(ctx, http.MethodPost, "groups/create", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListGroups retrieves all groups.
func (c *Client) ListGroups(ctx context.Context) ([]GroupResponse, error) {
	var resp []GroupResponse
	if err := c.do(ctx, http.MethodGet, "groups/", nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetGroup retrieves a group by identifier.
func (c *Client) GetGroup(ctx context.Context, id string) (*GroupResponse, error) {
	var resp GroupResponse
	path := fmt.Sprintf("groups/id/%s", url.PathEscape(id))
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateGroup updates fields on an existing group.
func (c *Client) UpdateGroup(ctx context.Context, id string, form GroupUpdateForm) (*GroupResponse, error) {
	var resp GroupResponse
	path := fmt.Sprintf("groups/id/%s/update", url.PathEscape(id))
	if err := c.do(ctx, http.MethodPost, path, nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// AddGroupUsers adds users to a group by their IDs.
func (c *Client) AddGroupUsers(ctx context.Context, id string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	body := map[string]any{
		"user_ids": userIDs,
	}

	path := fmt.Sprintf("groups/id/%s/users/add", url.PathEscape(id))
	return c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// RemoveGroupUsers removes users from a group by their IDs.
func (c *Client) RemoveGroupUsers(ctx context.Context, id string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	body := map[string]any{
		"user_ids": userIDs,
	}

	path := fmt.Sprintf("groups/id/%s/users/remove", url.PathEscape(id))
	return c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// DeleteGroup removes a group.
func (c *Client) DeleteGroup(ctx context.Context, id string) error {
	path := fmt.Sprintf("groups/id/%s/delete", url.PathEscape(id))
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
