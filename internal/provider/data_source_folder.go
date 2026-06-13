package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &folderDataSource{}
var _ datasource.DataSourceWithConfigure = &folderDataSource{}

type folderDataSource struct{ client *client.Client }
type folderDataSourceModel = folderModel

func NewFolderDataSource() datasource.DataSource { return &folderDataSource{} }
func (d *folderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}
func (d *folderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":              schema.StringAttribute{Required: true, Description: "Folder ID."},
		"name":            schema.StringAttribute{Computed: true},
		"parent_id":       schema.StringAttribute{Computed: true},
		"data_json":       schema.StringAttribute{Computed: true},
		"meta_json":       schema.StringAttribute{Computed: true},
		"is_expanded":     schema.BoolAttribute{Computed: true},
		"delete_contents": schema.BoolAttribute{Computed: true},
		"user_id":         schema.StringAttribute{Computed: true},
		"created_at":      schema.Int64Attribute{Computed: true},
		"updated_at":      schema.Int64Attribute{Computed: true},
	}}
}
func (d *folderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}
func (d *folderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the folder data source.")
		return
	}
	var config folderDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.ID.IsNull() || config.ID.IsUnknown() || config.ID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Missing folder ID", "The id argument must be supplied to query an existing folder.")
		return
	}
	current, err := d.client.GetFolder(ctx, config.ID.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.Diagnostics.AddAttributeError(path.Root("id"), "Folder not found", "No Open WebUI folder was found with the supplied id.")
			return
		}
		resp.Diagnostics.AddError("Read folder failed", err.Error())
		return
	}
	state := folderResponseToModel(current)
	state.DeleteContents = types.BoolNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
