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

var _ resource.Resource = &modelsConfigResource{}
var _ resource.ResourceWithConfigure = &modelsConfigResource{}
var _ resource.ResourceWithImportState = &modelsConfigResource{}

// modelsConfigResource manages default model configuration.
type modelsConfigResource struct {
	client *client.Client
}

type modelsConfigModel struct {
	ID                  types.String `tfsdk:"id"`
	DefaultModels       types.String `tfsdk:"default_models"`
	DefaultPinnedModels types.String `tfsdk:"default_pinned_models"`
	ModelOrderList      types.List   `tfsdk:"model_order_list"`
}

// NewModelsConfigResource constructs a new models config resource.
func NewModelsConfigResource() resource.Resource {
	return &modelsConfigResource{}
}

// Metadata sets the resource type name.
func (r *modelsConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_models_config"
}

// Schema defines the models config schema.
func (r *modelsConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Singleton identifier for the models config.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"default_models": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Default model IDs.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"default_pinned_models": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Default pinned model IDs.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"model_order_list": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Ordered list of model IDs.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

// Configure assigns the API client.
func (r *modelsConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}

// Create updates the models config.
func (r *modelsConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing models config.")
		return
	}

	var plan modelsConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diags := applyModelsConfig(ctx, r.client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the models config.
func (r *modelsConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing models config.")
		return
	}

	config, err := r.client.GetModelsConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Read models config failed", err.Error())
		return
	}

	orderList, listDiags := flattenStringSlice(ctx, config.ModelOrderList)
	resp.Diagnostics.Append(listDiags...)

	state := modelsConfigModel{
		ID:                  types.StringValue("models"),
		DefaultModels:       stringValueOrNull(config.DefaultModels),
		DefaultPinnedModels: stringValueOrNull(config.DefaultPinnedModels),
		ModelOrderList:      orderList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update applies new models config.
func (r *modelsConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing models config.")
		return
	}

	var plan modelsConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diags := applyModelsConfig(ctx, r.client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete removes the resource from state without changing remote configuration.
func (r *modelsConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing models config.")
		return
	}
}

// ImportState maps import identifiers onto the id attribute.
func (r *modelsConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func applyModelsConfig(ctx context.Context, apiClient *client.Client, plan modelsConfigModel) (modelsConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	var orderList []string
	if !plan.ModelOrderList.IsNull() && !plan.ModelOrderList.IsUnknown() {
		orderList = expandStringList(ctx, plan.ModelOrderList, path.Root("model_order_list"), &diags)
	}

	form := client.ModelsConfigForm{
		DefaultModels:       stringPtr(plan.DefaultModels),
		DefaultPinnedModels: stringPtr(plan.DefaultPinnedModels),
		ModelOrderList:      orderList,
	}

	updated, err := apiClient.SetModelsConfig(ctx, form)
	if err != nil {
		diags.AddError("Update models config failed", err.Error())
		return modelsConfigModel{}, diags
	}

	order, listDiags := flattenStringSlice(ctx, updated.ModelOrderList)
	diags.Append(listDiags...)
	if len(updated.ModelOrderList) == 0 && !plan.ModelOrderList.IsNull() && !plan.ModelOrderList.IsUnknown() {
		order = plan.ModelOrderList
	}

	state := modelsConfigModel{
		ID:                  types.StringValue("models"),
		DefaultModels:       stringValueOrNull(updated.DefaultModels),
		DefaultPinnedModels: stringValueOrNull(updated.DefaultPinnedModels),
		ModelOrderList:      order,
	}

	return state, diags
}
