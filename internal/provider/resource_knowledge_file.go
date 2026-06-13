package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ resource.Resource = &knowledgeFileResource{}
var _ resource.ResourceWithConfigure = &knowledgeFileResource{}
var _ resource.ResourceWithImportState = &knowledgeFileResource{}

// knowledgeFileResource manages knowledge base file attachments.
type knowledgeFileResource struct {
	client *client.Client
}

// knowledgeFileResourceModel captures Terraform state.
type knowledgeFileResourceModel struct {
	ID          types.String `tfsdk:"id"`
	KnowledgeID types.String `tfsdk:"knowledge_id"`
	FileID      types.String `tfsdk:"file_id"`
	DirectoryID types.String `tfsdk:"directory_id"`
	DeleteFile  types.Bool   `tfsdk:"delete_file"`
	FileJSON    types.String `tfsdk:"file_json"`
}

// NewKnowledgeFileResource constructs a new knowledge file resource.
func NewKnowledgeFileResource() resource.Resource {
	return &knowledgeFileResource{}
}

// Metadata sets the resource type name.
func (r *knowledgeFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge_file"
}

// Schema defines the knowledge file resource schema.
func (r *knowledgeFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Composite identifier in the form knowledge_id:file_id.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"knowledge_id": schema.StringAttribute{
				Required:    true,
				Description: "Knowledge base identifier to attach the file to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_id": schema.StringAttribute{
				Required:    true,
				Description: "File identifier to attach.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"directory_id": schema.StringAttribute{
				Optional:    true,
				Description: "Knowledge directory ID to place the file in. Updating this moves the attachment.",
			},
			"delete_file": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether to delete the file when detaching from the knowledge base.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"file_json": schema.StringAttribute{
				Computed:    true,
				Description: "Raw JSON describing the attached file.",
			},
		},
	}
}

// Configure assigns the API client.
func (r *knowledgeFileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}

// Create attaches a file to a knowledge base.
func (r *knowledgeFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge files.")
		return
	}

	var plan knowledgeFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var directoryID *string
	if !plan.DirectoryID.IsNull() && !plan.DirectoryID.IsUnknown() && plan.DirectoryID.ValueString() != "" {
		value := plan.DirectoryID.ValueString()
		directoryID = &value
	}

	if _, err := r.client.AddKnowledgeFileToDirectory(ctx, plan.KnowledgeID.ValueString(), plan.FileID.ValueString(), directoryID); err != nil {
		resp.Diagnostics.AddError("Attach knowledge file failed", err.Error())
		return
	}

	state, diags := readKnowledgeFileState(ctx, r.client, plan.KnowledgeID.ValueString(), plan.FileID.ValueString(), plan.DirectoryID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.ID.IsNull() || state.ID.IsUnknown() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Attach knowledge file failed", "Open WebUI did not return the attached file in the knowledge listing.")
		return
	}

	deleteFile := true
	if !plan.DeleteFile.IsNull() && !plan.DeleteFile.IsUnknown() {
		deleteFile = plan.DeleteFile.ValueBool()
	}
	state.DeleteFile = types.BoolValue(deleteFile)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the attachment state.
func (r *knowledgeFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge files.")
		return
	}

	var state knowledgeFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, diags := readKnowledgeFileState(ctx, r.client, state.KnowledgeID.ValueString(), state.FileID.ValueString(), state.DirectoryID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if updated.ID.IsNull() || updated.ID.IsUnknown() || updated.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	updated.DeleteFile = state.DeleteFile
	updated.DirectoryID = state.DirectoryID
	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

// Update mutates attachment options such as directory placement.
func (r *knowledgeFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge files.")
		return
	}

	var plan knowledgeFileResourceModel
	var state knowledgeFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.DirectoryID.Equal(state.DirectoryID) {
		var directoryID *string
		if !plan.DirectoryID.IsNull() && !plan.DirectoryID.IsUnknown() && plan.DirectoryID.ValueString() != "" {
			value := plan.DirectoryID.ValueString()
			directoryID = &value
		}
		if _, err := r.client.MoveKnowledgeFile(ctx, state.KnowledgeID.ValueString(), state.FileID.ValueString(), directoryID); err != nil {
			resp.Diagnostics.AddError("Move knowledge file failed", err.Error())
			return
		}
	}

	updated, diags := readKnowledgeFileState(ctx, r.client, state.KnowledgeID.ValueString(), state.FileID.ValueString(), plan.DirectoryID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updated.DeleteFile = plan.DeleteFile
	updated.DirectoryID = plan.DirectoryID
	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

// Delete detaches the file from the knowledge base.
func (r *knowledgeFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge files.")
		return
	}

	var state knowledgeFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteFile := true
	if !state.DeleteFile.IsNull() && !state.DeleteFile.IsUnknown() {
		deleteFile = state.DeleteFile.ValueBool()
	}

	if _, err := r.client.RemoveKnowledgeFile(ctx, state.KnowledgeID.ValueString(), state.FileID.ValueString(), deleteFile); err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Detach knowledge file failed", err.Error())
		return
	}
}

// ImportState maps composite IDs in the form knowledge_id:file_id.
func (r *knowledgeFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) == 2 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("knowledge_id"), parts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("file_id"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func readKnowledgeFileState(ctx context.Context, apiClient *client.Client, knowledgeID string, fileID string, directoryID types.String) (knowledgeFileResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	directoryFilter := ""
	if !directoryID.IsNull() && !directoryID.IsUnknown() {
		directoryFilter = directoryID.ValueString()
	}
	list, err := apiClient.ListKnowledgeFilesFiltered(ctx, knowledgeID, "", false, "", "", "", directoryFilter, 1)
	if err != nil {
		if err == client.ErrNotFound {
			return knowledgeFileResourceModel{}, diags
		}
		diags.AddError("Read knowledge files failed", err.Error())
		return knowledgeFileResourceModel{}, diags
	}

	for _, item := range list.Items {
		if item.ID == fileID {
			fileJSON, err := encodeOptionalJSONValue(item)
			if err != nil {
				diags.AddError("Serialize file attachment", err.Error())
			}
			state := knowledgeFileResourceModel{
				ID:          types.StringValue(knowledgeID + ":" + fileID),
				KnowledgeID: types.StringValue(knowledgeID),
				FileID:      types.StringValue(fileID),
				DirectoryID: directoryID,
				FileJSON:    fileJSON,
			}
			return state, diags
		}
	}

	return knowledgeFileResourceModel{}, diags
}
