package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &functionDataSource{}
var _ datasource.DataSourceWithConfigure = &functionDataSource{}

type functionDataSource struct{ client *client.Client }
type functionDataSourceModel = functionResourceModel

func NewFunctionDataSource() datasource.DataSource { return &functionDataSource{} }
func NewFunctionsDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("functions", "Reads the raw Open WebUI functions list payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ListFunctionsRaw(ctx) })
}
func (d *functionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function"
}
func (d *functionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{"function_id": schema.StringAttribute{Required: true}, "id": schema.StringAttribute{Computed: true}, "name": schema.StringAttribute{Computed: true}, "content": schema.StringAttribute{Computed: true}, "description": schema.StringAttribute{Computed: true}, "manifest_json": schema.StringAttribute{Computed: true}, "type": schema.StringAttribute{Computed: true}, "is_active": schema.BoolAttribute{Computed: true}, "is_global": schema.BoolAttribute{Computed: true}, "user_id": schema.StringAttribute{Computed: true}, "created_at": schema.Int64Attribute{Computed: true}, "updated_at": schema.Int64Attribute{Computed: true}}}
}
func (d *functionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		d.client = c
	}
}
func (d *functionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the function data source.")
		return
	}
	var id types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("function_id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cur, err := d.client.GetFunction(ctx, id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read function failed", err.Error())
		return
	}
	state, diags := functionResponseToModel(cur)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
