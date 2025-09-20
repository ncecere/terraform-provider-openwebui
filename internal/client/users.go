package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// User represents an Open WebUI user account.
type User struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Username     *string `json:"username"`
	Role         string  `json:"role"`
	ProfileImage string  `json:"profile_image_url"`
	Bio          *string `json:"bio"`
	OAuthSubject *string `json:"oauth_sub"`
	LastActiveAt int64   `json:"last_active_at"`
	UpdatedAt    int64   `json:"updated_at"`
	CreatedAt    int64   `json:"created_at"`
}

// listUsersResponse models the API contract for GET /users/.
type listUsersResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

// SearchUsers finds users whose username, email, or name matches the provided query.
func (c *Client) SearchUsers(ctx context.Context, query string, limit int) ([]User, error) {
	values := url.Values{}
	if query != "" {
		values.Set("query", query)
	}
	if limit > 0 {
		values.Set("limit", fmt.Sprintf("%d", limit))
	}

	var resp listUsersResponse
	if err := c.do(ctx, http.MethodGet, "users/", values, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Users, nil
}

// GetUser retrieves a user by identifier.
func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {
	var resp User
	path := fmt.Sprintf("users/%s", url.PathEscape(id))
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
