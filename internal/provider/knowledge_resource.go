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
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &KnowledgeResource{}
var _ resource.ResourceWithImportState = &KnowledgeResource{}

func NewKnowledgeResource() resource.Resource {
	return &KnowledgeResource{}
}

// KnowledgeResource defines the resource implementation.
type KnowledgeResource struct {
	client *Client
}

// KnowledgeResourceModel describes the resource data model.
type KnowledgeResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Data          types.Map    `tfsdk:"data"`
	AccessControl types.String `tfsdk:"access_control"`
	AccessGroups  types.List   `tfsdk:"access_groups"`
	AccessUsers   types.List   `tfsdk:"access_users"`
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
				MarkdownDescription: "Additional data for the knowledge base",
				Optional:            true,
			},
			"access_control": schema.StringAttribute{
				MarkdownDescription: "Access control type ('public' or 'private')",
				Required:            true,
			},
			"access_groups": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of group IDs with access (only used when access_control is 'private')",
				Optional:            true,
			},
			"access_users": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of user IDs with access (only used when access_control is 'private')",
				Optional:            true,
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

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func buildAccessControl(ctx context.Context, accessControl types.String, accessGroups types.List, accessUsers types.List) (map[string]interface{}, error) {
	if accessControl.ValueString() == "public" {
		return nil, nil
	}

	result := map[string]interface{}{
		"read": map[string]interface{}{
			"group_ids": []string{},
			"user_ids":  []string{},
		},
		"write": map[string]interface{}{
			"group_ids": []string{},
			"user_ids":  []string{},
		},
	}

	readAccess := result["read"].(map[string]interface{})
	writeAccess := result["write"].(map[string]interface{})

	// Convert group IDs
	var groupIds []string
	if !accessGroups.IsNull() {
		diags := accessGroups.ElementsAs(ctx, &groupIds, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting group IDs: %v", diags)
		}
	}
	readAccess["group_ids"] = groupIds
	writeAccess["group_ids"] = groupIds

	// Convert user IDs
	var userIds []string
	if !accessUsers.IsNull() {
		diags := accessUsers.ElementsAs(ctx, &userIds, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting user IDs: %v", diags)
		}
	}
	readAccess["user_ids"] = userIds
	writeAccess["user_ids"] = userIds

	return result, nil
}

func (r *KnowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KnowledgeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert from Terraform types to Go types
	knowledgeForm := &KnowledgeForm{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Convert data map if present
	if !data.Data.IsNull() {
		dataMap := make(map[string]string)
		resp.Diagnostics.Append(data.Data.ElementsAs(ctx, &dataMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		knowledgeForm.Data = dataMap
	}

	// Build access control structure
	accessControl, err := buildAccessControl(ctx, data.AccessControl, data.AccessGroups, data.AccessUsers)
	if err != nil {
		resp.Diagnostics.AddError("Access Control Error", err.Error())
		return
	}
	knowledgeForm.AccessControl = accessControl

	// Create new knowledge base
	result, err := r.client.CreateKnowledge(knowledgeForm)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create knowledge base, got error: %s", err))
		return
	}

	// Map response back into the model
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
	result, err := r.client.GetKnowledge(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read knowledge base, got error: %s", err))
		return
	}

	// Update model with the current state
	data.Name = types.StringValue(result.Name)
	data.Description = types.StringValue(result.Description)
	data.LastUpdated = types.StringValue(fmt.Sprint(result.UpdatedAt))

	// Update data map if present
	if result.Data != nil {
		dataMap, err := types.MapValueFrom(ctx, types.StringType, result.Data)
		if err != nil {
			resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Unable to convert data map, got error: %s", err))
			return
		}
		data.Data = dataMap
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

	// Convert from Terraform types to Go types
	knowledgeForm := &KnowledgeForm{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Convert data map if present
	if !data.Data.IsNull() {
		dataMap := make(map[string]string)
		resp.Diagnostics.Append(data.Data.ElementsAs(ctx, &dataMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		knowledgeForm.Data = dataMap
	}

	// Build access control structure
	accessControl, err := buildAccessControl(ctx, data.AccessControl, data.AccessGroups, data.AccessUsers)
	if err != nil {
		resp.Diagnostics.AddError("Access Control Error", err.Error())
		return
	}
	knowledgeForm.AccessControl = accessControl

	// Update knowledge base
	result, err := r.client.UpdateKnowledge(data.ID.ValueString(), knowledgeForm)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update knowledge base, got error: %s", err))
		return
	}

	// Update model with the response
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
	err := r.client.DeleteKnowledge(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete knowledge base, got error: %s", err))
		return
	}
}

func (r *KnowledgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
