package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model represents the Terraform schema model
type Model struct {
	ID            types.String   `tfsdk:"id"`
	UserID        types.String   `tfsdk:"user_id"`
	BaseModelID   types.String   `tfsdk:"base_model_id"`
	Name          types.String   `tfsdk:"name"`
	Params        *ModelParams   `tfsdk:"params"`
	Meta          *ModelMeta     `tfsdk:"meta"`
	AccessControl *AccessControl `tfsdk:"access_control"`
	IsActive      types.Bool     `tfsdk:"is_active"`
	UpdatedAt     types.Int64    `tfsdk:"updated_at"`
	CreatedAt     types.Int64    `tfsdk:"created_at"`
}

// APIModel represents the API response model
type APIModel struct {
	ID            string            `json:"id"`
	UserID        string            `json:"user_id"`
	BaseModelID   string            `json:"base_model_id"`
	Name          string            `json:"name"`
	Params        *APIModelParams   `json:"params,omitempty"`
	Meta          *APIModelMeta     `json:"meta,omitempty"`
	AccessControl *APIAccessControl `json:"access_control,omitempty"`
	IsActive      bool              `json:"is_active"`
	UpdatedAt     int64             `json:"updated_at"`
	CreatedAt     int64             `json:"created_at"`
}

type ModelParams struct {
	System           types.String  `tfsdk:"system"`
	StreamResponse   types.Bool    `tfsdk:"stream_response"`
	Seed             types.Int64   `tfsdk:"seed"`
	Temperature      types.Float64 `tfsdk:"temperature"`
	TopK             types.Int64   `tfsdk:"top_k"`
	TopP             types.Float64 `tfsdk:"top_p"`
	MinP             types.Float64 `tfsdk:"min_p"`
	FrequencyPenalty types.Int64   `tfsdk:"frequency_penalty"`
	RepeatLastN      types.Int64   `tfsdk:"repeat_last_n"`
	NumCtx           types.Int64   `tfsdk:"num_ctx"`
	NumBatch         types.Int64   `tfsdk:"num_batch"`
	NumKeep          types.Int64   `tfsdk:"num_keep"`
	MaxTokens        types.Int64   `tfsdk:"max_tokens"`
}

