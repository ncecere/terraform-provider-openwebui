package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/models"
)

var (
	_ resource.Resource                = &ModelResource{}
	_ resource.ResourceWithImportState = &ModelResource{}
)

func NewModelResource() resource.Resource {
	return &ModelResource{}
}

type ModelResource struct {
	client *models.Client
}

func (r *ModelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

func (r *ModelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clients, ok := req.ProviderData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected map[string]interface{}, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	client, ok := clients["models"].(*models.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *models.Client, got: %T. Please report this issue to the provider developers.", clients["models"]),
		)
		return
	}

	r.client = client
}

func (r *ModelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a model in OpenWebUI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the model.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The ID of the user who created the model.",
				Computed:    true,
			},
			"base_model_id": schema.StringAttribute{
				Description: "The ID of the base model.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the model.",
				Required:    true,
			},
			"params": schema.SingleNestedAttribute{
				Description: "Model parameters.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"system": schema.StringAttribute{
						Description: "System prompt for the model.",
						Optional:    true,
					},
					"stream_response": schema.BoolAttribute{
						Description: "Whether to stream responses.",
						Optional:    true,
					},
					"temperature": schema.Float64Attribute{
						Description: "Sampling temperature.",
						Optional:    true,
					},
					"top_p": schema.Float64Attribute{
						Description: "Top-p sampling parameter.",
						Optional:    true,
					},
					"top_k": schema.Int64Attribute{
						Description: "Top-k sampling parameter.",
						Optional:    true,
					},
					"min_p": schema.Float64Attribute{
						Description: "Minimum probability threshold.",
						Optional:    true,
					},
					"max_tokens": schema.Int64Attribute{
						Description: "Maximum number of tokens to generate.",
						Optional:    true,
					},
					"seed": schema.Int64Attribute{
						Description: "Random seed for reproducibility.",
						Optional:    true,
					},
					"frequency_penalty": schema.Int64Attribute{
						Description: "Frequency penalty.",
						Optional:    true,
					},
					"repeat_last_n": schema.Int64Attribute{
						Description: "Number of tokens to consider for repetition penalty.",
						Optional:    true,
					},
					"num_ctx": schema.Int64Attribute{
						Description: "Context window size.",
						Optional:    true,
					},
					"num_batch": schema.Int64Attribute{
						Description: "Batch size for processing.",
						Optional:    true,
					},
					"num_keep": schema.Int64Attribute{
						Description: "Number of tokens to keep from prompt.",
						Optional:    true,
					},
				},
			},
			"meta": schema.SingleNestedAttribute{
				Description: "Model metadata.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"profile_image_url": schema.StringAttribute{
						Description: "URL for the model's profile image.",
						Optional:    true,
					},
					"description": schema.StringAttribute{
						Description: "Description of the model.",
						Optional:    true,
					},
					"capabilities": schema.SingleNestedAttribute{
						Description: "Model capabilities.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"vision": schema.BoolAttribute{
								Description: "Whether the model supports vision tasks.",
								Optional:    true,
							},
							"usage": schema.BoolAttribute{
								Description: "Whether to track usage statistics.",
								Optional:    true,
							},
							"citations": schema.BoolAttribute{
								Description: "Whether the model supports citations.",
								Optional:    true,
							},
						},
					},
					"tags": schema.ListNestedAttribute{
						Description: "List of tags.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "Name of the tag.",
									Required:    true,
								},
							},
						},
					},
				},
			},
			"access_control": schema.SingleNestedAttribute{
				Description: "Access control settings.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"read": schema.SingleNestedAttribute{
						Description: "Read access settings.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"group_ids": schema.ListAttribute{
								Description: "List of group IDs with read access.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"user_ids": schema.ListAttribute{
								Description: "List of user IDs with read access.",
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
					"write": schema.SingleNestedAttribute{
						Description: "Write access settings.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"group_ids": schema.ListAttribute{
								Description: "List of group IDs with write access.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"user_ids": schema.ListAttribute{
								Description: "List of user IDs with write access.",
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
			"is_active": schema.BoolAttribute{
				Description: "Whether the model is active.",
				Optional:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Timestamp when the model was created.",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Timestamp when the model was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *ModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.CreateModel(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating model", err.Error())
		return
	}

	// Ensure the ID is set in the state
	if model.ID.IsNull() {
		resp.Diagnostics.AddError("Error creating model", "Model ID is null after creation")
		return
	}

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *ModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetModel(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading model", err.Error())
		return
	}

	// Ensure the ID is preserved
	if model.ID.IsNull() {
		model.ID = state.ID
	}

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *ModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.Model
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure we use the existing ID for the update
	plan.ID = state.ID

	model, err := r.client.UpdateModel(state.ID.ValueString(), &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error updating model", err.Error())
		return
	}

	// Ensure the ID is preserved
	if model.ID.IsNull() {
		model.ID = state.ID
	}

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *ModelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.Model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteModel(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting model", err.Error())
		return
	}
}

func (r *ModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
