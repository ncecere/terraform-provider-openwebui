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
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/knowledge"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &KnowledgeResource{}
var _ resource.ResourceWithImportState = &KnowledgeResource{}

func NewKnowledgeResource() resource.Resource {
	return &KnowledgeResource{}
}

// KnowledgeResource defines the resource implementation.
type KnowledgeResource struct {
	client *client.OpenWebUIClient
}

// KnowledgeResourceModel describes the resource data model.
type KnowledgeResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Data          types.Map    `tfsdk:"data"`
	AccessControl types.String `tfsdk:"access_control"`
	LastUpdated   types.String `tfsdk:"last_updated"`
}

func (r *KnowledgeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

func (r *KnowledgeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Knowledge resource for OpenWebUI",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Knowledge identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the knowledge base",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the knowledge base",
				Required:            true,
			},
			"data": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Additional data for the knowledge base",
			},
			"access_control": schema.StringAttribute{
				MarkdownDescription: "Access control type ('public' or 'private')",
				Optional:            true,
				Computed:            true,
			},
			"last_updated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of the last update",
			},
		},
	}
}

func (r *KnowledgeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.OpenWebUIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.OpenWebUIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *KnowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KnowledgeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert data model to API form
	form := &knowledge.KnowledgeForm{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Handle data map
	if !data.Data.IsNull() {
		dataMap := make(map[string]string)
		resp.Diagnostics.Append(data.Data.ElementsAs(ctx, &dataMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		form.Data = dataMap
	}

	// Handle access control
	if !data.AccessControl.IsNull() {
		accessType := data.AccessControl.ValueString()
		if accessType == "private" {
			form.AccessControl = map[string]interface{}{
				"type": "private",
			}
		}
	}

	// Create new knowledge base
	result, err := r.client.Create(form)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create knowledge base, got error: %s", err))
		return
	}

	// Map response to model
	data.ID = types.StringValue(result.ID)
	data.LastUpdated = types.StringValue(fmt.Sprint(result.UpdatedAt))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KnowledgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KnowledgeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get knowledge base from API
	result, err := r.client.Get(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read knowledge base, got error: %s", err))
		return
	}

	// Map response to model
	data.Name = types.StringValue(result.Name)
	data.Description = types.StringValue(result.Description)
	data.LastUpdated = types.StringValue(fmt.Sprint(result.UpdatedAt))

	// Convert data map
	if result.Data != nil {
		dataMap := make(map[string]string)
		for k, v := range result.Data {
			if str, ok := v.(string); ok {
				dataMap[k] = str
			}
		}
		convertedMap, diags := types.MapValueFrom(ctx, types.StringType, dataMap)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Data = convertedMap
	}

	// Handle access control
	if result.AccessControl == nil {
		data.AccessControl = types.StringValue("public")
	} else {
		data.AccessControl = types.StringValue("private")
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KnowledgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data KnowledgeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert data model to API form
	form := &knowledge.KnowledgeForm{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Handle data map
	if !data.Data.IsNull() {
		dataMap := make(map[string]string)
		resp.Diagnostics.Append(data.Data.ElementsAs(ctx, &dataMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		form.Data = dataMap
	}

	// Handle access control
	if !data.AccessControl.IsNull() {
		accessType := data.AccessControl.ValueString()
		if accessType == "private" {
			form.AccessControl = map[string]interface{}{
				"type": "private",
			}
		}
	}

	// Update knowledge base
	result, err := r.client.Update(data.ID.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update knowledge base, got error: %s", err))
		return
	}

	// Update last updated timestamp
	data.LastUpdated = types.StringValue(fmt.Sprint(result.UpdatedAt))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KnowledgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KnowledgeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete knowledge base
	err := r.client.Delete(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete knowledge base, got error: %s", err))
		return
	}
}

func (r *KnowledgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
