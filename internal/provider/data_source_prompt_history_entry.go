package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &promptHistoryEntryDataSource{}
var _ datasource.DataSourceWithConfigure = &promptHistoryEntryDataSource{}

type promptHistoryEntryDataSource struct{ client *client.Client }
type promptHistoryEntryDataSourceModel struct {
	Command     types.String `tfsdk:"command"`
	HistoryID   types.String `tfsdk:"history_id"`
	HistoryJSON types.String `tfsdk:"history_json"`
}

func NewPromptHistoryEntryDataSource() datasource.DataSource { return &promptHistoryEntryDataSource{} }
func (d *promptHistoryEntryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prompt_history_entry"
}
func (d *promptHistoryEntryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{"command": schema.StringAttribute{Required: true, Description: "Prompt command."}, "history_id": schema.StringAttribute{Required: true, Description: "Prompt history entry id."}, "history_json": schema.StringAttribute{Computed: true, Description: "Raw prompt history entry JSON."}}}
}
func (d *promptHistoryEntryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		d.client = c
	}
}
func (d *promptHistoryEntryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using prompt history entry.")
		return
	}
	var config promptHistoryEntryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.Command.IsNull() || config.Command.IsUnknown() || config.Command.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("command"), "Missing command", "command is required.")
		return
	}
	if config.HistoryID.IsNull() || config.HistoryID.IsUnknown() || config.HistoryID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("history_id"), "Missing history_id", "history_id is required.")
		return
	}
	entry, err := d.client.GetPromptHistoryEntryRaw(ctx, config.Command.ValueString(), config.HistoryID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read prompt history entry failed", err.Error())
		return
	}
	entryJSON, err := encodeOptionalJSONValue(entry)
	if err != nil {
		resp.Diagnostics.AddError("Serialize prompt history entry", err.Error())
		return
	}
	config.HistoryJSON = entryJSON
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
