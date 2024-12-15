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
		MarkdownDescription: "Interact with OpenWebUI.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint URL of the OpenWebUI instance. May also be provided via OPENWEBUI_ENDPOINT environment variable.",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The API token to authenticate with OpenWebUI. May also be provided via OPENWEBUI_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *OpenWebUIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OpenWebUIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Use environment variables if config attributes are not set
	endpoint := os.Getenv("OPENWEBUI_ENDPOINT")
	token := os.Getenv("OPENWEBUI_TOKEN")

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	// Error if required values are not set
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

	// Create a new OpenWebUI client using the configuration values
	client := NewClient(endpoint, token)

	// Make the client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpenWebUIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKnowledgeResource,
	}
}

func (p *OpenWebUIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenWebUIProvider{
			version: version,
		}
	}
}
