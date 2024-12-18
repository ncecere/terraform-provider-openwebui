package client

import (
	"fmt"

	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/groups"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/knowledge"
)

// OpenWebUIClient implements the Client interface
type OpenWebUIClient struct {
	BaseClient
	Knowledge *knowledge.Client
	Groups    *groups.Client
}

// New creates a new OpenWebUI client
func New(endpoint, token string) (*OpenWebUIClient, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	base := NewBaseClient(endpoint, token)
	return &OpenWebUIClient{
		BaseClient: base,
		Knowledge:  knowledge.NewClient(endpoint, token),
		Groups:     groups.NewClient(endpoint, token),
	}, nil
}

// Create implements KnowledgeClient
func (c *OpenWebUIClient) Create(form *knowledge.KnowledgeForm) (*knowledge.KnowledgeResponse, error) {
	return c.Knowledge.Create(form)
}

// Get implements KnowledgeClient
func (c *OpenWebUIClient) Get(id string) (*knowledge.KnowledgeResponse, error) {
	return c.Knowledge.Get(id)
}

// List implements KnowledgeClient
func (c *OpenWebUIClient) List() ([]knowledge.KnowledgeResponse, error) {
	return c.Knowledge.List()
}

// Update implements KnowledgeClient
func (c *OpenWebUIClient) Update(id string, form *knowledge.KnowledgeForm) (*knowledge.KnowledgeResponse, error) {
	return c.Knowledge.Update(id, form)
}

// Delete implements KnowledgeClient
func (c *OpenWebUIClient) Delete(id string) error {
	return c.Knowledge.Delete(id)
}
