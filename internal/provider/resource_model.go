package provider

import (
	"context"
	"fmt"

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

var _ resource.Resource = &modelResource{}
var _ resource.ResourceWithConfigure = &modelResource{}
var _ resource.ResourceWithImportState = &modelResource{}

// modelResource implements the Terraform resource for Open WebUI models.
type modelResource struct {
	client *client.Client
}

// modelResourceModel captures Terraform state and plan data.
type modelResourceModel struct {
	ID                   types.String            `tfsdk:"id"`
	ModelID              types.String            `tfsdk:"model_id"`
	Name                 types.String            `tfsdk:"name"`
	BaseModelID          types.String            `tfsdk:"base_model_id"`
	IsActive             types.Bool              `tfsdk:"is_active"`
	UserID               types.String            `tfsdk:"user_id"`
	CreatedAt            types.Int64             `tfsdk:"created_at"`
	UpdatedAt            types.Int64             `tfsdk:"updated_at"`
	MetaAdditionalJSON   types.String            `tfsdk:"meta_additional_json"`
	Params               *modelParamsModel       `tfsdk:"params"`
	ParamsAdditionalJSON types.String            `tfsdk:"params_additional_json"`
	ReadGroups           types.List              `tfsdk:"read_groups"`
	WriteGroups          types.List              `tfsdk:"write_groups"`
	ProfileImageURL      types.String            `tfsdk:"profile_image_url"`
	Description          types.String            `tfsdk:"description"`
	SuggestionPrompts    types.List              `tfsdk:"suggestion_prompts"`
	Tags                 types.List              `tfsdk:"tags"`
	ToolIDs              types.List              `tfsdk:"tool_ids"`
	DefaultFeatureIDs    types.List              `tfsdk:"default_feature_ids"`
	Capabilities         *modelCapabilitiesModel `tfsdk:"capabilities"`
}

// modelCapabilitiesModel captures boolean capability toggles.
type modelCapabilitiesModel struct {
	Vision          types.Bool `tfsdk:"vision"`
	FileUpload      types.Bool `tfsdk:"file_upload"`
	WebSearch       types.Bool `tfsdk:"web_search"`
	ImageGeneration types.Bool `tfsdk:"image_generation"`
	CodeInterpreter types.Bool `tfsdk:"code_interpreter"`
	Citations       types.Bool `tfsdk:"citations"`
	StatusUpdates   types.Bool `tfsdk:"status_updates"`
	Usage           types.Bool `tfsdk:"usage"`
}

// modelParamsModel captures configurable model parameters.
type modelParamsModel struct {
	System               types.String  `tfsdk:"system"`
	StreamResponse       types.Bool    `tfsdk:"stream_response"`
	StreamDeltaChunkSize types.Int64   `tfsdk:"stream_delta_chunk_size"`
	FunctionCalling      types.String  `tfsdk:"function_calling"`
	ReasoningTags        types.List    `tfsdk:"reasoning_tags"`
	Seed                 types.Int64   `tfsdk:"seed"`
	Temperature          types.Float64 `tfsdk:"temperature"`
	KeepAlive            types.String  `tfsdk:"keep_alive"`
	NumGPU               types.Int64   `tfsdk:"num_gpu"`
	NumThread            types.Int64   `tfsdk:"num_thread"`
	NumBatch             types.Int64   `tfsdk:"num_batch"`
	NumCtx               types.Int64   `tfsdk:"num_ctx"`
	NumKeep              types.Int64   `tfsdk:"num_keep"`
	Format               types.String  `tfsdk:"format"`
	Think                types.Bool    `tfsdk:"think"`
	UseMlock             types.Bool    `tfsdk:"use_mlock"`
	UseMmap              types.Bool    `tfsdk:"use_mmap"`
	RepeatPenalty        types.Float64 `tfsdk:"repeat_penalty"`
	TfsZ                 types.Float64 `tfsdk:"tfs_z"`
	RepeatLastN          types.Int64   `tfsdk:"repeat_last_n"`
	MirostatTau          types.Float64 `tfsdk:"mirostat_tau"`
	MirostatEta          types.Float64 `tfsdk:"mirostat_eta"`
	Mirostat             types.Int64   `tfsdk:"mirostat"`
	PresencePenalty      types.Float64 `tfsdk:"presence_penalty"`
	FrequencyPenalty     types.Float64 `tfsdk:"frequency_penalty"`
	MinP                 types.Float64 `tfsdk:"min_p"`
	TopP                 types.Float64 `tfsdk:"top_p"`
	TopK                 types.Int64   `tfsdk:"top_k"`
	MaxTokens            types.Int64   `tfsdk:"max_tokens"`
	ReasoningEffort      types.String  `tfsdk:"reasoning_effort"`
	CustomParams         types.Map     `tfsdk:"custom_params"`
}

// modelMetaState captures structured metadata fields exposed at the top level.
type modelMetaState struct {
	ProfileImageURL   types.String
	Description       types.String
	SuggestionPrompts types.List
	Tags              types.List
	ToolIDs           types.List
	DefaultFeatureIDs types.List
	Capabilities      *modelCapabilitiesModel
}

// NewModelResource constructs a new instance.
func NewModelResource() resource.Resource {
	return &modelResource{}
}

// Metadata implements resource.Resource.
func (r *modelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

// Schema defines the Terraform schema for the model resource.
func (r *modelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Unique identifier returned by the Open WebUI API for the model.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"model_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier supplied to the Open WebUI API when creating the model.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Human-readable model name.",
			},
			"base_model_id": schema.StringAttribute{
				Optional:    true,
				Description: "Optional base model identifier.",
			},
			"is_active": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Flag indicating whether the model is active.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the user who owns the model.",
			},
			"created_at": schema.Int64Attribute{
				Computed:      true,
				Description:   "Unix timestamp indicating model creation time.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.Int64Attribute{
				Computed:      true,
				Description:   "Unix timestamp indicating the last update time.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"meta_additional_json": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Raw JSON fragment merged into the metadata payload for fields not covered by dedicated arguments.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"params_additional_json": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Raw JSON fragment merged into the params payload for fields not covered by dedicated arguments.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"read_groups": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Group names or IDs granted read access to the model.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"write_groups": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Group names or IDs granted write access to the model.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"profile_image_url": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Profile image URL displayed for the model.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Human-readable description of the model.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"suggestion_prompts": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Prompt suggestions surfaced to end users when selecting the model.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Tags associated with the model.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"tool_ids": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Identifiers of tools made available to the model.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"default_feature_ids": schema.ListAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "Feature identifiers enabled by default for the model.",
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"capabilities": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Capability flags exposed by the model.",
				Attributes: map[string]schema.Attribute{
					"vision": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"file_upload": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"web_search": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"image_generation": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"code_interpreter": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"citations": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"status_updates": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"usage": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
				},
			},
			"params": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Parameter values that control model behaviour.",
				Attributes: map[string]schema.Attribute{
					"system":                  schema.StringAttribute{Optional: true},
					"stream_response":         schema.BoolAttribute{Optional: true},
					"stream_delta_chunk_size": schema.Int64Attribute{Optional: true},
					"function_calling":        schema.StringAttribute{Optional: true},
					"reasoning_tags": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"seed":              schema.Int64Attribute{Optional: true},
					"temperature":       schema.Float64Attribute{Optional: true},
					"keep_alive":        schema.StringAttribute{Optional: true},
					"num_gpu":           schema.Int64Attribute{Optional: true},
					"num_thread":        schema.Int64Attribute{Optional: true},
					"num_batch":         schema.Int64Attribute{Optional: true},
					"num_ctx":           schema.Int64Attribute{Optional: true},
					"num_keep":          schema.Int64Attribute{Optional: true},
					"format":            schema.StringAttribute{Optional: true},
					"think":             schema.BoolAttribute{Optional: true},
					"use_mlock":         schema.BoolAttribute{Optional: true},
					"use_mmap":          schema.BoolAttribute{Optional: true},
					"repeat_penalty":    schema.Float64Attribute{Optional: true},
					"tfs_z":             schema.Float64Attribute{Optional: true},
					"repeat_last_n":     schema.Int64Attribute{Optional: true},
					"mirostat_tau":      schema.Float64Attribute{Optional: true},
					"mirostat_eta":      schema.Float64Attribute{Optional: true},
					"mirostat":          schema.Int64Attribute{Optional: true},
					"presence_penalty":  schema.Float64Attribute{Optional: true},
					"frequency_penalty": schema.Float64Attribute{Optional: true},
					"min_p":             schema.Float64Attribute{Optional: true},
					"top_p":             schema.Float64Attribute{Optional: true},
					"top_k":             schema.Int64Attribute{Optional: true},
					"max_tokens":        schema.Int64Attribute{Optional: true},
					"reasoning_effort":  schema.StringAttribute{Optional: true},
					"custom_params": schema.MapAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Additional key/value parameters passed through to Open WebUI.",
					},
				},
			},
		},
	}
}

