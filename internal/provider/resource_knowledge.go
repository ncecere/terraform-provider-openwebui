package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ resource.Resource = &knowledgeResource{}
var _ resource.ResourceWithConfigure = &knowledgeResource{}
var _ resource.ResourceWithImportState = &knowledgeResource{}

// knowledgeResource implements the Terraform resource for Open WebUI knowledge bases.
type knowledgeResource struct {
	client *client.Client
}

// knowledgeResourceModel maps the resource schema data.
type knowledgeResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	DataJSON    types.String `tfsdk:"data_json"`
	MetaJSON    types.String `tfsdk:"meta_json"`
	ReadGroups  types.List   `tfsdk:"read_groups"`
	WriteGroups types.List   `tfsdk:"write_groups"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	UserID      types.String `tfsdk:"user_id"`
}

// NewKnowledgeResource returns a new instance.
func NewKnowledgeResource() resource.Resource {
	return &knowledgeResource{}
}

// Metadata implements resource.Resource.
func (r *knowledgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

// Schema describes the resource schema.
func (r *knowledgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Unique identifier for the knowledge base entry.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Human-readable name for the knowledge entry.",
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "Detailed description of the knowledge entry.",
			},
			"data_json": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "JSON payload describing additional data for the knowledge entry.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"meta_json": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "JSON payload describing metadata for the knowledge entry.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"read_groups": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Group names or IDs granted read access to the knowledge entry.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"write_groups": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Group names or IDs granted write access to the knowledge entry.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"created_at": schema.StringAttribute{
				Computed:      true,
				Description:   "Creation date in YYYY-MM-DD format.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.StringAttribute{
				Computed:      true,
				Description:   "Last update date in YYYY-MM-DD format.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the user who owns the knowledge entry.",
			},
		},
	}
}

// Configure receives provider configuration data.
func (r *knowledgeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}

// Create handles the creation of the resource.
func (r *knowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge resources.")
		return
	}

	var plan knowledgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	form := client.KnowledgeForm{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	readNames := expandStringList(ctx, plan.ReadGroups, path.Root("read_groups"), &resp.Diagnostics)
	writeNames := expandStringList(ctx, plan.WriteGroups, path.Root("write_groups"), &resp.Diagnostics)
	readIDs := resolveGroupNamesToIDs(ctx, r.client, readNames, path.Root("read_groups"), &resp.Diagnostics)
	writeIDs := resolveGroupNamesToIDs(ctx, r.client, writeNames, path.Root("write_groups"), &resp.Diagnostics)

	form.AccessControl = buildAccessControl(readIDs, writeIDs)
	form.Data = decodeOptionalJSON(plan.DataJSON, path.Root("data_json"), &resp.Diagnostics)
	form.Meta = decodeOptionalJSON(plan.MetaJSON, path.Root("meta_json"), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateKnowledge(ctx, form)
	if err != nil {
		resp.Diagnostics.AddError("Create knowledge entry failed", err.Error())
		return
	}

	// Fetch the latest representation to populate computed fields consistently.
	current, err := r.client.GetKnowledge(ctx, created.ID)
	if err != nil {
		resp.Diagnostics.AddError("Read knowledge entry failed", err.Error())
		return
	}

	state, diags := knowledgeResponseToModel(ctx, r.client, *current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest API data.
func (r *knowledgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge resources.")
		return
	}

	var state knowledgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	current, err := r.client.GetKnowledge(ctx, state.ID.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Read knowledge entry failed", err.Error())
		return
	}

	updated, diags := knowledgeResponseToModel(ctx, r.client, *current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

// Update applies plan changes.
func (r *knowledgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge resources.")
		return
	}

	var plan knowledgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	form := client.KnowledgeForm{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	readNames := expandStringList(ctx, plan.ReadGroups, path.Root("read_groups"), &resp.Diagnostics)
	writeNames := expandStringList(ctx, plan.WriteGroups, path.Root("write_groups"), &resp.Diagnostics)
	readIDs := resolveGroupNamesToIDs(ctx, r.client, readNames, path.Root("read_groups"), &resp.Diagnostics)
	writeIDs := resolveGroupNamesToIDs(ctx, r.client, writeNames, path.Root("write_groups"), &resp.Diagnostics)

	form.AccessControl = buildAccessControl(readIDs, writeIDs)
	form.Data = decodeOptionalJSON(plan.DataJSON, path.Root("data_json"), &resp.Diagnostics)
	form.Meta = decodeOptionalJSON(plan.MetaJSON, path.Root("meta_json"), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateKnowledge(ctx, plan.ID.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError("Update knowledge entry failed", err.Error())
		return
	}

	current, err := r.client.GetKnowledge(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read knowledge entry failed", err.Error())
		return
	}

	state, diags := knowledgeResponseToModel(ctx, r.client, *current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete removes the knowledge resource.
func (r *knowledgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing knowledge resources.")
		return
	}

	var state knowledgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteKnowledge(ctx, state.ID.ValueString()); err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Delete knowledge entry failed", err.Error())
		return
	}
}

// ImportState maps imported IDs to the id attribute.
func (r *knowledgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// knowledgeResponseToModel maps API structures to Terraform state.
func knowledgeResponseToModel(ctx context.Context, apiClient *client.Client, resp client.KnowledgeFilesResponse) (knowledgeResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, err := encodeOptionalJSON(resp.Data)
	if err != nil {
		diags.AddError("Serialize data", err.Error())
	}

	meta, err := encodeOptionalJSON(resp.Meta)
	if err != nil {
		diags.AddError("Serialize metadata", err.Error())
	}

	readIDs := extractGroupIDsFromAccessControl(resp.AccessControl, "read")
	writeIDs := extractGroupIDsFromAccessControl(resp.AccessControl, "write")

	readNames, readDiags := fetchGroupNamesForIDs(ctx, apiClient, readIDs)
	diags.Append(readDiags...)
	writeNames, writeDiags := fetchGroupNamesForIDs(ctx, apiClient, writeIDs)
	diags.Append(writeDiags...)

	readList := types.ListNull(types.StringType)
	if len(readNames) > 0 {
		l, listDiags := types.ListValueFrom(ctx, types.StringType, readNames)
		diags.Append(listDiags...)
		if !listDiags.HasError() {
			readList = l
		}
	}

	writeList := types.ListNull(types.StringType)
	if len(writeNames) > 0 {
		l, listDiags := types.ListValueFrom(ctx, types.StringType, writeNames)
		diags.Append(listDiags...)
		if !listDiags.HasError() {
			writeList = l
		}
	}

	model := knowledgeResourceModel{
		ID:          types.StringValue(resp.ID),
		Name:        types.StringValue(resp.Name),
		Description: types.StringValue(resp.Description),
		DataJSON:    data,
		MetaJSON:    meta,
		ReadGroups:  readList,
		WriteGroups: writeList,
		CreatedAt:   formatDateValue(resp.CreatedAt),
		UpdatedAt:   formatDateValue(resp.UpdatedAt),
		UserID:      types.StringValue(resp.UserID),
	}

	return model, diags
}
