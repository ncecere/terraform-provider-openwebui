package client

import (
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/knowledge"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/models"
)

// Client interface defines all client operations
type Client interface {
	knowledge.KnowledgeClient
	ModelsClient
}

// ModelsClient defines model operations
type ModelsClient interface {
	GetModel(id string) (*models.Model, error)
	GetModels() ([]models.Model, error)
	CreateModel(model *models.Model) (*models.Model, error)
	UpdateModel(id string, model *models.Model) (*models.Model, error)
	DeleteModel(id string) error
}

// BaseClient provides common functionality for all clients
type BaseClient struct {
	endpoint string
	token    string
}

// NewBaseClient creates a new base client
func NewBaseClient(endpoint, token string) BaseClient {
	return BaseClient{
		endpoint: endpoint,
		token:    token,
	}
}

// GetEndpoint returns the API endpoint
func (c *BaseClient) GetEndpoint() string {
	return c.endpoint
}

// GetToken returns the API token
func (c *BaseClient) GetToken() string {
	return c.token
}