// Configure wires the API client into the resource.
func (r *modelResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}

// Create provisions the model through the API.
func (r *modelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing models.")
		return
	}

	var plan modelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Params == nil {
		plan.Params = &modelParamsModel{}
	}

	paramsMap := expandModelParams(ctx, plan.Params, &resp.Diagnostics)
	metaMap := expandModelMeta(ctx, &plan, &resp.Diagnostics)

	additionalParams := decodeOptionalJSON(plan.ParamsAdditionalJSON, path.Root("params_additional_json"), &resp.Diagnostics)
	additionalMeta := decodeOptionalJSON(plan.MetaAdditionalJSON, path.Root("meta_additional_json"), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMap = mergeStringAnyMaps(paramsMap, additionalParams)
	metaMap = mergeStringAnyMaps(metaMap, additionalMeta)

	if paramsMap == nil {
		paramsMap = map[string]any{}
	}
	if metaMap == nil {
		metaMap = map[string]any{}
	}

	form := client.ModelForm{
		ID:     plan.ModelID.ValueString(),
		Name:   plan.Name.ValueString(),
		Meta:   copyStringAnyMap(metaMap),
		Params: copyStringAnyMap(paramsMap),
	}

	readNames := expandStringList(ctx, plan.ReadGroups, path.Root("read_groups"), &resp.Diagnostics)
	writeNames := expandStringList(ctx, plan.WriteGroups, path.Root("write_groups"), &resp.Diagnostics)
	readIDs := resolveGroupNamesToIDs(ctx, r.client, readNames, path.Root("read_groups"), &resp.Diagnostics)
	writeIDs := resolveGroupNamesToIDs(ctx, r.client, writeNames, path.Root("write_groups"), &resp.Diagnostics)
	form.AccessControl = buildAccessControl(readIDs, writeIDs)

	if !plan.BaseModelID.IsNull() && !plan.BaseModelID.IsUnknown() && plan.BaseModelID.ValueString() != "" {
		base := plan.BaseModelID.ValueString()
		form.BaseModelID = &base
	}

	if !plan.IsActive.IsNull() && !plan.IsActive.IsUnknown() {
		value := plan.IsActive.ValueBool()
		form.IsActive = &value
	}

	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateModel(ctx, form)
	if err != nil {
		resp.Diagnostics.AddError("Create model failed", err.Error())
		return
	}

	state, diags := modelResponseToModel(ctx, r.client, created, plan.ModelID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read fetches model details from the API.
func (r *modelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing models.")
		return
	}

	var state modelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	current, err := r.client.GetModel(ctx, state.ID.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read model failed", err.Error())
		return
	}

	updated, diags := modelResponseToModel(ctx, r.client, current, state.ModelID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

// Update mutates the model via the API.
func (r *modelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing models.")
		return
	}

	var plan modelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Params == nil {
		plan.Params = &modelParamsModel{}
	}

	paramsMap := expandModelParams(ctx, plan.Params, &resp.Diagnostics)
	metaMap := expandModelMeta(ctx, &plan, &resp.Diagnostics)

	additionalParams := decodeOptionalJSON(plan.ParamsAdditionalJSON, path.Root("params_additional_json"), &resp.Diagnostics)
	additionalMeta := decodeOptionalJSON(plan.MetaAdditionalJSON, path.Root("meta_additional_json"), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMap = mergeStringAnyMaps(paramsMap, additionalParams)
	metaMap = mergeStringAnyMaps(metaMap, additionalMeta)

	if paramsMap == nil {
		paramsMap = map[string]any{}
	}
	if metaMap == nil {
		metaMap = map[string]any{}
	}

	form := client.ModelForm{
		ID:     plan.ModelID.ValueString(),
		Name:   plan.Name.ValueString(),
		Meta:   copyStringAnyMap(metaMap),
		Params: copyStringAnyMap(paramsMap),
	}

	readNames := expandStringList(ctx, plan.ReadGroups, path.Root("read_groups"), &resp.Diagnostics)
	writeNames := expandStringList(ctx, plan.WriteGroups, path.Root("write_groups"), &resp.Diagnostics)
	readIDs := resolveGroupNamesToIDs(ctx, r.client, readNames, path.Root("read_groups"), &resp.Diagnostics)
	writeIDs := resolveGroupNamesToIDs(ctx, r.client, writeNames, path.Root("write_groups"), &resp.Diagnostics)
	form.AccessControl = buildAccessControl(readIDs, writeIDs)

	if !plan.BaseModelID.IsNull() && !plan.BaseModelID.IsUnknown() && plan.BaseModelID.ValueString() != "" {
		base := plan.BaseModelID.ValueString()
		form.BaseModelID = &base
	}

	if !plan.IsActive.IsNull() && !plan.IsActive.IsUnknown() {
		value := plan.IsActive.ValueBool()
		form.IsActive = &value
	}

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateModel(ctx, plan.ID.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError("Update model failed", err.Error())
		return
	}

	current, err := r.client.GetModel(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read model failed", err.Error())
		return
	}

	state, diags := modelResponseToModel(ctx, r.client, current, plan.ModelID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete removes the model.
func (r *modelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing models.")
		return
	}

	var state modelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteModel(ctx, state.ID.ValueString()); err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Delete model failed", err.Error())
		return
	}
}

// ImportState maps import IDs onto the id attribute.
func (r *modelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("model_id"), req.ID)...)
}

// modelResponseToModel maps API responses to Terraform state structures.
func modelResponseToModel(ctx context.Context, apiClient *client.Client, resp *client.ModelResponse, requestedID string) (modelResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	paramsModel, paramsAdditional, paramsDiags := flattenModelParams(ctx, resp.Params)
	diags.Append(paramsDiags...)
	metaState, metaAdditional, metaDiags := flattenModelMeta(ctx, resp.Meta)
	diags.Append(metaDiags...)

	readIDs := extractGroupIDsFromAccessControl(resp.AccessControl, "read")
	writeIDs := extractGroupIDsFromAccessControl(resp.AccessControl, "write")

	readNames, readDiags := fetchGroupNamesForIDs(ctx, apiClient, readIDs)
	diags.Append(readDiags...)
	writeNames, writeDiags := fetchGroupNamesForIDs(ctx, apiClient, writeIDs)
	diags.Append(writeDiags...)

	readList, readListDiags := flattenStringSlice(ctx, readNames)
	diags.Append(readListDiags...)
	writeList, writeListDiags := flattenStringSlice(ctx, writeNames)
	diags.Append(writeListDiags...)

	state := modelResourceModel{
		ID:                   types.StringValue(resp.ID),
		ModelID:              types.StringValue(requestedID),
		Name:                 types.StringValue(resp.Name),
		IsActive:             types.BoolValue(resp.IsActive),
		UserID:               types.StringValue(resp.UserID),
		CreatedAt:            types.Int64Value(resp.CreatedAt),
		UpdatedAt:            types.Int64Value(resp.UpdatedAt),
		Params:               paramsModel,
		ParamsAdditionalJSON: paramsAdditional,
		MetaAdditionalJSON:   metaAdditional,
		ReadGroups:           readList,
		WriteGroups:          writeList,
		ProfileImageURL:      metaState.ProfileImageURL,
		Description:          metaState.Description,
		SuggestionPrompts:    metaState.SuggestionPrompts,
		Tags:                 metaState.Tags,
		ToolIDs:              metaState.ToolIDs,
		DefaultFeatureIDs:    metaState.DefaultFeatureIDs,
		Capabilities:         metaState.Capabilities,
	}

	if resp.BaseModelID != nil && *resp.BaseModelID != "" {
		state.BaseModelID = types.StringValue(*resp.BaseModelID)
	} else {
		state.BaseModelID = types.StringNull()
	}

	return state, diags
}

func expandModelParams(ctx context.Context, params *modelParamsModel, diags *diag.Diagnostics) map[string]any {
	if params == nil {
		return nil
	}

	result := make(map[string]any)

	if !params.System.IsNull() && !params.System.IsUnknown() {
		result["system"] = params.System.ValueString()
	}
	if !params.StreamResponse.IsNull() && !params.StreamResponse.IsUnknown() {
		result["stream_response"] = params.StreamResponse.ValueBool()
	}
	if !params.StreamDeltaChunkSize.IsNull() && !params.StreamDeltaChunkSize.IsUnknown() {
		result["stream_delta_chunk_size"] = params.StreamDeltaChunkSize.ValueInt64()
	}
	if !params.FunctionCalling.IsNull() && !params.FunctionCalling.IsUnknown() {
		result["function_calling"] = params.FunctionCalling.ValueString()
	}
	if !params.ReasoningTags.IsNull() && !params.ReasoningTags.IsUnknown() {
		var tags []string
		diags.Append(params.ReasoningTags.ElementsAs(ctx, &tags, false)...)
		if len(tags) > 0 {
			result["reasoning_tags"] = tags
		}
	}
	if !params.Seed.IsNull() && !params.Seed.IsUnknown() {
		result["seed"] = params.Seed.ValueInt64()
	}
	if !params.Temperature.IsNull() && !params.Temperature.IsUnknown() {
		result["temperature"] = params.Temperature.ValueFloat64()
	}
	if !params.KeepAlive.IsNull() && !params.KeepAlive.IsUnknown() {
		result["keep_alive"] = params.KeepAlive.ValueString()
	}
	if !params.NumGPU.IsNull() && !params.NumGPU.IsUnknown() {
		result["num_gpu"] = params.NumGPU.ValueInt64()
	}
	if !params.NumThread.IsNull() && !params.NumThread.IsUnknown() {
		result["num_thread"] = params.NumThread.ValueInt64()
	}
	if !params.NumBatch.IsNull() && !params.NumBatch.IsUnknown() {
		result["num_batch"] = params.NumBatch.ValueInt64()
	}
	if !params.NumCtx.IsNull() && !params.NumCtx.IsUnknown() {
		result["num_ctx"] = params.NumCtx.ValueInt64()
	}
	if !params.NumKeep.IsNull() && !params.NumKeep.IsUnknown() {
		result["num_keep"] = params.NumKeep.ValueInt64()
	}
	if !params.Format.IsNull() && !params.Format.IsUnknown() {
		result["format"] = params.Format.ValueString()
	}
	if !params.Think.IsNull() && !params.Think.IsUnknown() {
		result["think"] = params.Think.ValueBool()
	}
	if !params.UseMlock.IsNull() && !params.UseMlock.IsUnknown() {
		result["use_mlock"] = params.UseMlock.ValueBool()
	}
	if !params.UseMmap.IsNull() && !params.UseMmap.IsUnknown() {
		result["use_mmap"] = params.UseMmap.ValueBool()
	}
	if !params.RepeatPenalty.IsNull() && !params.RepeatPenalty.IsUnknown() {
		result["repeat_penalty"] = params.RepeatPenalty.ValueFloat64()
	}
	if !params.TfsZ.IsNull() && !params.TfsZ.IsUnknown() {
		result["tfs_z"] = params.TfsZ.ValueFloat64()
	}
	if !params.RepeatLastN.IsNull() && !params.RepeatLastN.IsUnknown() {
		result["repeat_last_n"] = params.RepeatLastN.ValueInt64()
	}
	if !params.MirostatTau.IsNull() && !params.MirostatTau.IsUnknown() {
		result["mirostat_tau"] = params.MirostatTau.ValueFloat64()
	}
	if !params.MirostatEta.IsNull() && !params.MirostatEta.IsUnknown() {
		result["mirostat_eta"] = params.MirostatEta.ValueFloat64()
	}
	if !params.Mirostat.IsNull() && !params.Mirostat.IsUnknown() {
		result["mirostat"] = params.Mirostat.ValueInt64()
	}
	if !params.PresencePenalty.IsNull() && !params.PresencePenalty.IsUnknown() {
		result["presence_penalty"] = params.PresencePenalty.ValueFloat64()
	}
	if !params.FrequencyPenalty.IsNull() && !params.FrequencyPenalty.IsUnknown() {
		result["frequency_penalty"] = params.FrequencyPenalty.ValueFloat64()
	}
	if !params.MinP.IsNull() && !params.MinP.IsUnknown() {
		result["min_p"] = params.MinP.ValueFloat64()
	}
	if !params.TopP.IsNull() && !params.TopP.IsUnknown() {
		result["top_p"] = params.TopP.ValueFloat64()
	}
	if !params.TopK.IsNull() && !params.TopK.IsUnknown() {
		result["top_k"] = params.TopK.ValueInt64()
	}
	if !params.MaxTokens.IsNull() && !params.MaxTokens.IsUnknown() {
		result["max_tokens"] = params.MaxTokens.ValueInt64()
	}
	if !params.ReasoningEffort.IsNull() && !params.ReasoningEffort.IsUnknown() {
		result["reasoning_effort"] = params.ReasoningEffort.ValueString()
	}
	if !params.CustomParams.IsNull() && !params.CustomParams.IsUnknown() {
		var custom map[string]string
		diags.Append(params.CustomParams.ElementsAs(ctx, &custom, false)...)
		if len(custom) > 0 {
			customAny := make(map[string]any, len(custom))
			for k, v := range custom {
				customAny[k] = v
			}
			result["custom_params"] = customAny
		}
	}

	return result
}

func expandModelMeta(ctx context.Context, plan *modelResourceModel, diags *diag.Diagnostics) map[string]any {
	if plan == nil {
		return nil
	}

	var meta map[string]any

	ensureMap := func() {
		if meta == nil {
			meta = make(map[string]any)
		}
	}

	if !plan.ProfileImageURL.IsNull() && !plan.ProfileImageURL.IsUnknown() {
		ensureMap()
		meta["profile_image_url"] = plan.ProfileImageURL.ValueString()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ensureMap()
		meta["description"] = plan.Description.ValueString()
	}
	if plan.Capabilities != nil {
		if caps := expandModelCapabilities(plan.Capabilities); len(caps) > 0 {
			ensureMap()
			meta["capabilities"] = caps
		}
	}
	if !plan.SuggestionPrompts.IsNull() && !plan.SuggestionPrompts.IsUnknown() {
		prompts := expandStringList(ctx, plan.SuggestionPrompts, path.Root("suggestion_prompts"), diags)
		ensureMap()
		if len(prompts) > 0 {
			items := make([]map[string]any, 0, len(prompts))
			for _, prompt := range prompts {
				items = append(items, map[string]any{"content": prompt})
			}
			meta["suggestion_prompts"] = items
		} else {
			meta["suggestion_prompts"] = []any{}
		}
	}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		tags := expandStringList(ctx, plan.Tags, path.Root("tags"), diags)
		ensureMap()
		if len(tags) > 0 {
			items := make([]map[string]any, 0, len(tags))
			for _, tag := range tags {
				items = append(items, map[string]any{"name": tag})
			}
			meta["tags"] = items
		} else {
			meta["tags"] = []any{}
		}
	}
	if !plan.ToolIDs.IsNull() && !plan.ToolIDs.IsUnknown() {
		tools := expandStringList(ctx, plan.ToolIDs, path.Root("tool_ids"), diags)
		ensureMap()
		if len(tools) > 0 {
			meta["toolIds"] = tools
		} else {
			meta["toolIds"] = []any{}
		}
	}
	if !plan.DefaultFeatureIDs.IsNull() && !plan.DefaultFeatureIDs.IsUnknown() {
		features := expandStringList(ctx, plan.DefaultFeatureIDs, path.Root("default_feature_ids"), diags)
		ensureMap()
		if len(features) > 0 {
			meta["defaultFeatureIds"] = features
		} else {
			meta["defaultFeatureIds"] = []any{}
		}
	}

	return meta
}

