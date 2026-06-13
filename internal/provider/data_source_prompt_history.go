package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &promptHistoryDataSource{}
var _ datasource.DataSourceWithConfigure = &promptHistoryDataSource{}

type promptHistoryDataSource struct{ client *client.Client }

type promptHistoryDataSourceModel struct {
	Command     types.String `tfsdk:"command"`
	Page        types.Int64  `tfsdk:"page"`
	HistoryJSON types.String `tfsdk:"history_json"`
}

func NewPromptHistoryDataSource() datasource.DataSource { return &promptHistoryDataSource{} }
func (d *promptHistoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prompt_history"
}
func (d *promptHistoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"command":      schema.StringAttribute{Required: true, Description: "Prompt command."},
		"page":         schema.Int64Attribute{Optional: true},
		"history_json": schema.StringAttribute{Computed: true},
	}}
}
func (d *promptHistoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		d.client = c
	}
}
func (d *promptHistoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using prompt history.")
		return
	}
	var config promptHistoryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.Command.IsNull() || config.Command.IsUnknown() || config.Command.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("command"), "Missing command", "command is required.")
		return
	}
	page := 0
	if !config.Page.IsNull() && !config.Page.IsUnknown() {
		page = int(config.Page.ValueInt64())
	}
	history, err := d.client.GetPromptHistoryRaw(ctx, config.Command.ValueString(), page)
	if err != nil {
		resp.Diagnostics.AddError("Read prompt history failed", err.Error())
		return
	}
	historyJSON, err := encodeOptionalJSONValue(history)
	if err != nil {
		resp.Diagnostics.AddError("Serialize prompt history", err.Error())
		return
	}
	state := config
	state.HistoryJSON = historyJSON
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
