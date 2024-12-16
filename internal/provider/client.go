package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is the API client for OpenWebUI
type Client struct {
	endpoint string
	token    string
	http     *http.Client
}

// NewClient creates a new OpenWebUI API client
func NewClient(endpoint string, token string) *Client {
	return &Client{
		endpoint: endpoint,
		token:    token,
		http:     &http.Client{},
	}
}

// KnowledgeForm represents the request body for creating/updating a knowledge base
type KnowledgeForm struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Data          map[string]string      `json:"data,omitempty"`
	AccessControl map[string]interface{} `json:"access_control,omitempty"`
}

// KnowledgeResponse represents the response from the knowledge API
type KnowledgeResponse struct {
	ID            string                 `json:"id"`
	UserID        string                 `json:"user_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Data          map[string]string      `json:"data,omitempty"`
	Meta          map[string]string      `json:"meta,omitempty"`
	AccessControl map[string]interface{} `json:"access_control,omitempty"`
	CreatedAt     int64                  `json:"created_at"`
	UpdatedAt     int64                  `json:"updated_at"`
}

// CreateKnowledge creates a new knowledge base
func (c *Client) CreateKnowledge(req *KnowledgeForm) (*KnowledgeResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/knowledge/create", c.endpoint), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", response.StatusCode)
	}

	var result KnowledgeResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// GetKnowledge retrieves a knowledge base by ID
func (c *Client) GetKnowledge(id string) (*KnowledgeResponse, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/knowledge/%s", c.endpoint, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", response.StatusCode)
	}

	var result KnowledgeResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// UpdateKnowledge updates an existing knowledge base
func (c *Client) UpdateKnowledge(id string, req *KnowledgeForm) (*KnowledgeResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/knowledge/%s/update", c.endpoint, id), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", response.StatusCode)
	}

	var result KnowledgeResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// DeleteKnowledge deletes a knowledge base
func (c *Client) DeleteKnowledge(id string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/knowledge/%s/delete", c.endpoint, id), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	response, err := c.http.Do(request)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", response.StatusCode)
	}

	return nil
}
