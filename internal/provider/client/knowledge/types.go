package knowledge

import (
	"encoding/json"
)

// KnowledgeClient defines the interface for knowledge operations
type KnowledgeClient interface {
	Create(form *KnowledgeForm) (*KnowledgeResponse, error)
	Get(id string) (*KnowledgeResponse, error)
	List() ([]KnowledgeResponse, error)
	Update(id string, form *KnowledgeForm) (*KnowledgeResponse, error)
	Delete(id string) error
}

// KnowledgeForm represents the form data for creating/updating a knowledge base
type KnowledgeForm struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Data          map[string]string      `json:"data,omitempty"`
	AccessControl map[string]interface{} `json:"access_control,omitempty"`
}

// KnowledgeResponse represents the API response for a knowledge base
type KnowledgeResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Data          map[string]interface{} `json:"data,omitempty"`
	AccessControl interface{}            `json:"access_control,omitempty"`
	UpdatedAt     int64                  `json:"updated_at"`
	CreatedAt     int64                  `json:"created_at"`
}

// UnmarshalJSON implements custom JSON unmarshaling for KnowledgeResponse
func (k *KnowledgeResponse) UnmarshalJSON(data []byte) error {
	type Alias KnowledgeResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(k),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
