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

type rawConfigResource struct {
	client   *client.Client
	typeName string
	id       string
	desc     string
	apply    func(context.Context, *client.Client, map[string]any) (map[string]any, error)
}

type rawConfigResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ConfigJSON types.String `tfsdk:"config_json"`
}

func newRawConfigResource(typeName, id, desc string, apply func(context.Context, *client.Client, map[string]any) (map[string]any, error)) resource.Resource {
	return &rawConfigResource{typeName: typeName, id: id, desc: desc, apply: apply}
}

func NewRetrievalConfigResource() resource.Resource {
	return newRawConfigResource("retrieval_config", "retrieval", "Applies Open WebUI retrieval configuration from raw JSON.", func(ctx context.Context, c *client.Client, config map[string]any) (map[string]any, error) {
		return c.SetRetrievalConfig(ctx, config)
	})
}
func NewEvaluationsConfigResource() resource.Resource {
	return newRawConfigResource("evaluations_config", "evaluations", "Applies Open WebUI evaluations configuration from raw JSON.", func(ctx context.Context, c *client.Client, config map[string]any) (map[string]any, error) {
		return c.SetEvaluationsConfig(ctx, config)
	})
}
func NewDefaultUserPermissionsResource() resource.Resource {
	return newRawConfigResource("default_user_permissions", "default_user_permissions", "Applies Open WebUI default user permissions from raw JSON.", func(ctx context.Context, c *client.Client, config map[string]any) (map[string]any, error) {
		return c.SetDefaultUserPermissions(ctx, config)
	})
}

func (r *rawConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}
func (r *rawConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: r.desc, Attributes: map[string]schema.Attribute{
		"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"config_json": schema.StringAttribute{Required: true, Description: "Raw JSON configuration payload to apply. The provider stores the desired JSON in state to avoid drift from server-added defaults."},
	}}
}
func (r *rawConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		r.client = c
	}
}
func (r *rawConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.applyPlan(ctx, req.Plan, &resp.Diagnostics, func(state rawConfigResourceModel) { resp.Diagnostics.Append(resp.State.Set(ctx, &state)...) })
}
func (r *rawConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rawConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.ID.IsNull() || state.ID.IsUnknown() {
		state.ID = types.StringValue(r.id)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *rawConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.applyPlan(ctx, req.Plan, &resp.Diagnostics, func(state rawConfigResourceModel) { resp.Diagnostics.Append(resp.State.Set(ctx, &state)...) })
}
func (r *rawConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
func (r *rawConfigResource) applyPlan(ctx context.Context, planGetter interface {
	Get(context.Context, any) diag.Diagnostics
}, diags *diag.Diagnostics, set func(rawConfigResourceModel)) {
	if r.client == nil {
		diags.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing config.")
		return
	}
	var plan rawConfigResourceModel
	diags.Append(planGetter.Get(ctx, &plan)...)
	if diags.HasError() {
		return
	}
	config := decodeOptionalJSON(plan.ConfigJSON, path.Root("config_json"), diags)
	if diags.HasError() {
		return
	}
	if _, err := r.apply(ctx, r.client, config); err != nil {
		diags.AddError("Apply config failed", err.Error())
		return
	}
	set(rawConfigResourceModel{ID: types.StringValue(r.id), ConfigJSON: plan.ConfigJSON})
}
