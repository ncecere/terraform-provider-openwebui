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
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/groups"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/knowledge"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/models"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/users"
)

var (
	_ provider.Provider = &OpenWebUIProvider{}
)

type OpenWebUIProvider struct {
	version string
}

type OpenWebUIProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *OpenWebUIProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openwebui"
	resp.Version = p.version
}

func (p *OpenWebUIProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with OpenWebUI.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The endpoint URL of the OpenWebUI API. May also be provided via OPENWEBUI_ENDPOINT environment variable.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "The token to authenticate with the OpenWebUI API. May also be provided via OPENWEBUI_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *OpenWebUIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config OpenWebUIProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsNull() {
		endpoint := os.Getenv("OPENWEBUI_ENDPOINT")
		config.Endpoint = types.StringValue(endpoint)
	}

	if config.Token.IsNull() {
		token := os.Getenv("OPENWEBUI_TOKEN")
		config.Token = types.StringValue(token)
	}

	if config.Endpoint.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing OpenWebUI API Endpoint",
			"The provider cannot create the OpenWebUI API client as there is a missing or empty value for the OpenWebUI API endpoint. "+
				"Set the endpoint value in the configuration or use the OPENWEBUI_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if config.Token.IsNull() {
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

	// Create new OpenWebUI clients
	groupsClient := groups.NewClient(config.Endpoint.ValueString(), config.Token.ValueString())
	knowledgeClient := knowledge.NewClient(config.Endpoint.ValueString(), config.Token.ValueString())
	modelsClient := models.NewClient(config.Endpoint.ValueString(), config.Token.ValueString())
	usersClient := users.NewClient(config.Endpoint.ValueString(), config.Token.ValueString())

	// Create a map to store all clients
	clients := map[string]interface{}{
		"groups":    groupsClient,
		"knowledge": knowledgeClient,
		"models":    modelsClient,
		"users":     usersClient,
	}

	resp.DataSourceData = clients
	resp.ResourceData = clients
}

func (p *OpenWebUIProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGroupDataSource,
		NewKnowledgeDataSource,
		NewModelDataSource,
		NewUserDataSource,
	}
}

func (p *OpenWebUIProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGroupResource,
		NewKnowledgeResource,
		NewModelResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenWebUIProvider{
			version: version,
		}
	}
}
