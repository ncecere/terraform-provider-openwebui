package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client"
	"github.com/ncecere/terraform-provider-openwebui/internal/provider/client/groups"
)

var (
	_ resource.Resource                = &GroupResource{}
	_ resource.ResourceWithImportState = &GroupResource{}
)

type GroupResource struct {
	client *client.OpenWebUIClient
}

type GroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	UserIDs     types.List   `tfsdk:"user_ids"`
	Permissions types.Object `tfsdk:"permissions"`
}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

func (r *GroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a group in OpenWebUI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the group.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the group.",
				Optional:    true,
			},
			"user_ids": schema.ListAttribute{
				Description: "List of user IDs in the group.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"permissions": schema.SingleNestedAttribute{
				Description: "Permissions for the group.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"workspace": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"models":    schema.BoolAttribute{Required: true},
							"knowledge": schema.BoolAttribute{Required: true},
							"prompts":   schema.BoolAttribute{Required: true},
							"tools":     schema.BoolAttribute{Required: true},
						},
					},
					"chat": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"file_upload": schema.BoolAttribute{Required: true},
							"delete":      schema.BoolAttribute{Required: true},
							"edit":        schema.BoolAttribute{Required: true},
							"temporary":   schema.BoolAttribute{Required: true},
						},
					},
				},
			},
		},
	}
}

func (r *GroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// First, create the group with basic information
	createGroup := &groups.Group{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	createdGroup, err := r.client.Groups.Create(createGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
			fmt.Sprintf("Could not create group: %s", err),
		)
		return
	}

	// Now prepare the update with all the additional information
	updateGroup := &groups.Group{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Handle user IDs
	var userIDs []string
	diags = plan.UserIDs.ElementsAs(ctx, &userIDs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateGroup.UserIDs = userIDs

	// Handle permissions
	if !plan.Permissions.IsNull() {
		var permissions struct {
			Workspace struct {
				Models    bool `tfsdk:"models"`
				Knowledge bool `tfsdk:"knowledge"`
				Prompts   bool `tfsdk:"prompts"`
				Tools     bool `tfsdk:"tools"`
			} `tfsdk:"workspace"`
			Chat struct {
				FileUpload bool `tfsdk:"file_upload"`
				Delete     bool `tfsdk:"delete"`
				Edit       bool `tfsdk:"edit"`
				Temporary  bool `tfsdk:"temporary"`
			} `tfsdk:"chat"`
		}
		diags = plan.Permissions.As(ctx, &permissions, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		updateGroup.Permissions = &groups.GroupPermissions{
			Workspace: groups.WorkspacePermissions{
				Models:    permissions.Workspace.Models,
				Knowledge: permissions.Workspace.Knowledge,
				Prompts:   permissions.Workspace.Prompts,
				Tools:     permissions.Workspace.Tools,
			},
			Chat: groups.ChatPermissions{
				FileUpload: permissions.Chat.FileUpload,
				Delete:     permissions.Chat.Delete,
				Edit:       permissions.Chat.Edit,
				Temporary:  permissions.Chat.Temporary,
			},
		}
	}

	// Update the group with all the information
	updatedGroup, err := r.client.Groups.Update(createdGroup.ID, updateGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating group",
			fmt.Sprintf("Could not update group with ID %s: %s", createdGroup.ID, err),
		)
		return
	}

	plan.ID = types.StringValue(updatedGroup.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.Groups.Get(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading group",
			fmt.Sprintf("Could not read group ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	state.Name = types.StringValue(group.Name)
	state.Description = types.StringValue(group.Description)

	userIDs, diags := types.ListValueFrom(ctx, types.StringType, group.UserIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.UserIDs = userIDs

	if group.Permissions != nil {
		workspaceAttrs := map[string]attr.Value{
			"models":    types.BoolValue(group.Permissions.Workspace.Models),
			"knowledge": types.BoolValue(group.Permissions.Workspace.Knowledge),
			"prompts":   types.BoolValue(group.Permissions.Workspace.Prompts),
			"tools":     types.BoolValue(group.Permissions.Workspace.Tools),
		}

		chatAttrs := map[string]attr.Value{
			"file_upload": types.BoolValue(group.Permissions.Chat.FileUpload),
			"delete":      types.BoolValue(group.Permissions.Chat.Delete),
			"edit":        types.BoolValue(group.Permissions.Chat.Edit),
			"temporary":   types.BoolValue(group.Permissions.Chat.Temporary),
		}

		workspaceObj, diags := types.ObjectValue(
			map[string]attr.Type{
				"models":    types.BoolType,
				"knowledge": types.BoolType,
				"prompts":   types.BoolType,
				"tools":     types.BoolType,
			},
			workspaceAttrs,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		chatObj, diags := types.ObjectValue(
			map[string]attr.Type{
				"file_upload": types.BoolType,
				"delete":      types.BoolType,
				"edit":        types.BoolType,
				"temporary":   types.BoolType,
			},
			chatAttrs,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		permissionsObj, diags := types.ObjectValue(
			map[string]attr.Type{
				"workspace": types.ObjectType{AttrTypes: workspaceObj.AttributeTypes(ctx)},
				"chat":      types.ObjectType{AttrTypes: chatObj.AttributeTypes(ctx)},
			},
			map[string]attr.Value{
				"workspace": workspaceObj,
				"chat":      chatObj,
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		state.Permissions = permissionsObj
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := &groups.Group{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	var userIDs []string
	diags = plan.UserIDs.ElementsAs(ctx, &userIDs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	group.UserIDs = userIDs

	if !plan.Permissions.IsNull() {
		var permissions struct {
			Workspace struct {
				Models    bool `tfsdk:"models"`
				Knowledge bool `tfsdk:"knowledge"`
				Prompts   bool `tfsdk:"prompts"`
				Tools     bool `tfsdk:"tools"`
			} `tfsdk:"workspace"`
			Chat struct {
				FileUpload bool `tfsdk:"file_upload"`
				Delete     bool `tfsdk:"delete"`
				Edit       bool `tfsdk:"edit"`
				Temporary  bool `tfsdk:"temporary"`
			} `tfsdk:"chat"`
		}
		diags = plan.Permissions.As(ctx, &permissions, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		group.Permissions = &groups.GroupPermissions{
			Workspace: groups.WorkspacePermissions{
				Models:    permissions.Workspace.Models,
				Knowledge: permissions.Workspace.Knowledge,
				Prompts:   permissions.Workspace.Prompts,
				Tools:     permissions.Workspace.Tools,
			},
			Chat: groups.ChatPermissions{
				FileUpload: permissions.Chat.FileUpload,
				Delete:     permissions.Chat.Delete,
				Edit:       permissions.Chat.Edit,
				Temporary:  permissions.Chat.Temporary,
			},
		}
	}

	updatedGroup, err := r.client.Groups.Update(plan.ID.ValueString(), group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating group",
			fmt.Sprintf("Could not update group ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	plan.ID = types.StringValue(updatedGroup.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Groups.Delete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting group",
			fmt.Sprintf("Could not delete group ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
