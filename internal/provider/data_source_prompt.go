package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &promptDataSource{}
var _ datasource.DataSourceWithConfigure = &promptDataSource{}

// promptDataSource exposes prompt definitions.
type promptDataSource struct {
	client *client.Client
}

// promptDataSourceModel reuses the resource representation.
type promptDataSourceModel = promptResourceModel

// NewPromptDataSource constructs a new prompt data source.
func NewPromptDataSource() datasource.DataSource {
	return &promptDataSource{}
}

// Metadata sets the data source identifier.
func (d *promptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prompt"
}

// Schema describes the prompt data source schema.
func (d *promptDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"command": schema.StringAttribute{
				Required:    true,
				Description: "Command string that uniquely identifies the prompt.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the prompt (mirrors the command).",
			},
			"title": schema.StringAttribute{
				Computed:    true,
				Description: "Prompt title displayed in Open WebUI.",
			},
			"content": schema.StringAttribute{
				Computed:    true,
				Description: "Prompt content text.",
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
			"timestamp": schema.StringAttribute{
				Computed:    true,
				Description: "Prompt timestamp formatted as YYYY-MM-DD.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the user who owns the prompt.",
			},
		},
	}
}

// Configure attaches the API client.
func (d *promptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

// Read retrieves the prompt definition.
func (d *promptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the prompt data source.")
		return
	}

	var config promptDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Command.IsUnknown() || config.Command.IsNull() || config.Command.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("command"),
			"Missing prompt command",
			"The command argument must be supplied to query an existing prompt.",
		)
		return
	}

	current, err := d.client.GetPrompt(ctx, config.Command.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.Diagnostics.AddAttributeError(
				path.Root("command"),
				"Prompt not found",
				"No Open WebUI prompt was found with the supplied command.",
			)
			return
		}

		resp.Diagnostics.AddError("Read prompt failed", err.Error())
		return
	}

	state, diags := promptResponseToModel(ctx, d.client, current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
