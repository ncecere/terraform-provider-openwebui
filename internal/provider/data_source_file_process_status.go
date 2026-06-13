package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &fileProcessStatusDataSource{}
var _ datasource.DataSourceWithConfigure = &fileProcessStatusDataSource{}

type fileProcessStatusDataSource struct{ client *client.Client }

type fileProcessStatusDataSourceModel struct {
	FileID     types.String `tfsdk:"file_id"`
	StatusJSON types.String `tfsdk:"status_json"`
}

func NewFileProcessStatusDataSource() datasource.DataSource { return &fileProcessStatusDataSource{} }

func (d *fileProcessStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_process_status"
}

func (d *fileProcessStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"file_id":     schema.StringAttribute{Required: true, Description: "Identifier of the file to read processing status for."},
		"status_json": schema.StringAttribute{Computed: true, Description: "JSON response from the file process status endpoint."},
	}}
}

func (d *fileProcessStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

func (d *fileProcessStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the file process status data source.")
		return
	}
	var config fileProcessStatusDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.FileID.IsNull() || config.FileID.IsUnknown() || config.FileID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("file_id"), "Missing file identifier", "The file_id argument must be supplied.")
		return
	}
	status, err := d.client.GetFileProcessStatus(ctx, config.FileID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read file process status failed", err.Error())
		return
	}
	statusJSON, err := encodeOptionalJSON(status)
	if err != nil {
		resp.Diagnostics.AddError("Serialize file process status", err.Error())
		return
	}
	state := fileProcessStatusDataSourceModel{FileID: config.FileID, StatusJSON: statusJSON}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
