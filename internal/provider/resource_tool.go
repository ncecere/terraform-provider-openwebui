package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ resource.Resource = &toolResource{}
var _ resource.ResourceWithConfigure = &toolResource{}
var _ resource.ResourceWithImportState = &toolResource{}

// toolResource manages Open WebUI tools.
type toolResource struct {
	client *client.Client
}

// toolResourceModel captures Terraform state for tools.
type toolResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ToolID       types.String `tfsdk:"tool_id"`
	Name         types.String `tfsdk:"name"`
	Content      types.String `tfsdk:"content"`
	Description  types.String `tfsdk:"description"`
	ManifestJSON types.String `tfsdk:"manifest_json"`
	ReadGroups   types.List   `tfsdk:"read_groups"`
	WriteGroups  types.List   `tfsdk:"write_groups"`
	SpecsJSON    types.String `tfsdk:"specs_json"`
	UserID       types.String `tfsdk:"user_id"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	UpdatedAt    types.Int64  `tfsdk:"updated_at"`
	WriteAccess  types.Bool   `tfsdk:"write_access"`
}

// NewToolResource constructs a new tool resource.
func NewToolResource() resource.Resource {
	return &toolResource{}
}

// Metadata sets the resource type name.
func (r *toolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool"
}

// Schema defines the tool resource schema.
func (r *toolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Unique identifier assigned by Open WebUI.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tool_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier used when creating the tool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for the tool.",
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "Source content for the tool.",
			},
			"description": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Human-readable tool description.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"manifest_json": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "JSON manifest for the tool.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"read_groups": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Group names or IDs granted read access to the tool.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"write_groups": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Group names or IDs granted write access to the tool.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"specs_json": schema.StringAttribute{
				Computed:    true,
				Description: "Raw JSON specification returned by Open WebUI.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the user who owns the tool.",
			},
			"created_at": schema.Int64Attribute{
				Computed:      true,
				Description:   "Unix timestamp of tool creation.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.Int64Attribute{
				Computed:      true,
				Description:   "Unix timestamp of the last tool update.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"write_access": schema.BoolAttribute{
				Computed:      true,
				Description:   "Whether the current user has write access.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

// Configure assigns the API client.
func (r *toolResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}

// Create provisions a tool.
func (r *toolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing tools.")
		return
	}

	var plan toolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	form, readIDs, writeIDs, diags := toolFormFromPlan(ctx, r.client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateTool(ctx, form)
	if err != nil {
		resp.Diagnostics.AddError("Create tool failed", err.Error())
		return
	}

	if _, err := r.client.UpdateToolAccess(ctx, created.ID, client.BuildGroupGrants(readIDs, writeIDs)); err != nil {
		resp.Diagnostics.AddError("Set tool access failed", err.Error())
		return
	}

	access, err := r.client.GetTool(ctx, created.ID)
	if err != nil {
		resp.Diagnostics.AddError("Read tool failed", err.Error())
		return
	}

	content, specs, fetchDiags := fetchToolContent(ctx, r.client, created.ID)
	resp.Diagnostics.Append(fetchDiags...)

	state, stateDiags := toolResponseToModel(ctx, r.client, access, content, specs, plan.Content)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ToolID = types.StringValue(created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes tool state.
func (r *toolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing tools.")
		return
	}

	var state toolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	access, err := r.client.GetTool(ctx, state.ID.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read tool failed", err.Error())
		return
	}

	content, specs, fetchDiags := fetchToolContent(ctx, r.client, state.ID.ValueString())
	resp.Diagnostics.Append(fetchDiags...)

	updated, diags := toolResponseToModel(ctx, r.client, access, content, specs, state.Content)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

// Update mutates tool properties.
func (r *toolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing tools.")
		return
	}

	var plan toolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	form, readIDs, writeIDs, diags := toolFormFromPlan(ctx, r.client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdateTool(ctx, plan.ID.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError("Update tool failed", err.Error())
		return
	}

	if _, err := r.client.UpdateToolAccess(ctx, plan.ID.ValueString(), client.BuildGroupGrants(readIDs, writeIDs)); err != nil {
		resp.Diagnostics.AddError("Set tool access failed", err.Error())
		return
	}

	refreshed, err := r.client.GetTool(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read tool failed", err.Error())
		return
	}

	access := refreshed
	updated.AccessGrants = refreshed.AccessGrants

	state, stateDiags := toolResponseToModel(ctx, r.client, access, updated.Content, updated.Specs, plan.Content)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete removes a tool.
func (r *toolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing tools.")
		return
	}

	var state toolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteTool(ctx, state.ID.ValueString()); err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Delete tool failed", err.Error())
		return
	}
}

// ImportState maps an import identifier onto the id attribute.
func (r *toolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tool_id"), req.ID)...)
}

func toolFormFromPlan(ctx context.Context, apiClient *client.Client, plan toolResourceModel) (client.ToolForm, []string, []string, diag.Diagnostics) {
	var diags diag.Diagnostics

	manifest := decodeOptionalJSON(plan.ManifestJSON, path.Root("manifest_json"), &diags)

	var description *string
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		value := plan.Description.ValueString()
		description = &value
	}

	meta := client.ToolMeta{
		Description: description,
		Manifest:    manifest,
	}

	readNames := expandStringList(ctx, plan.ReadGroups, path.Root("read_groups"), &diags)
	writeNames := expandStringList(ctx, plan.WriteGroups, path.Root("write_groups"), &diags)
	readIDs := resolveGroupNamesToIDs(ctx, apiClient, readNames, path.Root("read_groups"), &diags)
	writeIDs := resolveGroupNamesToIDs(ctx, apiClient, writeNames, path.Root("write_groups"), &diags)

	return client.ToolForm{
		ID:      plan.ToolID.ValueString(),
		Name:    plan.Name.ValueString(),
		Content: plan.Content.ValueString(),
		Meta:    meta,
	}, readIDs, writeIDs, diags
}

func fetchToolContent(ctx context.Context, apiClient *client.Client, toolID string) (string, []map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	tools, err := apiClient.ExportTools(ctx)
	if err != nil {
		diags.AddWarning("Export tools failed", err.Error())
		return "", nil, diags
	}

	for _, tool := range tools {
		if tool.ID == toolID {
			return tool.Content, tool.Specs, diags
		}
	}

	return "", nil, diags
}

func toolResponseToModel(ctx context.Context, apiClient *client.Client, access *client.ToolAccessResponse, content string, specs []map[string]any, fallbackContent types.String) (toolResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if access == nil {
		diags.AddError("Missing tool response", "Tool details were not returned by the Open WebUI API.")
		return toolResourceModel{}, diags
	}

	readIDs := extractGroupIDsFromGrants(access.AccessGrants, "read")
	writeIDs := extractGroupIDsFromGrants(access.AccessGrants, "write")

	readNames := readIDs
	writeNames := writeIDs

	readList, readListDiags := flattenStringSlice(ctx, readNames)
	diags.Append(readListDiags...)
	writeList, writeListDiags := flattenStringSlice(ctx, writeNames)
	diags.Append(writeListDiags...)

	manifestJSON, err := encodeOptionalJSON(access.Meta.Manifest)
	if err != nil {
		diags.AddError("Serialize manifest", err.Error())
	}

	specsJSON, err := encodeOptionalJSONValue(specs)
	if err != nil {
		diags.AddError("Serialize specs", err.Error())
	}

	description := types.StringNull()
	if access.Meta.Description != nil {
		description = types.StringValue(*access.Meta.Description)
	}

	contentValue := types.StringNull()
	if content != "" {
		contentValue = types.StringValue(content)
	} else if !fallbackContent.IsNull() && !fallbackContent.IsUnknown() {
		contentValue = types.StringValue(fallbackContent.ValueString())
	}

	state := toolResourceModel{
		ID:           types.StringValue(access.ID),
		ToolID:       types.StringValue(access.ID),
		Name:         types.StringValue(access.Name),
		Content:      contentValue,
		Description:  description,
		ManifestJSON: manifestJSON,
		ReadGroups:   readList,
		WriteGroups:  writeList,
		SpecsJSON:    specsJSON,
		UserID:       types.StringValue(access.UserID),
		CreatedAt:    types.Int64Value(access.CreatedAt),
		UpdatedAt:    types.Int64Value(access.UpdatedAt),
		WriteAccess:  types.BoolNull(),
	}

	if access.WriteAccess != nil {
		state.WriteAccess = types.BoolValue(*access.WriteAccess)
	}

	return state, diags
}
