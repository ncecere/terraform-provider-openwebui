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

var _ resource.Resource = &terminalServersConfigResource{}
var _ resource.ResourceWithConfigure = &terminalServersConfigResource{}
var _ resource.ResourceWithImportState = &terminalServersConfigResource{}

type terminalServersConfigResource struct{ client *client.Client }

type terminalServerConnectionModel struct {
	ConnectionID types.String `tfsdk:"connection_id"`
	Name         types.String `tfsdk:"name"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	URL          types.String `tfsdk:"url"`
	Path         types.String `tfsdk:"path"`
	Key          types.String `tfsdk:"key"`
	AuthType     types.String `tfsdk:"auth_type"`
	ConfigJSON   types.String `tfsdk:"config_json"`
	ServerType   types.String `tfsdk:"server_type"`
	PolicyID     types.String `tfsdk:"policy_id"`
	PolicyJSON   types.String `tfsdk:"policy_json"`
}

type terminalServersConfigModel struct {
	ID          types.String                    `tfsdk:"id"`
	Connections []terminalServerConnectionModel `tfsdk:"connections"`
}

func NewTerminalServersConfigResource() resource.Resource { return &terminalServersConfigResource{} }

func (r *terminalServersConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_terminal_servers_config"
}

func (r *terminalServersConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{Computed: true, Description: "Singleton identifier for the terminal servers config.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"connections": schema.ListNestedAttribute{Required: true, Description: "Terminal server connection entries.", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"connection_id": schema.StringAttribute{Optional: true, Description: "Optional Open WebUI terminal server connection id."},
			"name":          schema.StringAttribute{Optional: true, Description: "Optional display name."},
			"enabled":       schema.BoolAttribute{Optional: true, Description: "Whether this terminal server is enabled."},
			"url":           schema.StringAttribute{Required: true, Description: "Terminal server base URL."},
			"path":          schema.StringAttribute{Optional: true, Description: "Terminal server OpenAPI/config path."},
			"key":           schema.StringAttribute{Optional: true, Sensitive: true, Description: "Optional authentication key."},
			"auth_type":     schema.StringAttribute{Optional: true, Description: "Optional authentication type."},
			"config_json":   schema.StringAttribute{Optional: true, Description: "Optional terminal server config JSON."},
			"server_type":   schema.StringAttribute{Optional: true, Computed: true, Description: "Detected or configured server type."},
			"policy_id":     schema.StringAttribute{Optional: true, Description: "Optional orchestrator policy id."},
			"policy_json":   schema.StringAttribute{Optional: true, Description: "Optional policy JSON."},
		}}},
	}}
}

func (r *terminalServersConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		r.client = client
	}
}

func (r *terminalServersConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan terminalServersConfigModel
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing terminal servers config.")
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	state, diags := applyTerminalServersConfig(ctx, r.client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *terminalServersConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing terminal servers config.")
		return
	}
	config, err := r.client.GetTerminalServersConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Read terminal servers config failed", err.Error())
		return
	}
	state := terminalServersConfigModel{ID: types.StringValue("terminal_servers"), Connections: flattenTerminalServerConnections(config.Connections)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *terminalServersConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan terminalServersConfigModel
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured API client", "Expected provider to configure the Open WebUI client before managing terminal servers config.")
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	state, diags := applyTerminalServersConfig(ctx, r.client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *terminalServersConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
func (r *terminalServersConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func applyTerminalServersConfig(ctx context.Context, apiClient *client.Client, plan terminalServersConfigModel) (terminalServersConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	connections := make([]client.TerminalServerConnection, 0, len(plan.Connections))
	for i, item := range plan.Connections {
		config := decodeOptionalJSON(item.ConfigJSON, path.Root("connections").AtListIndex(i).AtName("config_json"), &diags)
		policy := decodeOptionalJSON(item.PolicyJSON, path.Root("connections").AtListIndex(i).AtName("policy_json"), &diags)
		conn := client.TerminalServerConnection{URL: item.URL.ValueString(), Config: config, Policy: policy}
		if !item.ConnectionID.IsNull() && !item.ConnectionID.IsUnknown() {
			v := item.ConnectionID.ValueString()
			conn.ID = &v
		}
		if !item.Name.IsNull() && !item.Name.IsUnknown() {
			v := item.Name.ValueString()
			conn.Name = &v
		}
		if !item.Enabled.IsNull() && !item.Enabled.IsUnknown() {
			v := item.Enabled.ValueBool()
			conn.Enabled = &v
		}
		if !item.Path.IsNull() && !item.Path.IsUnknown() {
			v := item.Path.ValueString()
			conn.Path = &v
		}
		if !item.Key.IsNull() && !item.Key.IsUnknown() {
			v := item.Key.ValueString()
			conn.Key = &v
		}
		if !item.AuthType.IsNull() && !item.AuthType.IsUnknown() {
			v := item.AuthType.ValueString()
			conn.AuthType = &v
		}
		if !item.ServerType.IsNull() && !item.ServerType.IsUnknown() {
			v := item.ServerType.ValueString()
			conn.ServerType = &v
		}
		if !item.PolicyID.IsNull() && !item.PolicyID.IsUnknown() {
			v := item.PolicyID.ValueString()
			conn.PolicyID = &v
		}
		connections = append(connections, conn)
	}
	if diags.HasError() {
		return terminalServersConfigModel{}, diags
	}
	updated, err := apiClient.SetTerminalServersConfig(ctx, client.TerminalServersConfigForm{Connections: connections})
	if err != nil {
		diags.AddError("Update terminal servers config failed", err.Error())
		return terminalServersConfigModel{}, diags
	}
	return terminalServersConfigModel{ID: types.StringValue("terminal_servers"), Connections: flattenTerminalServerConnections(updated.Connections)}, diags
}

func flattenTerminalServerConnections(connections []client.TerminalServerConnection) []terminalServerConnectionModel {
	result := make([]terminalServerConnectionModel, 0, len(connections))
	for _, conn := range connections {
		configJSON, _ := encodeOptionalJSON(conn.Config)
		policyJSON, _ := encodeOptionalJSON(conn.Policy)
		item := terminalServerConnectionModel{ConnectionID: types.StringNull(), Name: types.StringNull(), Enabled: types.BoolNull(), URL: types.StringValue(conn.URL), Path: types.StringNull(), Key: types.StringNull(), AuthType: types.StringNull(), ConfigJSON: configJSON, ServerType: types.StringNull(), PolicyID: types.StringNull(), PolicyJSON: policyJSON}
		if conn.ID != nil {
			item.ConnectionID = types.StringValue(*conn.ID)
		}
		if conn.Name != nil {
			item.Name = types.StringValue(*conn.Name)
		}
		if conn.Enabled != nil {
			item.Enabled = types.BoolValue(*conn.Enabled)
		}
		if conn.Path != nil {
			item.Path = types.StringValue(*conn.Path)
		}
		if conn.Key != nil {
			item.Key = types.StringValue(*conn.Key)
		}
		if conn.AuthType != nil {
			item.AuthType = types.StringValue(*conn.AuthType)
		}
		if conn.ServerType != nil {
			item.ServerType = types.StringValue(*conn.ServerType)
		}
		if conn.PolicyID != nil {
			item.PolicyID = types.StringValue(*conn.PolicyID)
		}
		result = append(result, item)
	}
	return result
}