type APIModelParams struct {
	System           string  `json:"system,omitempty"`
	StreamResponse   bool    `json:"stream_response,omitempty"`
	Seed             int64   `json:"seed,omitempty"`
	Temperature      float64 `json:"temperature,omitempty"`
	TopK             int64   `json:"top_k,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	MinP             float64 `json:"min_p,omitempty"`
	FrequencyPenalty int64   `json:"frequency_penalty,omitempty"`
	RepeatLastN      int64   `json:"repeat_last_n,omitempty"`
	NumCtx           int64   `json:"num_ctx,omitempty"`
	NumBatch         int64   `json:"num_batch,omitempty"`
	NumKeep          int64   `json:"num_keep,omitempty"`
	MaxTokens        int64   `json:"max_tokens,omitempty"`
}

type ModelMeta struct {
	ProfileImageURL types.String       `tfsdk:"profile_image_url"`
	Description     types.String       `tfsdk:"description"`
	Capabilities    *ModelCapabilities `tfsdk:"capabilities"`
	Tags            []Tag              `tfsdk:"tags"`
}

type APIModelMeta struct {
	ProfileImageURL string                `json:"profile_image_url,omitempty"`
	Description     string                `json:"description,omitempty"`
	Capabilities    *APIModelCapabilities `json:"capabilities,omitempty"`
	Tags            []APITag              `json:"tags,omitempty"`
}

type ModelCapabilities struct {
	Vision    types.Bool `tfsdk:"vision"`
	Usage     types.Bool `tfsdk:"usage"`
	Citations types.Bool `tfsdk:"citations"`
}

type APIModelCapabilities struct {
	Vision    bool `json:"vision,omitempty"`
	Usage     bool `json:"usage,omitempty"`
	Citations bool `json:"citations,omitempty"`
}

type Tag struct {
	Name types.String `tfsdk:"name"`
}

type APITag struct {
	Name string `json:"name,omitempty"`
}

type AccessControl struct {
	Read  *AccessGroup `tfsdk:"read"`
	Write *AccessGroup `tfsdk:"write"`
}

type APIAccessControl struct {
	Read  *APIAccessGroup `json:"read,omitempty"`
	Write *APIAccessGroup `json:"write,omitempty"`
}

type AccessGroup struct {
	GroupIDs []types.String `tfsdk:"group_ids"`
	UserIDs  []types.String `tfsdk:"user_ids"`
}

type APIAccessGroup struct {
	GroupIDs []string `json:"group_ids,omitempty"`
	UserIDs  []string `json:"user_ids,omitempty"`
}

// Helper function to convert API model to Terraform model
func APIToModel(apiModel *APIModel) *Model {
	model := &Model{
		ID:          types.StringValue(apiModel.ID),
		UserID:      types.StringValue(apiModel.UserID),
		BaseModelID: types.StringValue(apiModel.BaseModelID),
		Name:        types.StringValue(apiModel.Name),
		IsActive:    types.BoolValue(apiModel.IsActive),
		UpdatedAt:   types.Int64Value(apiModel.UpdatedAt),
		CreatedAt:   types.Int64Value(apiModel.CreatedAt),
	}

	if apiModel.Params != nil {
		model.Params = &ModelParams{}
		if apiModel.Params.System != "" {
			model.Params.System = types.StringValue(apiModel.Params.System)
		}
		model.Params.StreamResponse = types.BoolValue(apiModel.Params.StreamResponse)
		if apiModel.Params.Temperature != 0 {
			model.Params.Temperature = types.Float64Value(apiModel.Params.Temperature)
		}
		if apiModel.Params.TopP != 0 {
			model.Params.TopP = types.Float64Value(apiModel.Params.TopP)
		}
		if apiModel.Params.MaxTokens != 0 {
			model.Params.MaxTokens = types.Int64Value(apiModel.Params.MaxTokens)
		}
		if apiModel.Params.Seed != 0 {
			model.Params.Seed = types.Int64Value(apiModel.Params.Seed)
		}
		if apiModel.Params.TopK != 0 {
			model.Params.TopK = types.Int64Value(apiModel.Params.TopK)
		}
		if apiModel.Params.MinP != 0 {
			model.Params.MinP = types.Float64Value(apiModel.Params.MinP)
		}
		if apiModel.Params.FrequencyPenalty != 0 {
			model.Params.FrequencyPenalty = types.Int64Value(apiModel.Params.FrequencyPenalty)
		}
		if apiModel.Params.RepeatLastN != 0 {
			model.Params.RepeatLastN = types.Int64Value(apiModel.Params.RepeatLastN)
		}
		if apiModel.Params.NumCtx != 0 {
			model.Params.NumCtx = types.Int64Value(apiModel.Params.NumCtx)
		}
		if apiModel.Params.NumBatch != 0 {
			model.Params.NumBatch = types.Int64Value(apiModel.Params.NumBatch)
		}
		if apiModel.Params.NumKeep != 0 {
			model.Params.NumKeep = types.Int64Value(apiModel.Params.NumKeep)
		}
	}

	if apiModel.Meta != nil {
		model.Meta = &ModelMeta{}
		if apiModel.Meta.ProfileImageURL != "" {
			model.Meta.ProfileImageURL = types.StringValue(apiModel.Meta.ProfileImageURL)
		}
		if apiModel.Meta.Description != "" {
			model.Meta.Description = types.StringValue(apiModel.Meta.Description)
		}

		if apiModel.Meta.Capabilities != nil {
			model.Meta.Capabilities = &ModelCapabilities{
				Vision:    types.BoolValue(apiModel.Meta.Capabilities.Vision),
				Usage:     types.BoolValue(apiModel.Meta.Capabilities.Usage),
				Citations: types.BoolValue(apiModel.Meta.Capabilities.Citations),
			}
		}

		if len(apiModel.Meta.Tags) > 0 {
			model.Meta.Tags = make([]Tag, len(apiModel.Meta.Tags))
			for i, tag := range apiModel.Meta.Tags {
				model.Meta.Tags[i] = Tag{
					Name: types.StringValue(tag.Name),
				}
			}
		}
	}

	if apiModel.AccessControl != nil {
		model.AccessControl = &AccessControl{}
		if apiModel.AccessControl.Read != nil {
			model.AccessControl.Read = &AccessGroup{
				GroupIDs: make([]types.String, len(apiModel.AccessControl.Read.GroupIDs)),
				UserIDs:  make([]types.String, len(apiModel.AccessControl.Read.UserIDs)),
			}
			for i, id := range apiModel.AccessControl.Read.GroupIDs {
				model.AccessControl.Read.GroupIDs[i] = types.StringValue(id)
			}
			for i, id := range apiModel.AccessControl.Read.UserIDs {
				model.AccessControl.Read.UserIDs[i] = types.StringValue(id)
			}
		}
		if apiModel.AccessControl.Write != nil {
			model.AccessControl.Write = &AccessGroup{
				GroupIDs: make([]types.String, len(apiModel.AccessControl.Write.GroupIDs)),
				UserIDs:  make([]types.String, len(apiModel.AccessControl.Write.UserIDs)),
			}
			for i, id := range apiModel.AccessControl.Write.GroupIDs {
				model.AccessControl.Write.GroupIDs[i] = types.StringValue(id)
			}
			for i, id := range apiModel.AccessControl.Write.UserIDs {
				model.AccessControl.Write.UserIDs[i] = types.StringValue(id)
			}
		}
	}

	return model
}
