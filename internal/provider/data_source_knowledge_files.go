package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &knowledgeFilesDataSource{}
var _ datasource.DataSourceWithConfigure = &knowledgeFilesDataSource{}

type knowledgeFilesDataSource struct{ client *client.Client }

type knowledgeFilesDataSourceModel struct {
	KnowledgeID     types.String `tfsdk:"knowledge_id"`
	Query           types.String `tfsdk:"query"`
	DirectoryID     types.String `tfsdk:"directory_id"`
	IncludeContent  types.Bool   `tfsdk:"include_content"`
	ViewOption      types.String `tfsdk:"view_option"`
	OrderBy         types.String `tfsdk:"order_by"`
	Direction       types.String `tfsdk:"direction"`
	Page            types.Int64  `tfsdk:"page"`
	FilesJSON       types.String `tfsdk:"files_json"`
	DirectoriesJSON types.String `tfsdk:"directories_json"`
	BreadcrumbsJSON types.String `tfsdk:"breadcrumbs_json"`
	Total           types.Int64  `tfsdk:"total"`
}

func NewKnowledgeFilesDataSource() datasource.DataSource { return &knowledgeFilesDataSource{} }
func (d *knowledgeFilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge_files"
}
func (d *knowledgeFilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"knowledge_id":     schema.StringAttribute{Required: true, Description: "Knowledge base ID."},
		"query":            schema.StringAttribute{Optional: true},
		"directory_id":     schema.StringAttribute{Optional: true, Description: "Directory filter. Use empty string for root."},
		"include_content":  schema.BoolAttribute{Optional: true},
		"view_option":      schema.StringAttribute{Optional: true},
		"order_by":         schema.StringAttribute{Optional: true},
		"direction":        schema.StringAttribute{Optional: true},
		"page":             schema.Int64Attribute{Optional: true},
		"files_json":       schema.StringAttribute{Computed: true},
		"directories_json": schema.StringAttribute{Computed: true},
		"breadcrumbs_json": schema.StringAttribute{Computed: true},
		"total":            schema.Int64Attribute{Computed: true},
	}}
}
func (d *knowledgeFilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}
func (d *knowledgeFilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the knowledge files data source.")
		return
	}
	var config knowledgeFilesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.KnowledgeID.IsNull() || config.KnowledgeID.IsUnknown() || config.KnowledgeID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("knowledge_id"), "Missing knowledge ID", "knowledge_id is required.")
		return
	}
	str := func(v types.String) string {
		if v.IsNull() || v.IsUnknown() {
			return ""
		}
		return v.ValueString()
	}
	include := false
	if !config.IncludeContent.IsNull() && !config.IncludeContent.IsUnknown() {
		include = config.IncludeContent.ValueBool()
	}
	page := 1
	if !config.Page.IsNull() && !config.Page.IsUnknown() {
		page = int(config.Page.ValueInt64())
	}
	list, err := d.client.ListKnowledgeFilesFiltered(ctx, config.KnowledgeID.ValueString(), str(config.Query), include, str(config.ViewOption), str(config.OrderBy), str(config.Direction), str(config.DirectoryID), page)
	if err != nil {
		resp.Diagnostics.AddError("List knowledge files failed", err.Error())
		return
	}
	filesJSON, err := encodeOptionalJSONValue(list.Items)
	if err != nil {
		resp.Diagnostics.AddError("Serialize knowledge files", err.Error())
		return
	}
	dirsJSON, err := encodeOptionalJSONValue(list.Directories)
	if err != nil {
		resp.Diagnostics.AddError("Serialize knowledge directories", err.Error())
		return
	}
	crumbsJSON, err := encodeOptionalJSONValue(list.Breadcrumbs)
	if err != nil {
		resp.Diagnostics.AddError("Serialize knowledge breadcrumbs", err.Error())
		return
	}
	state := config
	state.FilesJSON = filesJSON
	state.DirectoriesJSON = dirsJSON
	state.BreadcrumbsJSON = crumbsJSON
	state.Total = types.Int64Value(int64(list.Total))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