func expandModelCapabilities(caps *modelCapabilitiesModel) map[string]any {
	if caps == nil {
		return nil
	}

	result := make(map[string]any)

	if !caps.Vision.IsNull() && !caps.Vision.IsUnknown() {
		result["vision"] = caps.Vision.ValueBool()
	}
	if !caps.FileUpload.IsNull() && !caps.FileUpload.IsUnknown() {
		result["file_upload"] = caps.FileUpload.ValueBool()
	}
	if !caps.WebSearch.IsNull() && !caps.WebSearch.IsUnknown() {
		result["web_search"] = caps.WebSearch.ValueBool()
	}
	if !caps.ImageGeneration.IsNull() && !caps.ImageGeneration.IsUnknown() {
		result["image_generation"] = caps.ImageGeneration.ValueBool()
	}
	if !caps.CodeInterpreter.IsNull() && !caps.CodeInterpreter.IsUnknown() {
		result["code_interpreter"] = caps.CodeInterpreter.ValueBool()
	}
	if !caps.Citations.IsNull() && !caps.Citations.IsUnknown() {
		result["citations"] = caps.Citations.ValueBool()
	}
	if !caps.StatusUpdates.IsNull() && !caps.StatusUpdates.IsUnknown() {
		result["status_updates"] = caps.StatusUpdates.ValueBool()
	}
	if !caps.Usage.IsNull() && !caps.Usage.IsUnknown() {
		result["usage"] = caps.Usage.ValueBool()
	}

	return result
}

