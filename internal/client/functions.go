package client

import (
	"context"
	"net/http"
	"net/url"
)

type FunctionMeta struct {
	Description *string        `json:"description,omitempty"`
	Manifest    map[string]any `json:"manifest,omitempty"`
}

type FunctionForm struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Content string       `json:"content"`
	Meta    FunctionMeta `json:"meta"`
}

type FunctionModel struct {
	ID        string       `json:"id"`
	UserID    string       `json:"user_id"`
	Name      string       `json:"name"`
	Type      string       `json:"type"`
	Content   string       `json:"content"`
	Meta      FunctionMeta `json:"meta"`
	IsActive  bool         `json:"is_active"`
	IsGlobal  bool         `json:"is_global"`
	UpdatedAt int64        `json:"updated_at"`
	CreatedAt int64        `json:"created_at"`
}

func (c *Client) CreateFunction(ctx context.Context, form FunctionForm) (*FunctionModel, error) {
	var resp FunctionModel
	if err := c.do(ctx, http.MethodPost, "functions/create", nil, form, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (c *Client) GetFunction(ctx context.Context, id string) (*FunctionModel, error) {
	var resp FunctionModel
	if err := c.do(ctx, http.MethodGet, "functions/id/"+url.PathEscape(id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (c *Client) UpdateFunction(ctx context.Context, id string, form FunctionForm) (*FunctionModel, error) {
	var resp FunctionModel
	if err := c.do(ctx, http.MethodPost, "functions/id/"+url.PathEscape(id)+"/update", nil, form, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (c *Client) DeleteFunction(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "functions/id/"+url.PathEscape(id)+"/delete", nil, nil, nil)
}
func (c *Client) ToggleFunction(ctx context.Context, id string) (*FunctionModel, error) {
	var resp FunctionModel
	if err := c.do(ctx, http.MethodPost, "functions/id/"+url.PathEscape(id)+"/toggle", nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (c *Client) ToggleFunctionGlobal(ctx context.Context, id string) (*FunctionModel, error) {
	var resp FunctionModel
	if err := c.do(ctx, http.MethodPost, "functions/id/"+url.PathEscape(id)+"/toggle/global", nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (c *Client) ListFunctionsRaw(ctx context.Context) (any, error) {
	var resp any
	if err := c.do(ctx, http.MethodGet, "functions/list", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *Client) GetFunctionValves(ctx context.Context, id string) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodGet, "functions/id/"+url.PathEscape(id)+"/valves", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *Client) GetFunctionValvesSpec(ctx context.Context, id string) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodGet, "functions/id/"+url.PathEscape(id)+"/valves/spec", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *Client) UpdateFunctionValves(ctx context.Context, id string, valves map[string]any) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodPost, "functions/id/"+url.PathEscape(id)+"/valves/update", nil, valves, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
