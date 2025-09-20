package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ provider.Provider = &openWebUIProvider{}

const defaultEndpoint = "http://localhost:3000/api/v1"

// openWebUIProvider defines the provider implementation.
type openWebUIProvider struct {
	version string
}

// providerModel maps provider schema data to Go type.
type providerModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

// New instantiates a new provider.
func New() provider.Provider {
	return &openWebUIProvider{
		version: Version,
	}
}

// Metadata satisfies the provider.Provider interface.
func (p *openWebUIProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openwebui"
	resp.Version = p.version
}

// Schema defines the provider-level schema.
func (p *openWebUIProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Base URL for the Open WebUI API. Defaults to http://localhost:3000/api/v1.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "API token used to authenticate against the Open WebUI API. Can also be supplied via the OPENWEBUI_TOKEN environment variable.",
			},
		},
	}
}

// Configure prepares the Open WebUI API client for data sources and resources.
func (p *openWebUIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := defaultEndpoint
	if !data.Endpoint.IsNull() && !data.Endpoint.IsUnknown() {
		endpoint = data.Endpoint.ValueString()
	} else if envEndpoint := os.Getenv("OPENWEBUI_ENDPOINT"); envEndpoint != "" {
		endpoint = envEndpoint
	}

	token := ""
	if !data.Token.IsNull() && !data.Token.IsUnknown() {
		token = data.Token.ValueString()
	} else if envToken := os.Getenv("OPENWEBUI_TOKEN"); envToken != "" {
		token = envToken
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Open WebUI API token",
			"A valid API token must be supplied via the provider configuration or the OPENWEBUI_TOKEN environment variable.",
		)
		return
	}

	apiClient, err := client.NewClient(endpoint, token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Open WebUI API client",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Configured Open WebUI provider", map[string]any{
		"endpoint": endpoint,
	})

	resp.ResourceData = apiClient
	resp.DataSourceData = apiClient
}

// Resources defines provider-supported resources.
func (p *openWebUIProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKnowledgeResource,
		NewModelResource,
		NewPromptResource,
		NewGroupResource,
	}
}

// DataSources defines provider-supported data sources.
func (p *openWebUIProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewModelDataSource,
		NewKnowledgeDataSource,
		NewGroupDataSource,
		NewPromptDataSource,
	}
}
