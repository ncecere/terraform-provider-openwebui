package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

var _ datasource.DataSource = &modelDataSource{}
var _ datasource.DataSourceWithConfigure = &modelDataSource{}

// modelDataSource exposes read-only model information.
type modelDataSource struct {
	client *client.Client
}

// modelDataSourceModel reuses the resource model structure for mapping purposes.
type modelDataSourceModel = modelResourceModel

// NewModelDataSource constructs a new model data source instance.
func NewModelDataSource() datasource.DataSource {
	return &modelDataSource{}
}

// Metadata sets the data source identifier.
func (d *modelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

// Schema describes the model data source schema.
func (d *modelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"model_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the model to look up (matches the model_id used when creating the resource).",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier returned by Open WebUI for the model.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable model name.",
			},
			"base_model_id": schema.StringAttribute{
				Computed:    true,
				Description: "Base model identifier if supplied.",
			},
			"is_active": schema.BoolAttribute{
				Computed:    true,
				Description: "Flag indicating whether the model is active.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the user who owns the model.",
			},
			"created_at": schema.Int64Attribute{
				Computed:    true,
				Description: "Unix timestamp when the model was created.",
			},
			"updated_at": schema.Int64Attribute{
				Computed:    true,
				Description: "Unix timestamp when the model was last updated.",
			},
			"meta_additional_json": schema.StringAttribute{
				Computed:    true,
				Description: "Raw JSON fragment preserving metadata fields not modelled directly.",
			},
			"params_additional_json": schema.StringAttribute{
				Computed:    true,
				Description: "Raw JSON fragment preserving parameter fields not modelled directly.",
			},
			"read_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Group names granted read access to the model.",
			},
			"write_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Group names granted write access to the model.",
			},
			"params": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Parameter values returned by Open WebUI.",
				Attributes: map[string]schema.Attribute{
					"system":                  schema.StringAttribute{Computed: true},
					"stream_response":         schema.BoolAttribute{Computed: true},
					"stream_delta_chunk_size": schema.Int64Attribute{Computed: true},
					"function_calling":        schema.StringAttribute{Computed: true},
					"reasoning_tags": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"seed":              schema.Int64Attribute{Computed: true},
					"temperature":       schema.Float64Attribute{Computed: true},
					"keep_alive":        schema.StringAttribute{Computed: true},
					"num_gpu":           schema.Int64Attribute{Computed: true},
					"num_thread":        schema.Int64Attribute{Computed: true},
					"num_batch":         schema.Int64Attribute{Computed: true},
					"num_ctx":           schema.Int64Attribute{Computed: true},
					"num_keep":          schema.Int64Attribute{Computed: true},
					"format":            schema.StringAttribute{Computed: true},
					"think":             schema.BoolAttribute{Computed: true},
					"use_mlock":         schema.BoolAttribute{Computed: true},
					"use_mmap":          schema.BoolAttribute{Computed: true},
					"repeat_penalty":    schema.Float64Attribute{Computed: true},
					"tfs_z":             schema.Float64Attribute{Computed: true},
					"repeat_last_n":     schema.Int64Attribute{Computed: true},
					"mirostat_tau":      schema.Float64Attribute{Computed: true},
					"mirostat_eta":      schema.Float64Attribute{Computed: true},
					"mirostat":          schema.Int64Attribute{Computed: true},
					"presence_penalty":  schema.Float64Attribute{Computed: true},
					"frequency_penalty": schema.Float64Attribute{Computed: true},
					"min_p":             schema.Float64Attribute{Computed: true},
					"top_p":             schema.Float64Attribute{Computed: true},
					"top_k":             schema.Int64Attribute{Computed: true},
					"max_tokens":        schema.Int64Attribute{Computed: true},
					"reasoning_effort":  schema.StringAttribute{Computed: true},
					"custom_params": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "Additional key/value parameters returned by Open WebUI.",
					},
				},
			},
			"capabilities": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Capability flags advertised by the model.",
				Attributes: map[string]schema.Attribute{
					"vision":           schema.BoolAttribute{Computed: true},
					"file_upload":      schema.BoolAttribute{Computed: true},
					"web_search":       schema.BoolAttribute{Computed: true},
					"image_generation": schema.BoolAttribute{Computed: true},
					"code_interpreter": schema.BoolAttribute{Computed: true},
					"citations":        schema.BoolAttribute{Computed: true},
					"status_updates":   schema.BoolAttribute{Computed: true},
					"usage":            schema.BoolAttribute{Computed: true},
				},
			},
		},
	}
}

// Configure assigns the shared API client.
func (d *modelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

// Read fetches the model details from Open WebUI.
func (d *modelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before using the model data source.")
		return
	}

	var config modelDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ModelID.IsUnknown() || config.ModelID.IsNull() || config.ModelID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("model_id"),
			"Missing model identifier",
			"The model_id argument must be supplied to query an existing model.",
		)
		return
	}

	current, err := d.client.GetModel(ctx, config.ModelID.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.Diagnostics.AddAttributeError(
				path.Root("model_id"),
				"Model not found",
				"No Open WebUI model was found with the supplied model_id.",
			)
			return
		}

		resp.Diagnostics.AddError("Read model failed", err.Error())
		return
	}

	state, diags := modelResponseToModel(ctx, d.client, current, config.ModelID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
