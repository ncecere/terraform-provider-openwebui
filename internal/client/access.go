package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

// AccessGrant represents a single access grant returned by the API.
type AccessGrant struct {
	ID            string `json:"id,omitempty"`
	ResourceType  string `json:"resource_type,omitempty"`
	ResourceID    string `json:"resource_id,omitempty"`
	PrincipalType string `json:"principal_type"`
	PrincipalID   string `json:"principal_id"`
	Permission    string `json:"permission"`
	CreatedAt     int64  `json:"created_at,omitempty"`
}

// BuildGroupGrants converts read/write group ID lists into a flat slice of
// AccessGrant objects. Write implies read.
func BuildGroupGrants(readIDs, writeIDs []string) []AccessGrant {
	if len(readIDs) == 0 && len(writeIDs) == 0 {
		return []AccessGrant{}
	}

	dedupe := func(ids []string) []string {
		seen := make(map[string]struct{}, len(ids))
		out := make([]string, 0, len(ids))
		for _, id := range ids {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			out = append(out, id)
		}
		sort.Strings(out)
		return out
	}

	writeOrdered := dedupe(writeIDs)
	readOrdered := dedupe(append(append([]string{}, readIDs...), writeIDs...))

	grants := make([]AccessGrant, 0, len(readOrdered)+len(writeOrdered))
	for _, id := range readOrdered {
		grants = append(grants, AccessGrant{
			PrincipalType: "group",
			PrincipalID:   id,
			Permission:    "read",
		})
	}
	for _, id := range writeOrdered {
		grants = append(grants, AccessGrant{
			PrincipalType: "group",
			PrincipalID:   id,
			Permission:    "write",
		})
	}

	return grants
}

// ExtractGroupIDs returns unique group principal IDs holding the named permission.
func ExtractGroupIDs(grants []AccessGrant, permission string) []string {
	if len(grants) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	var ids []string
	for _, g := range grants {
		if g.PrincipalType != "group" {
			continue
		}
		if g.Permission != permission {
			continue
		}
		if _, ok := seen[g.PrincipalID]; ok {
			continue
		}
		seen[g.PrincipalID] = struct{}{}
		ids = append(ids, g.PrincipalID)
	}

	return ids
}

type accessUpdateResponse struct {
	AccessGrants []AccessGrant `json:"access_grants"`
}

// UpdateModelAccess replaces the access grants on a model and returns the new grants.
func (c *Client) UpdateModelAccess(ctx context.Context, id string, grants []AccessGrant) ([]AccessGrant, error) {
	body := map[string]any{
		"id":            id,
		"access_grants": grants,
	}
	var resp accessUpdateResponse
	if err := c.do(ctx, http.MethodPost, "models/model/access/update", nil, body, &resp); err != nil {
		return nil, err
	}
	return resp.AccessGrants, nil
}

// UpdateKnowledgeAccess replaces the access grants on a knowledge entry.
func (c *Client) UpdateKnowledgeAccess(ctx context.Context, id string, grants []AccessGrant) ([]AccessGrant, error) {
	body := map[string]any{"access_grants": grants}
	path := fmt.Sprintf("knowledge/%s/access/update", url.PathEscape(id))
	var resp accessUpdateResponse
	if err := c.do(ctx, http.MethodPost, path, nil, body, &resp); err != nil {
		return nil, err
	}
	return resp.AccessGrants, nil
}

// UpdatePromptAccess replaces the access grants on a prompt (path is by prompt id).
func (c *Client) UpdatePromptAccess(ctx context.Context, promptID string, grants []AccessGrant) ([]AccessGrant, error) {
	body := map[string]any{"access_grants": grants}
	path := fmt.Sprintf("prompts/id/%s/access/update", url.PathEscape(promptID))
	var resp accessUpdateResponse
	if err := c.do(ctx, http.MethodPost, path, nil, body, &resp); err != nil {
		return nil, err
	}
	return resp.AccessGrants, nil
}

// UpdateToolAccess replaces the access grants on a tool.
func (c *Client) UpdateToolAccess(ctx context.Context, id string, grants []AccessGrant) ([]AccessGrant, error) {
	body := map[string]any{"access_grants": grants}
	path := fmt.Sprintf("tools/id/%s/access/update", url.PathEscape(id))
	var resp accessUpdateResponse
	if err := c.do(ctx, http.MethodPost, path, nil, body, &resp); err != nil {
		return nil, err
	}
	return resp.AccessGrants, nil
}
