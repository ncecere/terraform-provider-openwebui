package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/users"
)

var (
	_ datasource.DataSource = &UserDataSource{}
)

type UserDataSourceModel struct {
	ID              types.String    `tfsdk:"id"`
	Name            types.String    `tfsdk:"name"`
	Email           types.String    `tfsdk:"email"`
	Role            types.String    `tfsdk:"role"`
	ProfileImageURL types.String    `tfsdk:"profile_image_url"`
	LastActiveAt    types.Int64     `tfsdk:"last_active_at"`
	UpdatedAt       types.Int64     `tfsdk:"updated_at"`
	CreatedAt       types.Int64     `tfsdk:"created_at"`
	APIKey          types.String    `tfsdk:"api_key"`
	Settings        *users.Settings `tfsdk:"settings"`
	Info            types.Map       `tfsdk:"info"`
	OAuthSub        types.String    `tfsdk:"oauth_sub"`
}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

type UserDataSource struct {
	client *users.Client
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a user by ID, email, or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the user.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the user.",
				Optional:    true,
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email address of the user.",
				Optional:    true,
				Computed:    true,
			},
			"role": schema.StringAttribute{
				Description: "The role of the user (pending, admin, or user).",
				Computed:    true,
			},
			"profile_image_url": schema.StringAttribute{
				Description: "URL of the user's profile image.",
				Computed:    true,
			},
			"last_active_at": schema.Int64Attribute{
				Description: "Timestamp of the user's last activity.",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Timestamp when the user was last updated.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Timestamp when the user was created.",
				Computed:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The API key of the user.",
				Computed:    true,
				Sensitive:   true,
			},
			"settings": schema.SingleNestedAttribute{
				Description: "User settings.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"ui": schema.MapAttribute{
						Description: "UI-specific settings.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"info": schema.MapAttribute{
				Description: "Additional user information.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"oauth_sub": schema.StringAttribute{
				Description: "OAuth subject identifier.",
				Computed:    true,
			},
		},
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clients, ok := req.ProviderData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected map[string]interface{}, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	client, ok := clients["users"].(*users.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *users.Client, got: %T. Please report this issue to the provider developers.", clients["users"]),
		)
		return
	}

	d.client = client
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that only one of id, email, or name is specified
	specifiedFields := 0
	if !config.ID.IsNull() {
		specifiedFields++
	}
	if !config.Email.IsNull() {
		specifiedFields++
	}
	if !config.Name.IsNull() {
		specifiedFields++
	}

	if specifiedFields == 0 {
		resp.Diagnostics.AddError(
			"Missing search criteria",
			"One of id, email, or name must be provided",
		)
		return
	}

	if specifiedFields > 1 {
		resp.Diagnostics.AddError(
			"Multiple search criteria",
			"Only one of id, email, or name can be provided",
		)
		return
	}

	var user *users.User
	var err error

	// Try to find user by ID first
	if !config.ID.IsNull() {
		user, err = d.client.GetUser(config.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading user by ID",
				fmt.Sprintf("Could not read user ID %s: %s", config.ID.ValueString(), err.Error()),
			)
			return
		}
	} else if !config.Email.IsNull() {
		// Try to find user by email
		user, err = d.client.FindUserByEmail(config.Email.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading user by email",
				fmt.Sprintf("Could not find user with email %s: %s", config.Email.ValueString(), err.Error()),
			)
			return
		}
	} else if !config.Name.IsNull() {
		// Try to find user by name
		user, err = d.client.FindUserByName(config.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading user by name",
				fmt.Sprintf("Could not find user with name %s: %s", config.Name.ValueString(), err.Error()),
			)
			return
		}
	}

	// Convert User to UserDataSourceModel
	var state UserDataSourceModel
	state.ID = user.ID
	state.Name = user.Name
	state.Email = user.Email
	state.Role = user.Role
	state.ProfileImageURL = user.ProfileImageURL
	state.LastActiveAt = user.LastActiveAt
	state.UpdatedAt = user.UpdatedAt
	state.CreatedAt = user.CreatedAt
	state.APIKey = user.APIKey
	state.Settings = user.Settings

	// Convert info map to types.Map with string values
	if user.Info.IsNull() {
		state.Info = types.MapNull(types.StringType)
	} else {
		infoMap := make(map[string]string)
		for k, v := range user.Info.Elements() {
			if str, ok := v.(types.String); ok {
				infoMap[k] = str.ValueString()
			}
		}
		state.Info, _ = types.MapValueFrom(ctx, types.StringType, infoMap)
	}

	state.OAuthSub = user.OAuthSub

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
