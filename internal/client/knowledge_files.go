package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// KnowledgeFileIDForm associates a file with a knowledge base.
type KnowledgeFileIDForm struct {
	FileID      string  `json:"file_id"`
	DirectoryID *string `json:"directory_id,omitempty"`
}

// KnowledgeFileListResponse captures paginated file results for a knowledge base.
type KnowledgeFileListResponse struct {
	Items       []FileUserResponse        `json:"items"`
	Directories []KnowledgeDirectoryModel `json:"directories,omitempty"`
	Breadcrumbs []KnowledgeDirectoryModel `json:"breadcrumbs,omitempty"`
	Total       int                       `json:"total"`
}

// FileUserResponse captures file details including user metadata.
type FileUserResponse struct {
	ID        string         `json:"id"`
	UserID    string         `json:"user_id"`
	Hash      *string        `json:"hash,omitempty"`
	Filename  string         `json:"filename"`
	Data      map[string]any `json:"data,omitempty"`
	Meta      FileMeta       `json:"meta"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	User      *User          `json:"user,omitempty"`
}

// ListKnowledgeFiles retrieves file attachments for a knowledge base.
func (c *Client) ListKnowledgeFiles(ctx context.Context, knowledgeID string, queryValue string, viewOption string, orderBy string, direction string, page int) (*KnowledgeFileListResponse, error) {
	return c.ListKnowledgeFilesFiltered(ctx, knowledgeID, queryValue, false, viewOption, orderBy, direction, "", page)
}

// ListKnowledgeFilesFiltered retrieves file attachments for a knowledge base with all supported filters.
func (c *Client) ListKnowledgeFilesFiltered(ctx context.Context, knowledgeID string, queryValue string, includeContent bool, viewOption string, orderBy string, direction string, directoryID string, page int) (*KnowledgeFileListResponse, error) {
	query := url.Values{}
	if queryValue != "" {
		query.Set("query", queryValue)
	}
	if includeContent {
		query.Set("include_content", "true")
	}
	if viewOption != "" {
		query.Set("view_option", viewOption)
	}
	if orderBy != "" {
		query.Set("order_by", orderBy)
	}
	if direction != "" {
		query.Set("direction", direction)
	}
	if directoryID != "" {
		query.Set("directory_id", directoryID)
	}
	if page > 0 {
		query.Set("page", fmt.Sprintf("%d", page))
	}

	var resp KnowledgeFileListResponse
	path := fmt.Sprintf("knowledge/%s/files", url.PathEscape(knowledgeID))
	if err := c.do(ctx, http.MethodGet, path, query, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// AddKnowledgeFile associates a file with a knowledge base.
func (c *Client) AddKnowledgeFile(ctx context.Context, knowledgeID string, fileID string) (*KnowledgeFilesResponse, error) {
	return c.AddKnowledgeFileToDirectory(ctx, knowledgeID, fileID, nil)
}

// AddKnowledgeFileToDirectory associates a file with a knowledge base and optional directory.
func (c *Client) AddKnowledgeFileToDirectory(ctx context.Context, knowledgeID string, fileID string, directoryID *string) (*KnowledgeFilesResponse, error) {
	var resp KnowledgeFilesResponse
	form := KnowledgeFileIDForm{FileID: fileID, DirectoryID: directoryID}
	path := fmt.Sprintf("knowledge/%s/file/add", url.PathEscape(knowledgeID))
	if err := c.do(ctx, http.MethodPost, path, nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateKnowledgeFile updates a file association within a knowledge base.
func (c *Client) UpdateKnowledgeFile(ctx context.Context, knowledgeID string, fileID string) (*KnowledgeFilesResponse, error) {
	var resp KnowledgeFilesResponse
	form := KnowledgeFileIDForm{FileID: fileID}
	path := fmt.Sprintf("knowledge/%s/file/update", url.PathEscape(knowledgeID))
	if err := c.do(ctx, http.MethodPost, path, nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// MoveKnowledgeFile moves a file attachment to a directory.
func (c *Client) MoveKnowledgeFile(ctx context.Context, knowledgeID string, fileID string, directoryID *string) (*KnowledgeFilesResponse, error) {
	var resp KnowledgeFilesResponse
	form := KnowledgeFileIDForm{FileID: fileID, DirectoryID: directoryID}
	path := fmt.Sprintf("knowledge/%s/file/move", url.PathEscape(knowledgeID))
	if err := c.do(ctx, http.MethodPost, path, nil, form, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveKnowledgeFile detaches a file from a knowledge base.
func (c *Client) RemoveKnowledgeFile(ctx context.Context, knowledgeID string, fileID string, deleteFile bool) (*KnowledgeFilesResponse, error) {
	var resp KnowledgeFilesResponse
	form := KnowledgeFileIDForm{FileID: fileID}
	query := url.Values{"delete_file": []string{fmt.Sprintf("%t", deleteFile)}}
	path := fmt.Sprintf("knowledge/%s/file/remove", url.PathEscape(knowledgeID))
	if err := c.do(ctx, http.MethodPost, path, query, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
