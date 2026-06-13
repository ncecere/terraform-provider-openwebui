package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ resource.Resource = &fileResource{}
var _ resource.ResourceWithConfigure = &fileResource{}
var _ resource.ResourceWithImportState = &fileResource{}

// fileResource manages uploaded files.
type fileResource struct {
	client *client.Client
}

// fileResourceModel maps Terraform state for files.
type fileResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	SourcePath          types.String `tfsdk:"source_path"`
	MetadataJSON        types.String `tfsdk:"metadata_json"`
	Process             types.Bool   `tfsdk:"process"`
	ProcessInBackground types.Bool   `tfsdk:"process_in_background"`
	TargetFilename      types.String `tfsdk:"target_filename"`
	Content             types.String `tfsdk:"content"`
	Filename            types.String `tfsdk:"filename"`
	Hash                types.String `tfsdk:"hash"`
	UserID              types.String `tfsdk:"user_id"`
	DataJSON            types.String `tfsdk:"data_json"`
	MetaJSON            types.String `tfsdk:"meta_json"`
	CreatedAt           types.Int64  `tfsdk:"created_at"`
	UpdatedAt           types.Int64  `tfsdk:"updated_at"`
}

// NewFileResource constructs a new file resource.
func NewFileResource() resource.Resource {
	return &fileResource{}
}

// Metadata sets the resource type name.
func (r *fileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

// Schema defines the file resource schema.
func (r *fileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Unique identifier assigned by Open WebUI.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"source_path": schema.StringAttribute{
				Required:      true,
				Description:   "Local path to the file to upload.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"metadata_json": schema.StringAttribute{
				Optional:      true,
				Description:   "Optional JSON metadata sent during upload.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"process": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether Open WebUI should process the file after upload.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown(), boolplanmodifier.RequiresReplace()},
			},
			"process_in_background": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether file processing should be queued in the background.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown(), boolplanmodifier.RequiresReplace()},
			},
			"target_filename": schema.StringAttribute{
				Optional:    true,
				Description: "Optional filename to set after upload. This renames the uploaded file without forcing replacement.",
			},
			"content": schema.StringAttribute{
				Optional:    true,
				Description: "Optional extracted text content to store for the file after upload.",
			},
			"filename": schema.StringAttribute{
				Computed:    true,
				Description: "Filename as stored by Open WebUI.",
			},
			"hash": schema.StringAttribute{
				Computed:    true,
				Description: "Hash returned by Open WebUI for the uploaded file.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the user who owns the file.",
			},
			"data_json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON data payload returned by Open WebUI.",
			},
			"meta_json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON metadata returned by Open WebUI.",
			},
			"created_at": schema.Int64Attribute{
				Computed:      true,
				Description:   "Unix timestamp of file creation.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.Int64Attribute{
				Computed:      true,
				Description:   "Unix timestamp of last update.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
		},
	}
}

// Configure assigns the API client.
func (r *fileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}

