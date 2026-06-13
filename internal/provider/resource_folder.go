package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ resource.Resource = &folderResource{}
var _ resource.ResourceWithConfigure = &folderResource{}
var _ resource.ResourceWithImportState = &folderResource{}

type folderResource struct{ client *client.Client }

type folderModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ParentID       types.String `tfsdk:"parent_id"`
	DataJSON       types.String `tfsdk:"data_json"`
	MetaJSON       types.String `tfsdk:"meta_json"`
	IsExpanded     types.Bool   `tfsdk:"is_expanded"`
	DeleteContents types.Bool   `tfsdk:"delete_contents"`
	UserID         types.String `tfsdk:"user_id"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewFolderResource() resource.Resource { return &folderResource{} }
func (r *folderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}
func (r *folderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"name":            schema.StringAttribute{Required: true, Description: "Folder name."},
		"parent_id":       schema.StringAttribute{Optional: true, Description: "Parent folder ID."},
		"data_json":       schema.StringAttribute{Optional: true, Computed: true, Description: "Folder data as JSON."},
		"meta_json":       schema.StringAttribute{Optional: true, Computed: true, Description: "Folder metadata as JSON."},
		"is_expanded":     schema.BoolAttribute{Optional: true, Computed: true, Description: "Whether the folder is expanded."},
		"delete_contents": schema.BoolAttribute{Optional: true, Description: "Whether to delete folder contents on destroy. Defaults to true."},
		"user_id":         schema.StringAttribute{Computed: true},
		"created_at":      schema.Int64Attribute{Computed: true},
		"updated_at":      schema.Int64Attribute{Computed: true},
	}}
}
func (r *folderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}
func (r *folderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing folders.")
		return
	}
	var plan folderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	form, diags := folderFormFromModel(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := r.client.CreateFolder(ctx, form)
	if err != nil {
		resp.Diagnostics.AddError("Create folder failed", err.Error())
		return
	}
	if !plan.IsExpanded.IsNull() && !plan.IsExpanded.IsUnknown() {
		created, err = r.client.UpdateFolderExpanded(ctx, created.ID, plan.IsExpanded.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Update folder expanded failed", err.Error())
			return
		}
	}
	state := folderResponseToModel(created)
	state.DeleteContents = plan.DeleteContents
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *folderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing folders.")
		return
	}
	var state folderModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	current, err := r.client.GetFolder(ctx, state.ID.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read folder failed", err.Error())
		return
	}
	updated := folderResponseToModel(current)
	updated.DeleteContents = state.DeleteContents
	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}
func (r *folderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing folders.")
		return
	}
	var plan folderModel
	var state folderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data := decodeOptionalJSON(plan.DataJSON, path.Root("data_json"), &resp.Diagnostics)
	meta := decodeOptionalJSON(plan.MetaJSON, path.Root("meta_json"), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	name := plan.Name.ValueString()
	current, err := r.client.UpdateFolder(ctx, state.ID.ValueString(), client.FolderUpdateForm{Name: &name, Data: data, Meta: meta})
	if err != nil {
		resp.Diagnostics.AddError("Update folder failed", err.Error())
		return
	}
	var parentID *string
	if !plan.ParentID.IsNull() && !plan.ParentID.IsUnknown() {
		v := plan.ParentID.ValueString()
		parentID = &v
	}
	current, err = r.client.UpdateFolderParent(ctx, state.ID.ValueString(), parentID)
	if err != nil {
		resp.Diagnostics.AddError("Update folder parent failed", err.Error())
		return
	}
	if !plan.IsExpanded.IsNull() && !plan.IsExpanded.IsUnknown() {
		current, err = r.client.UpdateFolderExpanded(ctx, state.ID.ValueString(), plan.IsExpanded.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Update folder expanded failed", err.Error())
			return
		}
	}
	updated := folderResponseToModel(current)
	updated.DeleteContents = plan.DeleteContents
	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}
func (r *folderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing folders.")
		return
	}
	var state folderModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteContents := true
	if !state.DeleteContents.IsNull() && !state.DeleteContents.IsUnknown() {
		deleteContents = state.DeleteContents.ValueBool()
	}
	if err := r.client.DeleteFolder(ctx, state.ID.ValueString(), deleteContents); err != nil && err != client.ErrNotFound {
		resp.Diagnostics.AddError("Delete folder failed", err.Error())
	}
}
func (r *folderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func folderFormFromModel(m folderModel) (client.FolderForm, diag.Diagnostics) {
	var diags diag.Diagnostics
	data := decodeOptionalJSON(m.DataJSON, path.Root("data_json"), &diags)
	meta := decodeOptionalJSON(m.MetaJSON, path.Root("meta_json"), &diags)
	var parentID *string
	if !m.ParentID.IsNull() && !m.ParentID.IsUnknown() {
		v := m.ParentID.ValueString()
		parentID = &v
	}
	return client.FolderForm{Name: m.Name.ValueString(), Data: data, Meta: meta, ParentID: parentID}, diags
}
func folderResponseToModel(f *client.FolderModel) folderModel {
	dataJSON, _ := encodeOptionalJSON(f.Data)
	metaJSON, _ := encodeOptionalJSON(f.Meta)
	m := folderModel{ID: types.StringValue(f.ID), Name: types.StringValue(f.Name), ParentID: types.StringNull(), DataJSON: dataJSON, MetaJSON: metaJSON, IsExpanded: types.BoolValue(f.IsExpanded), DeleteContents: types.BoolNull(), UserID: types.StringValue(f.UserID), CreatedAt: types.Int64Value(f.CreatedAt), UpdatedAt: types.Int64Value(f.UpdatedAt)}
	if f.ParentID != nil {
		m.ParentID = types.StringValue(*f.ParentID)
	}
	return m
}
