package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &promptHistoryDiffDataSource{}
var _ datasource.DataSourceWithConfigure = &promptHistoryDiffDataSource{}

type promptHistoryDiffDataSource struct{ client *client.Client }
type promptHistoryDiffDataSourceModel struct {
	Command  types.String `tfsdk:"command"`
	FromID   types.String `tfsdk:"from_id"`
	ToID     types.String `tfsdk:"to_id"`
	DiffJSON types.String `tfsdk:"diff_json"`
}

func NewPromptHistoryDiffDataSource() datasource.DataSource { return &promptHistoryDiffDataSource{} }
func (d *promptHistoryDiffDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prompt_history_diff"
}
func (d *promptHistoryDiffDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{"command": schema.StringAttribute{Required: true, Description: "Prompt command."}, "from_id": schema.StringAttribute{Required: true, Description: "Source prompt history entry id."}, "to_id": schema.StringAttribute{Required: true, Description: "Target prompt history entry id."}, "diff_json": schema.StringAttribute{Computed: true, Description: "Raw prompt history diff JSON."}}}
}
func (d *promptHistoryDiffDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		d.client = c
	}
}
func (d *promptHistoryDiffDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using prompt history diff.")
		return
	}
	var config promptHistoryDiffDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.Command.IsNull() || config.Command.IsUnknown() || config.Command.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("command"), "Missing command", "command is required.")
		return
	}
	if config.FromID.IsNull() || config.FromID.IsUnknown() || config.FromID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("from_id"), "Missing from_id", "from_id is required.")
		return
	}
	if config.ToID.IsNull() || config.ToID.IsUnknown() || config.ToID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("to_id"), "Missing to_id", "to_id is required.")
		return
	}
	diff, err := d.client.GetPromptHistoryDiffRaw(ctx, config.Command.ValueString(), config.FromID.ValueString(), config.ToID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read prompt history diff failed", err.Error())
		return
	}
	diffJSON, err := encodeOptionalJSONValue(diff)
	if err != nil {
		resp.Diagnostics.AddError("Serialize prompt history diff", err.Error())
		return
	}
	config.DiffJSON = diffJSON
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
