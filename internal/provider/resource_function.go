package provider

import (
	"context"

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

var _ resource.Resource = &functionResource{}
var _ resource.ResourceWithConfigure = &functionResource{}
var _ resource.ResourceWithImportState = &functionResource{}

type functionResource struct{ client *client.Client }

type functionResourceModel struct {
	ID           types.String `tfsdk:"id"`
	FunctionID   types.String `tfsdk:"function_id"`
	Name         types.String `tfsdk:"name"`
	Content      types.String `tfsdk:"content"`
	Description  types.String `tfsdk:"description"`
	ManifestJSON types.String `tfsdk:"manifest_json"`
	Type         types.String `tfsdk:"type"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	IsGlobal     types.Bool   `tfsdk:"is_global"`
	UserID       types.String `tfsdk:"user_id"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	UpdatedAt    types.Int64  `tfsdk:"updated_at"`
}

func NewFunctionResource() resource.Resource { return &functionResource{} }
func (r *functionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function"
}
func (r *functionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"function_id":   schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Function identifier."},
		"name":          schema.StringAttribute{Required: true},
		"content":       schema.StringAttribute{Required: true},
		"description":   schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"manifest_json": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"type":          schema.StringAttribute{Computed: true},
		"is_active":     schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}},
		"is_global":     schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}},
		"user_id":       schema.StringAttribute{Computed: true},
		"created_at":    schema.Int64Attribute{Computed: true},
		"updated_at":    schema.Int64Attribute{Computed: true},
	}}
}
func (r *functionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		r.client = c
	}
}
func (r *functionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing functions.")
		return
	}
	var plan functionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	form := functionFormFromModel(plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := r.client.CreateFunction(ctx, form)
	if err != nil {
		resp.Diagnostics.AddError("Create function failed", err.Error())
		return
	}
	state, diags := functionResponseToModel(created)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state = r.ensureToggles(ctx, state, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		state.Description = plan.Description
	}
	if !plan.ManifestJSON.IsNull() && !plan.ManifestJSON.IsUnknown() {
		state.ManifestJSON = plan.ManifestJSON
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *functionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing functions.")
		return
	}
	var state functionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cur, err := r.client.GetFunction(ctx, state.FunctionID.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read function failed", err.Error())
		return
	}
	updated, diags := functionResponseToModel(cur)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}
func (r *functionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing functions.")
		return
	}
	var plan functionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	form := functionFormFromModel(plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	updated, err := r.client.UpdateFunction(ctx, plan.FunctionID.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError("Update function failed", err.Error())
		return
	}
	state, diags := functionResponseToModel(updated)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state = r.ensureToggles(ctx, state, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		state.Description = plan.Description
	}
	if !plan.ManifestJSON.IsNull() && !plan.ManifestJSON.IsUnknown() {
		state.ManifestJSON = plan.ManifestJSON
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *functionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing functions.")
		return
	}
	var state functionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteFunction(ctx, state.FunctionID.ValueString()); err != nil && err != client.ErrNotFound {
		resp.Diagnostics.AddError("Delete function failed", err.Error())
	}
}
func (r *functionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("function_id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
func (r *functionResource) ensureToggles(ctx context.Context, state functionResourceModel, plan functionResourceModel, diags *diag.Diagnostics) functionResourceModel {
	if !plan.IsActive.IsNull() && !plan.IsActive.IsUnknown() && !state.IsActive.IsNull() && state.IsActive.ValueBool() != plan.IsActive.ValueBool() {
		cur, err := r.client.ToggleFunction(ctx, state.FunctionID.ValueString())
		if err != nil {
			diags.AddError("Toggle function failed", err.Error())
			return state
		}
		state, _ = functionResponseToModel(cur)
	}
	if !plan.IsGlobal.IsNull() && !plan.IsGlobal.IsUnknown() && !state.IsGlobal.IsNull() && state.IsGlobal.ValueBool() != plan.IsGlobal.ValueBool() {
		cur, err := r.client.ToggleFunctionGlobal(ctx, state.FunctionID.ValueString())
		if err != nil {
			diags.AddError("Toggle function global failed", err.Error())
			return state
		}
		state, _ = functionResponseToModel(cur)
	}
	return state
}
func functionFormFromModel(m functionResourceModel, diags *diag.Diagnostics) client.FunctionForm {
	manifest := decodeOptionalJSON(m.ManifestJSON, path.Root("manifest_json"), diags)
	meta := client.FunctionMeta{Manifest: manifest}
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		v := m.Description.ValueString()
		meta.Description = &v
	}
	return client.FunctionForm{ID: m.FunctionID.ValueString(), Name: m.Name.ValueString(), Content: m.Content.ValueString(), Meta: meta}
}
func functionResponseToModel(f *client.FunctionModel) (functionResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	manifestJSON, err := encodeOptionalJSON(f.Meta.Manifest)
	if err != nil {
		diags.AddError("Serialize function manifest", err.Error())
	}
	desc := types.StringNull()
	if f.Meta.Description != nil {
		desc = types.StringValue(*f.Meta.Description)
	}
	return functionResourceModel{ID: types.StringValue(f.ID), FunctionID: types.StringValue(f.ID), Name: types.StringValue(f.Name), Content: types.StringValue(f.Content), Description: desc, ManifestJSON: manifestJSON, Type: types.StringValue(f.Type), IsActive: types.BoolValue(f.IsActive), IsGlobal: types.BoolValue(f.IsGlobal), UserID: types.StringValue(f.UserID), CreatedAt: types.Int64Value(f.CreatedAt), UpdatedAt: types.Int64Value(f.UpdatedAt)}, diags
}