func flattenModelParams(ctx context.Context, data map[string]any) (*modelParamsModel, types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &modelParamsModel{
		System:               types.StringNull(),
		StreamResponse:       types.BoolNull(),
		StreamDeltaChunkSize: types.Int64Null(),
		FunctionCalling:      types.StringNull(),
		ReasoningTags:        types.ListNull(types.StringType),
		Seed:                 types.Int64Null(),
		Temperature:          types.Float64Null(),
		KeepAlive:            types.StringNull(),
		NumGPU:               types.Int64Null(),
		NumThread:            types.Int64Null(),
		NumBatch:             types.Int64Null(),
		NumCtx:               types.Int64Null(),
		NumKeep:              types.Int64Null(),
		Format:               types.StringNull(),
		Think:                types.BoolNull(),
		UseMlock:             types.BoolNull(),
		UseMmap:              types.BoolNull(),
		RepeatPenalty:        types.Float64Null(),
		TfsZ:                 types.Float64Null(),
		RepeatLastN:          types.Int64Null(),
		MirostatTau:          types.Float64Null(),
		MirostatEta:          types.Float64Null(),
		Mirostat:             types.Int64Null(),
		PresencePenalty:      types.Float64Null(),
		FrequencyPenalty:     types.Float64Null(),
		MinP:                 types.Float64Null(),
		TopP:                 types.Float64Null(),
		TopK:                 types.Int64Null(),
		MaxTokens:            types.Int64Null(),
		ReasoningEffort:      types.StringNull(),
		CustomParams:         types.MapNull(types.StringType),
	}

	if len(data) == 0 {
		return model, types.StringNull(), diags
	}

	additional := copyStringAnyMap(data)

	if value, ok := toStringValue(data["system"]); ok {
		model.System = types.StringValue(value)
		delete(additional, "system")
	}
	if value, ok := toBoolValue(data["stream_response"]); ok {
		model.StreamResponse = types.BoolValue(value)
		delete(additional, "stream_response")
	}
	if value, ok := toInt64Value(data["stream_delta_chunk_size"]); ok {
		model.StreamDeltaChunkSize = types.Int64Value(value)
		delete(additional, "stream_delta_chunk_size")
	}
	if value, ok := toStringValue(data["function_calling"]); ok {
		model.FunctionCalling = types.StringValue(value)
		delete(additional, "function_calling")
	}
	if raw, ok := data["reasoning_tags"]; ok && raw != nil {
		tags, convOK := toStringSlice(raw)
		if convOK {
			list, listDiags := types.ListValueFrom(ctx, types.StringType, tags)
			diags.Append(listDiags...)
			if !listDiags.HasError() {
				model.ReasoningTags = list
			}
			delete(additional, "reasoning_tags")
		} else {
			diags.AddError("Unexpected params value", fmt.Sprintf("Expected params.reasoning_tags to be a list of strings, received %T", raw))
		}
	}
	if value, ok := toInt64Value(data["seed"]); ok {
		model.Seed = types.Int64Value(value)
		delete(additional, "seed")
	}
	if value, ok := toFloat64Value(data["temperature"]); ok {
		model.Temperature = types.Float64Value(value)
		delete(additional, "temperature")
	}
	if value, ok := toStringValue(data["keep_alive"]); ok {
		model.KeepAlive = types.StringValue(value)
		delete(additional, "keep_alive")
	}
	if value, ok := toInt64Value(data["num_gpu"]); ok {
		model.NumGPU = types.Int64Value(value)
		delete(additional, "num_gpu")
	}
	if value, ok := toInt64Value(data["num_thread"]); ok {
		model.NumThread = types.Int64Value(value)
		delete(additional, "num_thread")
	}
	if value, ok := toInt64Value(data["num_batch"]); ok {
		model.NumBatch = types.Int64Value(value)
		delete(additional, "num_batch")
	}
	if value, ok := toInt64Value(data["num_ctx"]); ok {
		model.NumCtx = types.Int64Value(value)
		delete(additional, "num_ctx")
	}
	if value, ok := toInt64Value(data["num_keep"]); ok {
		model.NumKeep = types.Int64Value(value)
		delete(additional, "num_keep")
	}
	if value, ok := toStringValue(data["format"]); ok {
		model.Format = types.StringValue(value)
		delete(additional, "format")
	}
	if value, ok := toBoolValue(data["think"]); ok {
		model.Think = types.BoolValue(value)
		delete(additional, "think")
	}
	if value, ok := toBoolValue(data["use_mlock"]); ok {
		model.UseMlock = types.BoolValue(value)
		delete(additional, "use_mlock")
	}
	if value, ok := toBoolValue(data["use_mmap"]); ok {
		model.UseMmap = types.BoolValue(value)
		delete(additional, "use_mmap")
	}
	if value, ok := toFloat64Value(data["repeat_penalty"]); ok {
		model.RepeatPenalty = types.Float64Value(value)
		delete(additional, "repeat_penalty")
	}
	if value, ok := toFloat64Value(data["tfs_z"]); ok {
		model.TfsZ = types.Float64Value(value)
		delete(additional, "tfs_z")
	}
	if value, ok := toInt64Value(data["repeat_last_n"]); ok {
		model.RepeatLastN = types.Int64Value(value)
		delete(additional, "repeat_last_n")
	}
	if value, ok := toFloat64Value(data["mirostat_tau"]); ok {
		model.MirostatTau = types.Float64Value(value)
		delete(additional, "mirostat_tau")
	}
	if value, ok := toFloat64Value(data["mirostat_eta"]); ok {
		model.MirostatEta = types.Float64Value(value)
		delete(additional, "mirostat_eta")
	}
	if value, ok := toInt64Value(data["mirostat"]); ok {
		model.Mirostat = types.Int64Value(value)
		delete(additional, "mirostat")
	}
	if value, ok := toFloat64Value(data["presence_penalty"]); ok {
		model.PresencePenalty = types.Float64Value(value)
		delete(additional, "presence_penalty")
	}
	if value, ok := toFloat64Value(data["frequency_penalty"]); ok {
		model.FrequencyPenalty = types.Float64Value(value)
		delete(additional, "frequency_penalty")
	}
	if value, ok := toFloat64Value(data["min_p"]); ok {
		model.MinP = types.Float64Value(value)
		delete(additional, "min_p")
	}
	if value, ok := toFloat64Value(data["top_p"]); ok {
		model.TopP = types.Float64Value(value)
		delete(additional, "top_p")
	}
	if value, ok := toInt64Value(data["top_k"]); ok {
		model.TopK = types.Int64Value(value)
		delete(additional, "top_k")
	}
	if value, ok := toInt64Value(data["max_tokens"]); ok {
		model.MaxTokens = types.Int64Value(value)
		delete(additional, "max_tokens")
	}
	if value, ok := toStringValue(data["reasoning_effort"]); ok {
		model.ReasoningEffort = types.StringValue(value)
		delete(additional, "reasoning_effort")
	}
	if raw, ok := data["custom_params"]; ok && raw != nil {
		stringMap, convOK := toStringMap(raw)
		if convOK {
			mapVal, mapDiags := types.MapValueFrom(ctx, types.StringType, stringMap)
			diags.Append(mapDiags...)
			if !mapDiags.HasError() {
				model.CustomParams = mapVal
			}
			delete(additional, "custom_params")
		} else {
			diags.AddError("Unexpected params value", fmt.Sprintf("Expected params.custom_params to be a map of string values, received %T", raw))
		}
	}

	if len(additional) == 0 {
		return model, types.StringNull(), diags
	}

	encoded, err := encodeOptionalJSON(additional)
	if err != nil {
		diags.AddError("Serialize params", err.Error())
		return model, types.StringNull(), diags
	}

	return model, encoded, diags
}

