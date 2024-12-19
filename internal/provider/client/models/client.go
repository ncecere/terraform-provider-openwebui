package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Client implements the models operations
type Client struct {
	endpoint string
	token    string
}

// NewClient creates a new models client
func NewClient(endpoint, token string) *Client {
	return &Client{
		endpoint: endpoint,
		token:    token,
	}
}

func (c *Client) GetModel(id string) (*Model, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/models/model?id=%s", c.endpoint, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] GetModel response: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiModel APIModel
	if err := json.Unmarshal(bodyBytes, &apiModel); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return APIToModel(&apiModel), nil
}

func (c *Client) GetModels() ([]Model, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/models/", c.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] GetModels response: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiModels []APIModel
	if err := json.Unmarshal(bodyBytes, &apiModels); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var models []Model
	for _, apiModel := range apiModels {
		models = append(models, *APIToModel(&apiModel))
	}

	return models, nil
}

func (c *Client) CreateModel(model *Model) (*Model, error) {
	// Generate a new UUID for the model
	modelID := uuid.New().String()

	// Convert to API model
	apiModel := &APIModel{
		ID:          modelID,
		BaseModelID: model.BaseModelID.ValueString(),
		Name:        model.Name.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	// Handle Params
	if model.Params != nil {
		apiModel.Params = &APIModelParams{}
		if !model.Params.System.IsNull() {
			apiModel.Params.System = model.Params.System.ValueString()
		}
		if !model.Params.StreamResponse.IsNull() {
			apiModel.Params.StreamResponse = model.Params.StreamResponse.ValueBool()
		}
		if !model.Params.Temperature.IsNull() {
			apiModel.Params.Temperature = model.Params.Temperature.ValueFloat64()
		}
		if !model.Params.TopP.IsNull() {
			apiModel.Params.TopP = model.Params.TopP.ValueFloat64()
		}
		if !model.Params.MaxTokens.IsNull() {
			apiModel.Params.MaxTokens = model.Params.MaxTokens.ValueInt64()
		}
		if !model.Params.Seed.IsNull() {
			apiModel.Params.Seed = model.Params.Seed.ValueInt64()
		}
		if !model.Params.TopK.IsNull() {
			apiModel.Params.TopK = model.Params.TopK.ValueInt64()
		}
		if !model.Params.MinP.IsNull() {
			apiModel.Params.MinP = model.Params.MinP.ValueFloat64()
		}
		if !model.Params.FrequencyPenalty.IsNull() {
			apiModel.Params.FrequencyPenalty = model.Params.FrequencyPenalty.ValueInt64()
		}
		if !model.Params.RepeatLastN.IsNull() {
			apiModel.Params.RepeatLastN = model.Params.RepeatLastN.ValueInt64()
		}
		if !model.Params.NumCtx.IsNull() {
			apiModel.Params.NumCtx = model.Params.NumCtx.ValueInt64()
		}
		if !model.Params.NumBatch.IsNull() {
			apiModel.Params.NumBatch = model.Params.NumBatch.ValueInt64()
		}
		if !model.Params.NumKeep.IsNull() {
			apiModel.Params.NumKeep = model.Params.NumKeep.ValueInt64()
		}
	}

	// Handle Meta
	if model.Meta != nil {
		apiModel.Meta = &APIModelMeta{}
		if !model.Meta.ProfileImageURL.IsNull() {
			apiModel.Meta.ProfileImageURL = model.Meta.ProfileImageURL.ValueString()
		}
		if !model.Meta.Description.IsNull() {
			apiModel.Meta.Description = model.Meta.Description.ValueString()
		}

		if model.Meta.Capabilities != nil {
			apiModel.Meta.Capabilities = &APIModelCapabilities{
				Vision:    model.Meta.Capabilities.Vision.ValueBool(),
				Usage:     model.Meta.Capabilities.Usage.ValueBool(),
				Citations: model.Meta.Capabilities.Citations.ValueBool(),
			}
		}

		if len(model.Meta.Tags) > 0 {
			apiModel.Meta.Tags = make([]APITag, len(model.Meta.Tags))
			for i, tag := range model.Meta.Tags {
				if !tag.Name.IsNull() {
					apiModel.Meta.Tags[i] = APITag{
						Name: tag.Name.ValueString(),
					}
				}
			}
		}
	}

	// Handle AccessControl
	if model.AccessControl != nil {
		apiModel.AccessControl = &APIAccessControl{}
		if model.AccessControl.Read != nil {
			apiModel.AccessControl.Read = &APIAccessGroup{
				GroupIDs: make([]string, 0),
				UserIDs:  make([]string, 0),
			}
			for _, id := range model.AccessControl.Read.GroupIDs {
				if !id.IsNull() {
					apiModel.AccessControl.Read.GroupIDs = append(apiModel.AccessControl.Read.GroupIDs, id.ValueString())
				}
			}
			for _, id := range model.AccessControl.Read.UserIDs {
				if !id.IsNull() {
					apiModel.AccessControl.Read.UserIDs = append(apiModel.AccessControl.Read.UserIDs, id.ValueString())
				}
			}
		}
		if model.AccessControl.Write != nil {
			apiModel.AccessControl.Write = &APIAccessGroup{
				GroupIDs: make([]string, 0),
				UserIDs:  make([]string, 0),
			}
			for _, id := range model.AccessControl.Write.GroupIDs {
				if !id.IsNull() {
					apiModel.AccessControl.Write.GroupIDs = append(apiModel.AccessControl.Write.GroupIDs, id.ValueString())
				}
			}
			for _, id := range model.AccessControl.Write.UserIDs {
				if !id.IsNull() {
					apiModel.AccessControl.Write.UserIDs = append(apiModel.AccessControl.Write.UserIDs, id.ValueString())
				}
			}
		}
	}

	payload, err := json.Marshal(apiModel)
	if err != nil {
		return nil, fmt.Errorf("error marshaling model: %v", err)
	}

	log.Printf("[DEBUG] CreateModel request payload: %s", string(payload))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/models/create", c.endpoint), bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] CreateModel response: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var createdAPIModel APIModel
	if err := json.Unmarshal(bodyBytes, &createdAPIModel); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Ensure the ID is set in the response
	if createdAPIModel.ID == "" {
		createdAPIModel.ID = modelID
	}

	return APIToModel(&createdAPIModel), nil
}

func (c *Client) UpdateModel(id string, model *Model) (*Model, error) {
	// Convert to API model
	apiModel := &APIModel{
		ID:          id,
		BaseModelID: model.BaseModelID.ValueString(),
		Name:        model.Name.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	// Handle Params
	if model.Params != nil {
		apiModel.Params = &APIModelParams{}
		if !model.Params.System.IsNull() {
			apiModel.Params.System = model.Params.System.ValueString()
		}
		if !model.Params.StreamResponse.IsNull() {
			apiModel.Params.StreamResponse = model.Params.StreamResponse.ValueBool()
		}
		if !model.Params.Temperature.IsNull() {
			apiModel.Params.Temperature = model.Params.Temperature.ValueFloat64()
		}
		if !model.Params.TopP.IsNull() {
			apiModel.Params.TopP = model.Params.TopP.ValueFloat64()
		}
		if !model.Params.MaxTokens.IsNull() {
			apiModel.Params.MaxTokens = model.Params.MaxTokens.ValueInt64()
		}
		if !model.Params.Seed.IsNull() {
			apiModel.Params.Seed = model.Params.Seed.ValueInt64()
		}
		if !model.Params.TopK.IsNull() {
			apiModel.Params.TopK = model.Params.TopK.ValueInt64()
		}
		if !model.Params.MinP.IsNull() {
			apiModel.Params.MinP = model.Params.MinP.ValueFloat64()
		}
		if !model.Params.FrequencyPenalty.IsNull() {
			apiModel.Params.FrequencyPenalty = model.Params.FrequencyPenalty.ValueInt64()
		}
		if !model.Params.RepeatLastN.IsNull() {
			apiModel.Params.RepeatLastN = model.Params.RepeatLastN.ValueInt64()
		}
		if !model.Params.NumCtx.IsNull() {
			apiModel.Params.NumCtx = model.Params.NumCtx.ValueInt64()
		}
		if !model.Params.NumBatch.IsNull() {
			apiModel.Params.NumBatch = model.Params.NumBatch.ValueInt64()
		}
		if !model.Params.NumKeep.IsNull() {
			apiModel.Params.NumKeep = model.Params.NumKeep.ValueInt64()
		}
	}

	// Handle Meta
	if model.Meta != nil {
		apiModel.Meta = &APIModelMeta{}
		if !model.Meta.ProfileImageURL.IsNull() {
			apiModel.Meta.ProfileImageURL = model.Meta.ProfileImageURL.ValueString()
		}
		if !model.Meta.Description.IsNull() {
			apiModel.Meta.Description = model.Meta.Description.ValueString()
		}

		if model.Meta.Capabilities != nil {
			apiModel.Meta.Capabilities = &APIModelCapabilities{
				Vision:    model.Meta.Capabilities.Vision.ValueBool(),
				Usage:     model.Meta.Capabilities.Usage.ValueBool(),
				Citations: model.Meta.Capabilities.Citations.ValueBool(),
			}
		}

		if len(model.Meta.Tags) > 0 {
			apiModel.Meta.Tags = make([]APITag, len(model.Meta.Tags))
			for i, tag := range model.Meta.Tags {
				if !tag.Name.IsNull() {
					apiModel.Meta.Tags[i] = APITag{
						Name: tag.Name.ValueString(),
					}
				}
			}
		}
	}

	// Handle AccessControl
	if model.AccessControl != nil {
		apiModel.AccessControl = &APIAccessControl{}
		if model.AccessControl.Read != nil {
			apiModel.AccessControl.Read = &APIAccessGroup{
				GroupIDs: make([]string, 0),
				UserIDs:  make([]string, 0),
			}
			for _, id := range model.AccessControl.Read.GroupIDs {
				if !id.IsNull() {
					apiModel.AccessControl.Read.GroupIDs = append(apiModel.AccessControl.Read.GroupIDs, id.ValueString())
				}
			}
			for _, id := range model.AccessControl.Read.UserIDs {
				if !id.IsNull() {
					apiModel.AccessControl.Read.UserIDs = append(apiModel.AccessControl.Read.UserIDs, id.ValueString())
				}
			}
		}
		if model.AccessControl.Write != nil {
			apiModel.AccessControl.Write = &APIAccessGroup{
				GroupIDs: make([]string, 0),
				UserIDs:  make([]string, 0),
			}
			for _, id := range model.AccessControl.Write.GroupIDs {
				if !id.IsNull() {
					apiModel.AccessControl.Write.GroupIDs = append(apiModel.AccessControl.Write.GroupIDs, id.ValueString())
				}
			}
			for _, id := range model.AccessControl.Write.UserIDs {
				if !id.IsNull() {
					apiModel.AccessControl.Write.UserIDs = append(apiModel.AccessControl.Write.UserIDs, id.ValueString())
				}
			}
		}
	}

	payload, err := json.Marshal(apiModel)
	if err != nil {
		return nil, fmt.Errorf("error marshaling model: %v", err)
	}

	log.Printf("[DEBUG] UpdateModel request payload: %s", string(payload))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/models/model/update?id=%s", c.endpoint, id), bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] UpdateModel response: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var updatedAPIModel APIModel
	if err := json.Unmarshal(bodyBytes, &updatedAPIModel); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Ensure the ID is preserved
	if updatedAPIModel.ID == "" {
		updatedAPIModel.ID = id
	}

	return APIToModel(&updatedAPIModel), nil
}

func (c *Client) DeleteModel(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/models/model/delete?id=%s", c.endpoint, id), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] DeleteModel response: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
