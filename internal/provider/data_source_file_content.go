package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &fileContentDataSource{}
var _ datasource.DataSourceWithConfigure = &fileContentDataSource{}

type fileContentDataSource struct{ client *client.Client }

type fileContentDataSourceModel struct {
	FileID      types.String `tfsdk:"file_id"`
	ContentJSON types.String `tfsdk:"content_json"`
}

func NewFileContentDataSource() datasource.DataSource { return &fileContentDataSource{} }

func (d *fileContentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_content"
}

func (d *fileContentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"file_id":      schema.StringAttribute{Required: true, Description: "Identifier of the file to read content for."},
		"content_json": schema.StringAttribute{Computed: true, Description: "JSON response from the file data content endpoint."},
	}}
}

func (d *fileContentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

func (d *fileContentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the file content data source.")
		return
	}
	var config fileContentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.FileID.IsNull() || config.FileID.IsUnknown() || config.FileID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("file_id"), "Missing file identifier", "The file_id argument must be supplied.")
		return
	}
	content, err := d.client.GetFileDataContent(ctx, config.FileID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read file content failed", err.Error())
		return
	}
	contentJSON, err := encodeOptionalJSON(content)
	if err != nil {
		resp.Diagnostics.AddError("Serialize file content", err.Error())
		return
	}
	state := fileContentDataSourceModel{FileID: config.FileID, ContentJSON: contentJSON}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