func flattenModelMeta(ctx context.Context, data map[string]any) (modelMetaState, types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	state := modelMetaState{
		ProfileImageURL:   types.StringNull(),
		Description:       types.StringNull(),
		SuggestionPrompts: types.ListNull(types.StringType),
		Tags:              types.ListNull(types.StringType),
		ToolIDs:           types.ListNull(types.StringType),
		DefaultFeatureIDs: types.ListNull(types.StringType),
		Capabilities:      nil,
	}

	if data == nil {
		return state, types.StringNull(), diags
	}

	additional := copyStringAnyMap(data)

	if value, ok := toStringValue(data["profile_image_url"]); ok {
		state.ProfileImageURL = types.StringValue(value)
		delete(additional, "profile_image_url")
	}
	if value, ok := toStringValue(data["description"]); ok {
		state.Description = types.StringValue(value)
		delete(additional, "description")
	}
	if raw, ok := data["capabilities"]; ok {
		switch v := raw.(type) {
		case map[string]any:
			state.Capabilities = flattenModelCapabilities(v)
			delete(additional, "capabilities")
		case nil:
			delete(additional, "capabilities")
		default:
			diags.AddError("Unexpected meta value", fmt.Sprintf("Expected meta.capabilities to be an object, received %T", raw))
		}
	}
	if raw, ok := data["suggestion_prompts"]; ok && raw != nil {
		prompts, convOK := toKeyedStringSlice(raw, "content")
		if convOK {
			list, listDiags := types.ListValueFrom(ctx, types.StringType, prompts)
			diags.Append(listDiags...)
			if !listDiags.HasError() {
				state.SuggestionPrompts = list
			}
			delete(additional, "suggestion_prompts")
		} else {
			diags.AddError("Unexpected meta value", fmt.Sprintf("Expected meta.suggestion_prompts to be a list of objects with content, received %T", raw))
		}
	}
	if raw, ok := data["tags"]; ok && raw != nil {
		tags, convOK := toKeyedStringSlice(raw, "name")
		if convOK {
			list, listDiags := types.ListValueFrom(ctx, types.StringType, tags)
			diags.Append(listDiags...)
			if !listDiags.HasError() {
				state.Tags = list
			}
			delete(additional, "tags")
		} else {
			diags.AddError("Unexpected meta value", fmt.Sprintf("Expected meta.tags to be a list of objects with name, received %T", raw))
		}
	}
	if raw, ok := data["toolIds"]; ok && raw != nil {
		tools, convOK := toStringSlice(raw)
		if convOK {
			list, listDiags := types.ListValueFrom(ctx, types.StringType, tools)
			diags.Append(listDiags...)
			if !listDiags.HasError() {
				state.ToolIDs = list
			}
			delete(additional, "toolIds")
		} else {
			diags.AddError("Unexpected meta value", fmt.Sprintf("Expected meta.toolIds to be a list of strings, received %T", raw))
		}
	}
	if raw, ok := data["defaultFeatureIds"]; ok && raw != nil {
		features, convOK := toStringSlice(raw)
		if convOK {
			list, listDiags := types.ListValueFrom(ctx, types.StringType, features)
			diags.Append(listDiags...)
			if !listDiags.HasError() {
				state.DefaultFeatureIDs = list
			}
			delete(additional, "defaultFeatureIds")
		} else {
			diags.AddError("Unexpected meta value", fmt.Sprintf("Expected meta.defaultFeatureIds to be a list of strings, received %T", raw))
		}
	}

	if len(additional) == 0 {
		return state, types.StringNull(), diags
	}

	encoded, err := encodeOptionalJSON(additional)
	if err != nil {
		diags.AddError("Serialize metadata", err.Error())
		return state, types.StringNull(), diags
	}

	return state, encoded, diags
}

