package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &groupDataSource{}
var _ datasource.DataSourceWithConfigure = &groupDataSource{}

// groupDataSource exposes read-only group information.
type groupDataSource struct {
	client *client.Client
}

// groupDataSourceModel embeds the resource representation and adds the lookup identifier.
type groupDataSourceModel struct {
	GroupID types.String `tfsdk:"group_id"`
	groupResourceModel
}

// NewGroupDataSource constructs a new group data source.
func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

// Metadata sets the data source type name.
func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema describes the group data source schema.
func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the group to retrieve. If omitted, `name` must be provided.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier assigned by Open WebUI.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Group name used to locate the record when `group_id` is not supplied.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Group description.",
			},
			"users": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "User identifiers (or usernames/emails) that belong to the group.",
			},
			"permissions": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Fine-grained permission flags organised by category.",
				Attributes: map[string]schema.Attribute{
					"workspace": schema.MapAttribute{
						ElementType: types.BoolType,
						Computed:    true,
					},
					"sharing": schema.MapAttribute{
						ElementType: types.BoolType,
						Computed:    true,
					},
					"chat": schema.MapAttribute{
						ElementType: types.BoolType,
						Computed:    true,
					},
					"features": schema.MapAttribute{
						ElementType: types.BoolType,
						Computed:    true,
					},
				},
			},
			"meta_json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON metadata associated with the group.",
			},
			"data_json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON payload containing additional group data.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the user who owns the group.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Creation date assigned by Open WebUI (YYYY-MM-DD).",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Last update date assigned by Open WebUI (YYYY-MM-DD).",
			},
		},
	}
}

// Configure attaches the API client.
func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

// Read retrieves the group details.
func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the group data source.")
		return
	}

	var config groupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := ""
	if !config.GroupID.IsNull() && !config.GroupID.IsUnknown() {
		id = strings.TrimSpace(config.GroupID.ValueString())
	}

	name := ""
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		name = strings.TrimSpace(config.Name.ValueString())
	}

	if id == "" {
		if name == "" {
			resp.Diagnostics.AddError(
				"Missing group lookup value",
				"Either group_id or name must be provided to query an existing group.",
			)
			return
		}

		groups, err := d.client.ListGroups(ctx)
		if err != nil {
			resp.Diagnostics.AddError("List groups failed", err.Error())
			return
		}

		var matches []client.GroupResponse
		for _, g := range groups {
			if strings.EqualFold(g.Name, name) {
				matches = append(matches, g)
			}
		}

		if len(matches) == 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("name"),
				"Group not found",
				fmt.Sprintf("No Open WebUI group was found with the name %q.", name),
			)
			return
		}
		if len(matches) > 1 {
			resp.Diagnostics.AddAttributeError(
				path.Root("name"),
				"Group name not unique",
				fmt.Sprintf("Multiple groups share the name %q. Provide group_id instead.", name),
			)
			return
		}

		id = matches[0].ID
	}

	current, err := d.client.GetGroup(ctx, id)
	if err != nil {
		if err == client.ErrNotFound {
			resp.Diagnostics.AddAttributeError(
				path.Root("group_id"),
				"Group not found",
				"No Open WebUI group was found with the supplied group_id.",
			)
			return
		}

		resp.Diagnostics.AddError("Read group failed", err.Error())
		return
	}

	model, diags := groupResponseToModel(ctx, d.client, current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := groupDataSourceModel{
		GroupID:            types.StringValue(current.ID),
		groupResourceModel: model,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
