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
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client"
)

// Ensure OpenWebUIProvider satisfies various provider interfaces.
var _ provider.Provider = &OpenWebUIProvider{}

// OpenWebUIProvider defines the provider implementation.
type OpenWebUIProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// OpenWebUIProviderModel describes the provider data model.
type OpenWebUIProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *OpenWebUIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openwebui"
	resp.Version = p.version
}

func (p *OpenWebUIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "OpenWebUI API endpoint URL. Can also be set with OPENWEBUI_ENDPOINT environment variable.",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "OpenWebUI API token. Can also be set with OPENWEBUI_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *OpenWebUIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config OpenWebUIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check environment variables
	endpoint := os.Getenv("OPENWEBUI_ENDPOINT")
	token := os.Getenv("OPENWEBUI_TOKEN")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing OpenWebUI API Endpoint",
			"The provider cannot create the OpenWebUI API client as there is a missing or empty value for the OpenWebUI API endpoint. "+
				"Set the endpoint value in the configuration or use the OPENWEBUI_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing OpenWebUI API Token",
			"The provider cannot create the OpenWebUI API client as there is a missing or empty value for the OpenWebUI API token. "+
				"Set the token value in the configuration or use the OPENWEBUI_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new OpenWebUI client using the configuration
	client, err := client.New(endpoint, token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OpenWebUI API Client",
			"An unexpected error occurred when creating the OpenWebUI API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OpenWebUI Client Error: "+err.Error(),
		)
		return
	}

	// Make the client available during DataSource and Resource Configure methods
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpenWebUIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKnowledgeResource,
	}
}

func (p *OpenWebUIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKnowledgeDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenWebUIProvider{
			version: version,
		}
	}
}