func flattenModelCapabilities(data map[string]any) *modelCapabilitiesModel {
	caps := &modelCapabilitiesModel{
		Vision:          types.BoolNull(),
		FileUpload:      types.BoolNull(),
		WebSearch:       types.BoolNull(),
		ImageGeneration: types.BoolNull(),
		CodeInterpreter: types.BoolNull(),
		Citations:       types.BoolNull(),
		StatusUpdates:   types.BoolNull(),
		Usage:           types.BoolNull(),
	}

	if data == nil {
		return caps
	}

	if value, ok := toBoolValue(data["vision"]); ok {
		caps.Vision = types.BoolValue(value)
	}
	if value, ok := toBoolValue(data["file_upload"]); ok {
		caps.FileUpload = types.BoolValue(value)
	}
	if value, ok := toBoolValue(data["web_search"]); ok {
		caps.WebSearch = types.BoolValue(value)
	}
	if value, ok := toBoolValue(data["image_generation"]); ok {
		caps.ImageGeneration = types.BoolValue(value)
	}
	if value, ok := toBoolValue(data["code_interpreter"]); ok {
		caps.CodeInterpreter = types.BoolValue(value)
	}
	if value, ok := toBoolValue(data["citations"]); ok {
		caps.Citations = types.BoolValue(value)
	}
	if value, ok := toBoolValue(data["status_updates"]); ok {
		caps.StatusUpdates = types.BoolValue(value)
	}
	if value, ok := toBoolValue(data["usage"]); ok {
		caps.Usage = types.BoolValue(value)
	}

	return caps
}

