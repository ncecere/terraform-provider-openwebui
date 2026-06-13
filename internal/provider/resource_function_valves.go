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

var _ resource.Resource = &functionValvesResource{}
var _ resource.ResourceWithConfigure = &functionValvesResource{}
var _ resource.ResourceWithImportState = &functionValvesResource{}

type functionValvesResource struct{ client *client.Client }
type functionValvesModel struct {
	ID         types.String `tfsdk:"id"`
	FunctionID types.String `tfsdk:"function_id"`
	ValvesJSON types.String `tfsdk:"valves_json"`
	SpecJSON   types.String `tfsdk:"spec_json"`
}

func NewFunctionValvesResource() resource.Resource { return &functionValvesResource{} }
func (r *functionValvesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function_valves"
}
func (r *functionValvesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}}, "function_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}}, "valves_json": schema.StringAttribute{Optional: true, Computed: true}, "spec_json": schema.StringAttribute{Computed: true}}}
}
func (r *functionValvesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		r.client = c
	}
}
func (r *functionValvesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.apply(ctx, req.Plan, &resp.Diagnostics, func(s functionValvesModel) { resp.Diagnostics.Append(resp.State.Set(ctx, &s)...) })
}
func (r *functionValvesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state functionValvesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState, diags := readFunctionValves(ctx, r.client, state.FunctionID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
func (r *functionValvesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.apply(ctx, req.Plan, &resp.Diagnostics, func(s functionValvesModel) { resp.Diagnostics.Append(resp.State.Set(ctx, &s)...) })
}
func (r *functionValvesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
func (r *functionValvesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("function_id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
func (r *functionValvesResource) apply(ctx context.Context, getter interface {
	Get(context.Context, any) diag.Diagnostics
}, diags *diag.Diagnostics, set func(functionValvesModel)) {
	if r.client == nil {
		diags.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing function valves.")
		return
	}
	var plan functionValvesModel
	diags.Append(getter.Get(ctx, &plan)...)
	if diags.HasError() {
		return
	}
	valves := decodeOptionalJSON(plan.ValvesJSON, path.Root("valves_json"), diags)
	if diags.HasError() {
		return
	}
	if valves == nil {
		valves = map[string]any{}
	}
	if _, err := r.client.UpdateFunctionValves(ctx, plan.FunctionID.ValueString(), valves); err != nil {
		diags.AddError("Update function valves failed", err.Error())
		return
	}
	state, d := readFunctionValves(ctx, r.client, plan.FunctionID.ValueString())
	diags.Append(d...)
	set(state)
}
func readFunctionValves(ctx context.Context, c *client.Client, id string) (functionValvesModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	valves, err := c.GetFunctionValves(ctx, id)
	if err != nil {
		diags.AddError("Read function valves failed", err.Error())
		return functionValvesModel{}, diags
	}
	spec, err := c.GetFunctionValvesSpec(ctx, id)
	if err != nil {
		spec = map[string]any{}
	}
	valvesJSON, _ := encodeOptionalJSON(valves)
	specJSON, _ := encodeOptionalJSON(spec)
	return functionValvesModel{ID: types.StringValue(id), FunctionID: types.StringValue(id), ValvesJSON: valvesJSON, SpecJSON: specJSON}, diags
}
