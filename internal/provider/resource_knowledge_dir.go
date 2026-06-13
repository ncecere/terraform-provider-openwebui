package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ resource.Resource = &knowledgeDirResource{}
var _ resource.ResourceWithConfigure = &knowledgeDirResource{}
var _ resource.ResourceWithImportState = &knowledgeDirResource{}

type knowledgeDirResource struct{ client *client.Client }

type knowledgeDirModel struct {
	ID          types.String `tfsdk:"id"`
	KnowledgeID types.String `tfsdk:"knowledge_id"`
	DirectoryID types.String `tfsdk:"directory_id"`
	Name        types.String `tfsdk:"name"`
	ParentID    types.String `tfsdk:"parent_id"`
	MoveFiles   types.Bool   `tfsdk:"move_files_on_delete"`
	UserID      types.String `tfsdk:"user_id"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

func NewKnowledgeDirResource() resource.Resource { return &knowledgeDirResource{} }
func (r *knowledgeDirResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge_dir"
}
func (r *knowledgeDirResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":                   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"knowledge_id":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Knowledge base ID."},
		"directory_id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Knowledge directory ID."},
		"name":                 schema.StringAttribute{Required: true, Description: "Directory name."},
		"parent_id":            schema.StringAttribute{Optional: true, Description: "Parent directory ID."},
		"move_files_on_delete": schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}, Description: "Whether files should be moved to the parent directory on delete. Defaults to true."},
		"user_id":              schema.StringAttribute{Computed: true},
		"created_at":           schema.Int64Attribute{Computed: true},
		"updated_at":           schema.Int64Attribute{Computed: true},
	}}
}
func (r *knowledgeDirResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}
func (r *knowledgeDirResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge directories.")
		return
	}
	var plan knowledgeDirModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var parentID *string
	if !plan.ParentID.IsNull() && !plan.ParentID.IsUnknown() {
		v := plan.ParentID.ValueString()
		parentID = &v
	}
	created, err := r.client.CreateKnowledgeDirectory(ctx, plan.KnowledgeID.ValueString(), plan.Name.ValueString(), parentID)
	if err != nil {
		resp.Diagnostics.AddError("Create knowledge directory failed", err.Error())
		return
	}
	state := knowledgeDirResponseToModel(created)
	state.MoveFiles = boolDefault(plan.MoveFiles, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *knowledgeDirResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge directories.")
		return
	}
	var state knowledgeDirModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if _, err := r.client.GetKnowledge(ctx, state.KnowledgeID.ValueString()); err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read knowledge failed", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *knowledgeDirResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge directories.")
		return
	}
	var plan knowledgeDirModel
	var state knowledgeDirModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var parentID *string
	if !plan.ParentID.IsNull() && !plan.ParentID.IsUnknown() {
		v := plan.ParentID.ValueString()
		parentID = &v
	}
	updated, err := r.client.UpdateKnowledgeDirectory(ctx, state.KnowledgeID.ValueString(), state.DirectoryID.ValueString(), plan.Name.ValueString(), parentID)
	if err != nil {
		resp.Diagnostics.AddError("Update knowledge directory failed", err.Error())
		return
	}
	newState := knowledgeDirResponseToModel(updated)
	newState.MoveFiles = boolDefault(plan.MoveFiles, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
func (r *knowledgeDirResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge directories.")
		return
	}
	var state knowledgeDirModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	moveFiles := true
	if !state.MoveFiles.IsNull() && !state.MoveFiles.IsUnknown() {
		moveFiles = state.MoveFiles.ValueBool()
	}
	if err := r.client.DeleteKnowledgeDirectory(ctx, state.KnowledgeID.ValueString(), state.DirectoryID.ValueString(), moveFiles); err != nil && err != client.ErrNotFound {
		resp.Diagnostics.AddError("Delete knowledge directory failed", err.Error())
	}
}
func (r *knowledgeDirResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Use knowledge_id:directory_id.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("knowledge_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_id"), parts[1])...)
}
func knowledgeDirResponseToModel(d *client.KnowledgeDirectoryModel) knowledgeDirModel {
	m := knowledgeDirModel{ID: types.StringValue(d.KnowledgeID + ":" + d.ID), KnowledgeID: types.StringValue(d.KnowledgeID), DirectoryID: types.StringValue(d.ID), Name: types.StringValue(d.Name), ParentID: types.StringNull(), MoveFiles: types.BoolValue(true), UserID: types.StringValue(d.UserID), CreatedAt: types.Int64Value(d.CreatedAt), UpdatedAt: types.Int64Value(d.UpdatedAt)}
	if d.ParentID != nil {
		m.ParentID = types.StringValue(*d.ParentID)
	}
	return m
}
func boolDefault(v types.Bool, def bool) types.Bool {
	if v.IsNull() || v.IsUnknown() {
		return types.BoolValue(def)
	}
	return v
}
