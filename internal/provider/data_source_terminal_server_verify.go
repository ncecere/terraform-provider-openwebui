package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &terminalServerVerifyDataSource{}
var _ datasource.DataSourceWithConfigure = &terminalServerVerifyDataSource{}

type terminalServerVerifyDataSource struct{ client *client.Client }

type terminalServerVerifyModel struct {
	URL        types.String `tfsdk:"url"`
	Path       types.String `tfsdk:"path"`
	Key        types.String `tfsdk:"key"`
	AuthType   types.String `tfsdk:"auth_type"`
	ConfigJSON types.String `tfsdk:"config_json"`
	Verified   types.Bool   `tfsdk:"verified"`
}

func NewTerminalServerVerifyDataSource() datasource.DataSource {
	return &terminalServerVerifyDataSource{}
}

func (d *terminalServerVerifyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_terminal_server_verify"
}

func (d *terminalServerVerifyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"url":         schema.StringAttribute{Required: true, Description: "Terminal server base URL to verify."},
		"path":        schema.StringAttribute{Optional: true, Description: "Optional terminal server path."},
		"key":         schema.StringAttribute{Optional: true, Sensitive: true, Description: "Optional authentication key."},
		"auth_type":   schema.StringAttribute{Optional: true, Description: "Optional authentication type."},
		"config_json": schema.StringAttribute{Optional: true, Description: "Optional terminal server config JSON."},
		"verified":    schema.BoolAttribute{Computed: true, Description: "Whether terminal server verification succeeded."},
	}}
}

func (d *terminalServerVerifyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

func (d *terminalServerVerifyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the terminal server verify data source.")
		return
	}

	var config terminalServerVerifyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configMap := decodeOptionalJSON(config.ConfigJSON, path.Root("config_json"), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	connection := client.TerminalServerConnection{URL: config.URL.ValueString(), Config: configMap}
	if !config.Path.IsNull() && !config.Path.IsUnknown() {
		value := config.Path.ValueString()
		connection.Path = &value
	}
	if !config.Key.IsNull() && !config.Key.IsUnknown() {
		value := config.Key.ValueString()
		connection.Key = &value
	}
	if !config.AuthType.IsNull() && !config.AuthType.IsUnknown() {
		value := config.AuthType.ValueString()
		connection.AuthType = &value
	}

	if err := d.client.VerifyTerminalServer(ctx, connection); err != nil {
		resp.Diagnostics.AddError("Verify terminal server failed", err.Error())
		return
	}

	state := config
	state.Verified = types.BoolValue(true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
