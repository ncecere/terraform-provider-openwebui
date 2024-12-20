package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &KnowledgeDataSource{}

func NewKnowledgeDataSource() datasource.DataSource {
	return &KnowledgeDataSource{}
}

// KnowledgeDataSource defines the data source implementation.
type KnowledgeDataSource struct {
	client *client.OpenWebUIClient
}

// KnowledgeDataSourceModel describes the data source data model.
type KnowledgeDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Data          types.Map    `tfsdk:"data"`
	AccessControl types.String `tfsdk:"access_control"`
	AccessGroups  types.List   `tfsdk:"access_groups"`
	AccessUsers   types.List   `tfsdk:"access_users"`
	LastUpdated   types.String `tfsdk:"last_updated"`
}

func (d *KnowledgeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

func (d *KnowledgeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Knowledge data source for OpenWebUI",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Knowledge identifier",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the knowledge base to look up",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the knowledge base",
			},
			"data": schema.MapAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Additional data for the knowledge base",
			},
			"access_control": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Access control type ('public' or 'private')",
			},
			"access_groups": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of group IDs with access",
			},
			"access_users": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of user IDs with access",
			},
			"last_updated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of the last update",
			},
		},
	}
}

func (d *KnowledgeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KnowledgeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KnowledgeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get knowledge bases from API
	knowledgeBases, err := d.client.List()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read knowledge bases, got error: %s", err))
		return
	}

	// Find the knowledge base with matching name
	var found bool
	for _, kb := range knowledgeBases {
		if kb.Name == data.Name.ValueString() {
			// Convert API response to model
			data.ID = types.StringValue(kb.ID)
			data.Description = types.StringValue(kb.Description)
			data.LastUpdated = types.StringValue(fmt.Sprint(kb.UpdatedAt))

			// Convert data map
			if kb.Data != nil {
				dataMap := make(map[string]string)
				for k, v := range kb.Data {
					if str, ok := v.(string); ok {
						dataMap[k] = str
					}
				}
				convertedMap, err := types.MapValueFrom(ctx, types.StringType, dataMap)
				if err != nil {
					resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Unable to convert data map, got error: %s", err))
					return
				}
				data.Data = convertedMap
			}

			// Handle access control
			if kb.AccessControl == nil {
				data.AccessControl = types.StringValue("public")
			} else {
				data.AccessControl = types.StringValue("private")

				// Extract groups and users from access control
				if ac, ok := kb.AccessControl.(map[string]interface{}); ok {
					if read, ok := ac["read"].(map[string]interface{}); ok {
						// Handle groups
						if groups, ok := read["group_ids"].([]interface{}); ok {
							groupList := make([]string, 0, len(groups))
							for _, g := range groups {
								if str, ok := g.(string); ok {
									groupList = append(groupList, str)
								}
							}
							groupsValue, err := types.ListValueFrom(ctx, types.StringType, groupList)
							if err != nil {
								resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Unable to convert groups list, got error: %s", err))
								return
							}
							data.AccessGroups = groupsValue
						}

						// Handle users
						if users, ok := read["user_ids"].([]interface{}); ok {
							userList := make([]string, 0, len(users))
							for _, u := range users {
								if str, ok := u.(string); ok {
									userList = append(userList, str)
								}
							}
							usersValue, err := types.ListValueFrom(ctx, types.StringType, userList)
							if err != nil {
								resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Unable to convert users list, got error: %s", err))
								return
							}
							data.AccessUsers = usersValue
						}
					}
				}
			}

			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError(
			"Knowledge Base Not Found",
			fmt.Sprintf("No knowledge base found with name: %s", data.Name.ValueString()),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
