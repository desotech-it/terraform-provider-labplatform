package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &connectionTemplateResource{}

type connectionTemplateResource struct{ client *Client }

type connectionTemplateResourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	CourseID          types.Int64  `tfsdk:"course_id"`
	Name              types.String `tfsdk:"name"`
	Protocol          types.String `tfsdk:"protocol"`
	Hostname          types.String `tfsdk:"hostname"`
	Port              types.Int64  `tfsdk:"port"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	Parameters        types.String `tfsdk:"parameters"`
	VsphereEndpointID types.Int64  `tfsdk:"vsphere_endpoint_id"`
	GuestID           types.String `tfsdk:"guest_id"`
}

func NewConnectionTemplateResource() resource.Resource { return &connectionTemplateResource{} }

func (r *connectionTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection_template"
}

func (r *connectionTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a connection template (VNC, RDP, SSH, or vSphere). Templates define how students connect to lab machines.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"course_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Optional course ID to associate this template with a specific course.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Template name (e.g. 'Linux VNC Desktop 1').",
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Connection protocol: vnc, rdp, ssh, or vsphere.",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Target hostname or IP. For vsphere, this is the VM path.",
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Target port (e.g. 5901 for VNC, 3389 for RDP, 2222 for SSH).",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Connection username.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Connection password (write-only, not returned by API).",
			},
			"parameters": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Additional connection parameters as JSON string.",
			},
			"vsphere_endpoint_id": schema.Int64Attribute{
				Optional:    true,
				Description: "vSphere endpoint ID (required for vsphere protocol).",
			},
			"guest_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "vSphere guest OS ID for the VM.",
			},
		},
	}
}

func (r *connectionTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *connectionTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan connectionTemplateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"name":     plan.Name.ValueString(),
		"protocol": plan.Protocol.ValueString(),
	}
	setOptional(body, "course_id", plan.CourseID)
	setOptionalStr(body, "hostname", plan.Hostname)
	setOptional(body, "port", plan.Port)
	setOptionalStr(body, "username", plan.Username)
	setOptionalStr(body, "password", plan.Password)
	setOptionalStr(body, "parameters", plan.Parameters)
	setOptional(body, "vsphere_endpoint_id", plan.VsphereEndpointID)
	setOptionalStr(body, "guest_id", plan.GuestID)

	var result APIConnectionTemplate
	if err := r.client.Post("/api/templates", body, &result); err != nil {
		resp.Diagnostics.AddError("Create template failed", err.Error())
		return
	}

	r.apiToState(&plan, &result)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *connectionTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state connectionTemplateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API doesn't have a single-template GET, use list and filter
	var templates []APIConnectionTemplate
	if err := r.client.Get("/api/templates", &templates); err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	var found *APIConnectionTemplate
	for _, t := range templates {
		if int64(t.ID) == state.ID.ValueInt64() {
			found = &t
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	r.apiToState(&state, found)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *connectionTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan connectionTemplateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"name":     plan.Name.ValueString(),
		"protocol": plan.Protocol.ValueString(),
	}
	setOptional(body, "course_id", plan.CourseID)
	setOptionalStr(body, "hostname", plan.Hostname)
	setOptional(body, "port", plan.Port)
	setOptionalStr(body, "username", plan.Username)
	setOptionalStr(body, "password", plan.Password)
	setOptionalStr(body, "parameters", plan.Parameters)
	setOptional(body, "vsphere_endpoint_id", plan.VsphereEndpointID)
	setOptionalStr(body, "guest_id", plan.GuestID)

	var result APIConnectionTemplate
	if err := r.client.Put(fmt.Sprintf("/api/templates/%d", plan.ID.ValueInt64()), body, &result); err != nil {
		resp.Diagnostics.AddError("Update template failed", err.Error())
		return
	}

	r.apiToState(&plan, &result)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *connectionTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state connectionTemplateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(fmt.Sprintf("/api/templates/%d", state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Delete template failed", err.Error())
	}
}

func (r *connectionTemplateResource) apiToState(state *connectionTemplateResourceModel, api *APIConnectionTemplate) {
	state.ID = types.Int64Value(int64(api.ID))
	state.Name = types.StringValue(api.Name)
	state.Protocol = types.StringValue(api.Protocol)
	state.Hostname = types.StringValue(api.Hostname)
	state.Port = types.Int64Value(int64(api.Port))
	state.Username = types.StringValue(api.Username)
	// Password is write-only — API never returns it, preserve from state/plan
	if api.Password != "" {
		state.Password = types.StringValue(api.Password)
	}
	state.Parameters = types.StringValue(api.Parameters)
	state.GuestID = types.StringValue(api.GuestID)
	if api.CourseID != nil {
		state.CourseID = types.Int64Value(int64(*api.CourseID))
	}
	if api.VsphereEndpointID != nil {
		state.VsphereEndpointID = types.Int64Value(int64(*api.VsphereEndpointID))
	}
}

// Helpers for optional fields
func setOptional(m map[string]interface{}, key string, v types.Int64) {
	if !v.IsNull() && !v.IsUnknown() {
		m[key] = v.ValueInt64()
	}
}

func setOptionalStr(m map[string]interface{}, key string, v types.String) {
	if !v.IsNull() && !v.IsUnknown() {
		m[key] = v.ValueString()
	}
}