// Create uploads a new file.
func (r *fileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing files.")
		return
	}

	var plan fileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	process := true
	if !plan.Process.IsNull() && !plan.Process.IsUnknown() {
		process = plan.Process.ValueBool()
	}
	processInBackground := true
	if !plan.ProcessInBackground.IsNull() && !plan.ProcessInBackground.IsUnknown() {
		processInBackground = plan.ProcessInBackground.ValueBool()
	}

	metadata := ""
	if !plan.MetadataJSON.IsNull() && !plan.MetadataJSON.IsUnknown() {
		metadata = plan.MetadataJSON.ValueString()
	}

	uploaded, err := r.client.UploadFile(ctx, plan.SourcePath.ValueString(), metadata, process, processInBackground)
	if err != nil {
		resp.Diagnostics.AddError("Upload file failed", err.Error())
		return
	}

	if !plan.TargetFilename.IsNull() && !plan.TargetFilename.IsUnknown() && plan.TargetFilename.ValueString() != "" {
		if err := r.client.RenameFile(ctx, uploaded.ID, plan.TargetFilename.ValueString()); err != nil {
			resp.Diagnostics.AddError("Rename file failed", err.Error())
			return
		}
	}
	if !plan.Content.IsNull() && !plan.Content.IsUnknown() {
		if err := r.client.UpdateFileContent(ctx, uploaded.ID, plan.Content.ValueString()); err != nil {
			resp.Diagnostics.AddError("Update file content failed", err.Error())
			return
		}
	}

	state, diags := fileStateFromAPI(ctx, r.client, uploaded.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.SourcePath = plan.SourcePath
	state.MetadataJSON = plan.MetadataJSON
	state.Process = types.BoolValue(process)
	state.ProcessInBackground = types.BoolValue(processInBackground)
	state.TargetFilename = plan.TargetFilename
	state.Content = plan.Content
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes file state.
func (r *fileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing files.")
		return
	}

	var state fileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, diags := fileStateFromAPI(ctx, r.client, state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if updated.ID.IsNull() || updated.ID.IsUnknown() || updated.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	updated.SourcePath = state.SourcePath
	updated.MetadataJSON = state.MetadataJSON
	updated.Process = state.Process
	updated.ProcessInBackground = state.ProcessInBackground
	updated.TargetFilename = state.TargetFilename
	updated.Content = state.Content
	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

// Update mutates editable file metadata/content. Upload inputs still require replacement.
func (r *fileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing files.")
		return
	}

	var plan fileResourceModel
	var state fileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if !plan.TargetFilename.IsNull() && !plan.TargetFilename.IsUnknown() && plan.TargetFilename.ValueString() != "" && !plan.TargetFilename.Equal(state.TargetFilename) {
		if err := r.client.RenameFile(ctx, id, plan.TargetFilename.ValueString()); err != nil {
			resp.Diagnostics.AddError("Rename file failed", err.Error())
			return
		}
	}
	if !plan.Content.IsNull() && !plan.Content.IsUnknown() && !plan.Content.Equal(state.Content) {
		if err := r.client.UpdateFileContent(ctx, id, plan.Content.ValueString()); err != nil {
			resp.Diagnostics.AddError("Update file content failed", err.Error())
			return
		}
	}

	updated, diags := fileStateFromAPI(ctx, r.client, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updated.SourcePath = state.SourcePath
	updated.MetadataJSON = state.MetadataJSON
	updated.Process = state.Process
	updated.ProcessInBackground = state.ProcessInBackground
	updated.TargetFilename = plan.TargetFilename
	updated.Content = plan.Content
	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

// Delete removes the file.
func (r *fileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing files.")
		return
	}

	var state fileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteFile(ctx, state.ID.ValueString()); err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Delete file failed", err.Error())
		return
	}
}

// ImportState maps import identifiers to the id attribute.
func (r *fileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func fileStateFromAPI(ctx context.Context, apiClient *client.Client, id string) (fileResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	file, err := apiClient.GetFile(ctx, id)
	if err != nil {
		if err == client.ErrNotFound {
			return fileResourceModel{}, diags
		}
		diags.AddError("Read file failed", err.Error())
		return fileResourceModel{}, diags
	}

	dataJSON, err := encodeOptionalJSON(file.Data)
	if err != nil {
		diags.AddError("Serialize file data", err.Error())
	}
	metaJSON, err := encodeOptionalJSON(file.Meta)
	if err != nil {
		diags.AddError("Serialize file metadata", err.Error())
	}

	createdAt := types.Int64Value(file.CreatedAt)
	updatedAt := types.Int64Value(file.UpdatedAt)

	hash := types.StringNull()
	if file.Hash != nil {
		hash = types.StringValue(*file.Hash)
	}

	state := fileResourceModel{
		ID:        types.StringValue(file.ID),
		Filename:  types.StringValue(file.Filename),
		Hash:      hash,
		UserID:    types.StringValue(file.UserID),
		DataJSON:  dataJSON,
		MetaJSON:  metaJSON,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return state, diags
}
