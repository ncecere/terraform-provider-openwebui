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

var _ datasource.DataSource = &knowledgeDataSource{}
var _ datasource.DataSourceWithConfigure = &knowledgeDataSource{}

// knowledgeDataSource surfaces knowledge base entries.
type knowledgeDataSource struct {
	client *client.Client
}

// knowledgeDataSourceModel embeds the resource fields and adds the lookup identifier.
type knowledgeDataSourceModel struct {
	KnowledgeID types.String `tfsdk:"knowledge_id"`
	knowledgeResourceModel
}

// NewKnowledgeDataSource creates a new knowledge data source instance.
func NewKnowledgeDataSource() datasource.DataSource {
	return &knowledgeDataSource{}
}

// Metadata sets the data source type name.
func (d *knowledgeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

// Schema describes the knowledge data source schema.
func (d *knowledgeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"knowledge_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the knowledge entry to retrieve. If omitted, `name` must be provided.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Case-insensitive knowledge name used to locate the entry when `knowledge_id` is not supplied.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier assigned by Open WebUI.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Detailed description of the knowledge entry.",
			},
			"data_json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON payload describing additional knowledge data.",
			},
			"meta_json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON payload describing metadata for the knowledge entry.",
			},
			"read_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Group names granted read access.",
			},
			"write_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Group names granted write access.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Creation date in YYYY-MM-DD format.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Last update date in YYYY-MM-DD format.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the user who owns the knowledge entry.",
			},
		},
	}
}

// Configure attaches the provider client.
func (d *knowledgeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

// Read looks up the requested knowledge entry.
func (d *knowledgeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the knowledge data source.")
		return
	}

	var config knowledgeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := ""
	if !config.KnowledgeID.IsNull() && !config.KnowledgeID.IsUnknown() {
		id = strings.TrimSpace(config.KnowledgeID.ValueString())
	}

	name := ""
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		name = strings.TrimSpace(config.Name.ValueString())
	}

	if id == "" {
		if name == "" {
			resp.Diagnostics.AddError(
				"Missing knowledge lookup value",
				"Either knowledge_id or name must be provided to query an existing knowledge entry.",
			)
			return
		}

		entries, err := d.client.ListKnowledge(ctx)
		if err != nil {
			resp.Diagnostics.AddError("List knowledge entries failed", err.Error())
			return
		}

		var matches []client.KnowledgeListItem
		for _, entry := range entries {
			if strings.EqualFold(entry.Name, name) {
				matches = append(matches, entry)
			}
		}

		if len(matches) == 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("name"),
				"Knowledge entry not found",
				fmt.Sprintf("No Open WebUI knowledge entry was found with the name %q.", name),
			)
			return
		}
		if len(matches) > 1 {
			resp.Diagnostics.AddAttributeError(
				path.Root("name"),
				"Knowledge entry name not unique",
				fmt.Sprintf("Multiple knowledge entries share the name %q. Provide knowledge_id instead.", name),
			)
			return
		}

		id = matches[0].ID
	}

	current, err := d.client.GetKnowledge(ctx, id)
	if err != nil {
		if err == client.ErrNotFound {
			resp.Diagnostics.AddAttributeError(
				path.Root("knowledge_id"),
				"Knowledge entry not found",
				"No Open WebUI knowledge entry was found with the supplied knowledge_id.",
			)
			return
		}

		resp.Diagnostics.AddError("Read knowledge entry failed", err.Error())
		return
	}

	model, diags := knowledgeResponseToModel(ctx, d.client, *current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := knowledgeDataSourceModel{
		KnowledgeID:            types.StringValue(current.ID),
		knowledgeResourceModel: model,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
