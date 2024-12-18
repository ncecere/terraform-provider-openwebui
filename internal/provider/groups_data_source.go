package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &GroupDataSource{}

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

// GroupDataSource defines the data source implementation.
type GroupDataSource struct {
	client *client.OpenWebUIClient
}

// GroupDataSourceModel describes the data source data model.
type GroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	UserIDs     types.List   `tfsdk:"user_ids"`
	Permissions types.Object `tfsdk:"permissions"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

func (d *GroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *GroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Group data source for OpenWebUI",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Group identifier",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the group to look up",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the group",
			},
			"user_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of user IDs in the group",
			},
			"permissions": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Permissions for the group",
				Attributes: map[string]schema.Attribute{
					"workspace": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"models":    schema.BoolAttribute{Computed: true},
							"knowledge": schema.BoolAttribute{Computed: true},
							"prompts":   schema.BoolAttribute{Computed: true},
							"tools":     schema.BoolAttribute{Computed: true},
						},
					},
					"chat": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"file_upload": schema.BoolAttribute{Computed: true},
							"delete":      schema.BoolAttribute{Computed: true},
							"edit":        schema.BoolAttribute{Computed: true},
							"temporary":   schema.BoolAttribute{Computed: true},
						},
					},
				},
			},
			"created_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the group was created",
			},
			"updated_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the group was last updated",
			},
		},
	}
}

func (d *GroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.OpenWebUIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.OpenWebUIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get groups from API
	groups, err := d.client.Groups.List()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read groups, got error: %s", err))
		return
	}

	// Find the group with matching name
	var found bool
	for _, group := range groups {
		if group.Name == data.Name.ValueString() {
			// Convert API response to model
			data.ID = types.StringValue(group.ID)
			data.Description = types.StringValue(group.Description)
			data.CreatedAt = types.Int64Value(group.CreatedAt)
			data.UpdatedAt = types.Int64Value(group.UpdatedAt)

			// Handle user IDs
			userIDs, diags := types.ListValueFrom(ctx, types.StringType, group.UserIDs)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			data.UserIDs = userIDs

			// Handle permissions
			if group.Permissions != nil {
				workspaceAttrs := map[string]attr.Value{
					"models":    types.BoolValue(group.Permissions.Workspace.Models),
					"knowledge": types.BoolValue(group.Permissions.Workspace.Knowledge),
					"prompts":   types.BoolValue(group.Permissions.Workspace.Prompts),
					"tools":     types.BoolValue(group.Permissions.Workspace.Tools),
				}

				chatAttrs := map[string]attr.Value{
					"file_upload": types.BoolValue(group.Permissions.Chat.FileUpload),
					"delete":      types.BoolValue(group.Permissions.Chat.Delete),
					"edit":        types.BoolValue(group.Permissions.Chat.Edit),
					"temporary":   types.BoolValue(group.Permissions.Chat.Temporary),
				}

				workspaceObj, diags := types.ObjectValue(
					map[string]attr.Type{
						"models":    types.BoolType,
						"knowledge": types.BoolType,
						"prompts":   types.BoolType,
						"tools":     types.BoolType,
					},
					workspaceAttrs,
				)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}

				chatObj, diags := types.ObjectValue(
					map[string]attr.Type{
						"file_upload": types.BoolType,
						"delete":      types.BoolType,
						"edit":        types.BoolType,
						"temporary":   types.BoolType,
					},
					chatAttrs,
				)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}

				permissionsObj, diags := types.ObjectValue(
					map[string]attr.Type{
						"workspace": types.ObjectType{AttrTypes: workspaceObj.AttributeTypes(ctx)},
						"chat":      types.ObjectType{AttrTypes: chatObj.AttributeTypes(ctx)},
					},
					map[string]attr.Value{
						"workspace": workspaceObj,
						"chat":      chatObj,
					},
				)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}

				data.Permissions = permissionsObj
			}

			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError(
			"Group Not Found",
			fmt.Sprintf("No group found with name: %s", data.Name.ValueString()),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