func mergeStringAnyMaps(primary, secondary map[string]any) map[string]any {
	if primary == nil && secondary == nil {
		return nil
	}

	if primary == nil {
		result := make(map[string]any, len(secondary))
		for k, v := range secondary {
			result[k] = v
		}
		return result
	}

	result := make(map[string]any, len(primary)+len(secondary))
	for k, v := range primary {
		result[k] = v
	}
	for k, v := range secondary {
		result[k] = v
	}

	return result
}

func copyStringAnyMap(input map[string]any) map[string]any {
	if input == nil {
		return nil
	}

	result := make(map[string]any, len(input))
	for k, v := range input {
		result[k] = v
	}

	return result
}

func toStringValue(value any) (string, bool) {
	if value == nil {
		return "", false
	}

	switch v := value.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}

func toBoolValue(value any) (bool, bool) {
	if value == nil {
		return false, false
	}

	switch v := value.(type) {
	case bool:
		return v, true
	default:
		return false, false
	}
}

func toFloat64Value(value any) (float64, bool) {
	if value == nil {
		return 0, false
	}

	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

func toInt64Value(value any) (int64, bool) {
	f, ok := toFloat64Value(value)
	if !ok {
		return 0, false
	}

	return int64(f), true
}

func toStringSlice(value any) ([]string, bool) {
	switch v := value.(type) {
	case []string:
		return v, true
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			str, ok := toStringValue(item)
			if !ok {
				return nil, false
			}
			result = append(result, str)
		}
		return result, true
	default:
		return nil, false
	}
}

func toKeyedStringSlice(value any, key string) ([]string, bool) {
	switch v := value.(type) {
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if item == nil {
				continue
			}
			if str, ok := toStringValue(item); ok {
				result = append(result, str)
				continue
			}
			itemMap, ok := item.(map[string]any)
			if !ok {
				return nil, false
			}
			str, ok := toStringValue(itemMap[key])
			if !ok {
				return nil, false
			}
			result = append(result, str)
		}
		return result, true
	case []map[string]any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			str, ok := toStringValue(item[key])
			if !ok {
				return nil, false
			}
			result = append(result, str)
		}
		return result, true
	case []string:
		return v, true
	default:
		return nil, false
	}
}

func toStringMap(value any) (map[string]string, bool) {
	switch v := value.(type) {
	case map[string]string:
		return v, true
	case map[string]any:
		result := make(map[string]string, len(v))
		for key, raw := range v {
			if str, ok := toStringValue(raw); ok {
				result[key] = str
				continue
			}
			if raw == nil {
				result[key] = ""
				continue
			}
			result[key] = fmt.Sprint(raw)
		}
		return result, true
	default:
		return nil, false
	}
}
