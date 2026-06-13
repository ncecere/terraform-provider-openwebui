package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

type jsonEndpointDataSource struct {
	client   *client.Client
	typeName string
	desc     string
	read     func(context.Context, *client.Client) (any, error)
}

type jsonEndpointModel struct {
	JSON types.String `tfsdk:"json"`
}

func newJSONEndpointDataSource(typeName string, desc string, read func(context.Context, *client.Client) (any, error)) datasource.DataSource {
	return &jsonEndpointDataSource{typeName: typeName, desc: desc, read: read}
}

func NewModelsDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("models", "Reads the raw Open WebUI models list payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ListModelsRaw(ctx) })
}
func NewBaseModelsDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("base_models", "Reads the raw Open WebUI base models payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ListBaseModelsRaw(ctx) })
}
func NewModelTagsDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("model_tags", "Reads the raw Open WebUI model tags payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ListModelTagsRaw(ctx) })
}
func NewModelsExportDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("models_export", "Reads the raw Open WebUI model export payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ExportModelsRaw(ctx) })
}
func NewPromptsDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("prompts", "Reads the raw Open WebUI prompts list payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ListPromptsRaw(ctx) })
}
func NewPromptTagsDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("prompt_tags", "Reads the raw Open WebUI prompt tags payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ListPromptTagsRaw(ctx) })
}
func NewToolsDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("tools", "Reads the raw Open WebUI tools list payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ListToolsRaw(ctx) })
}
func NewToolsExportDataSource() datasource.DataSource {
	return newJSONEndpointDataSource("tools_export", "Reads the raw Open WebUI tools export payload.", func(ctx context.Context, c *client.Client) (any, error) { return c.ExportToolsRaw(ctx) })
}

func (d *jsonEndpointDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.typeName
}

func (d *jsonEndpointDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: d.desc, Attributes: map[string]schema.Attribute{
		"json": schema.StringAttribute{Computed: true, Description: "Raw JSON response payload."},
	}}
}

func (d *jsonEndpointDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

func (d *jsonEndpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using this data source.")
		return
	}
	payload, err := d.read(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Read data source failed", err.Error())
		return
	}
	payloadJSON, err := encodeOptionalJSONValue(payload)
	if err != nil {
		resp.Diagnostics.AddError("Serialize data source payload", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &jsonEndpointModel{JSON: payloadJSON})...)
}
