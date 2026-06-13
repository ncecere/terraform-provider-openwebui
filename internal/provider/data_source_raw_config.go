package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

type rawConfigDataSource struct {
	client   *client.Client
	typeName string
	desc     string
	read     func(context.Context, *client.Client) (map[string]any, error)
}

type rawConfigDataSourceModel struct {
	ConfigJSON types.String `tfsdk:"config_json"`
}

func newRawConfigDataSource(typeName, desc string, read func(context.Context, *client.Client) (map[string]any, error)) datasource.DataSource {
	return &rawConfigDataSource{typeName: typeName, desc: desc, read: read}
}

func NewRetrievalConfigDataSource() datasource.DataSource {
	return newRawConfigDataSource("retrieval_config", "Reads Open WebUI retrieval configuration.", func(ctx context.Context, c *client.Client) (map[string]any, error) { return c.GetRetrievalConfig(ctx) })
}
func NewEvaluationsConfigDataSource() datasource.DataSource {
	return newRawConfigDataSource("evaluations_config", "Reads Open WebUI evaluations configuration.", func(ctx context.Context, c *client.Client) (map[string]any, error) {
		return c.GetEvaluationsConfig(ctx)
	})
}
func NewDefaultUserPermissionsDataSource() datasource.DataSource {
	return newRawConfigDataSource("default_user_permissions", "Reads Open WebUI default user permissions.", func(ctx context.Context, c *client.Client) (map[string]any, error) {
		return c.GetDefaultUserPermissions(ctx)
	})
}

func (d *rawConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.typeName
}
func (d *rawConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: d.desc, Attributes: map[string]schema.Attribute{
		"config_json": schema.StringAttribute{Computed: true, Description: "Raw JSON configuration payload."},
	}}
}
func (d *rawConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}
func (d *rawConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using this data source.")
		return
	}
	config, err := d.read(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Read config failed", err.Error())
		return
	}
	configJSON, err := encodeOptionalJSON(config)
	if err != nil {
		resp.Diagnostics.AddError("Serialize config", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &rawConfigDataSourceModel{ConfigJSON: configJSON})...)
}
