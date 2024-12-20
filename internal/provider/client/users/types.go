package users

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// User represents the Terraform schema model for users
type User struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Email           types.String `tfsdk:"email"`
	Role            types.String `tfsdk:"role"`
	ProfileImageURL types.String `tfsdk:"profile_image_url"`
	LastActiveAt    types.Int64  `tfsdk:"last_active_at"`
	UpdatedAt       types.Int64  `tfsdk:"updated_at"`
	CreatedAt       types.Int64  `tfsdk:"created_at"`
	APIKey          types.String `tfsdk:"api_key"`
	Settings        *Settings    `tfsdk:"settings"`
	Info            types.Map    `tfsdk:"info"`
	OAuthSub        types.String `tfsdk:"oauth_sub"`
}

// Settings represents user settings
type Settings struct {
	UI types.Map `tfsdk:"ui"`
}

// APIUser represents the API response model for users
type APIUser struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Email           string                 `json:"email"`
	Role            string                 `json:"role"`
	ProfileImageURL string                 `json:"profile_image_url"`
	LastActiveAt    int64                  `json:"last_active_at"`
	UpdatedAt       int64                  `json:"updated_at"`
	CreatedAt       int64                  `json:"created_at"`
	APIKey          string                 `json:"api_key"`
	Settings        *APISettings           `json:"settings,omitempty"`
	Info            map[string]interface{} `json:"info,omitempty"`
	OAuthSub        string                 `json:"oauth_sub"`
}

// APISettings represents the API response model for user settings
type APISettings struct {
	UI map[string]interface{} `json:"ui,omitempty"`
}

// APIToUser converts an API user to a Terraform user
func APIToUser(apiUser *APIUser) *User {
	user := &User{
		ID:              types.StringValue(apiUser.ID),
		Name:            types.StringValue(apiUser.Name),
		Email:           types.StringValue(apiUser.Email),
		Role:            types.StringValue(apiUser.Role),
		ProfileImageURL: types.StringValue(apiUser.ProfileImageURL),
		LastActiveAt:    types.Int64Value(apiUser.LastActiveAt),
		UpdatedAt:       types.Int64Value(apiUser.UpdatedAt),
		CreatedAt:       types.Int64Value(apiUser.CreatedAt),
		APIKey:          types.StringValue(apiUser.APIKey),
		OAuthSub:        types.StringValue(apiUser.OAuthSub),
	}

	if apiUser.Settings != nil {
		user.Settings = &Settings{}
		if apiUser.Settings.UI != nil {
			ui, _ := types.MapValueFrom(nil, types.StringType, apiUser.Settings.UI)
			user.Settings.UI = ui
		}
	}

	if apiUser.Info != nil {
		info, _ := types.MapValueFrom(nil, types.StringType, apiUser.Info)
		user.Info = info
	}

	return user
}
