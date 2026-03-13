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

var _ resource.Resource = &vsphereEndpointResource{}

type vsphereEndpointResource struct{ client *Client }

type vsphereEndpointResourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	URL        types.String `tfsdk:"url"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Datacenter types.String `tfsdk:"datacenter"`
	Insecure   types.Bool   `tfsdk:"insecure"`
}

func NewVsphereEndpointResource() resource.Resource { return &vsphereEndpointResource{} }

func (r *vsphereEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vsphere_endpoint"
}

func (r *vsphereEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a vSphere vCenter endpoint for VM console access in labs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for this vCenter endpoint.",
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "vCenter URL (e.g. https://vcenter.example.com).",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "vCenter username.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "vCenter password (write-only).",
			},
			"datacenter": schema.StringAttribute{
				Required:    true,
				Description: "vCenter datacenter name.",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Skip TLS certificate verification (default: false).",
			},
		},
	}
}

func (r *vsphereEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *vsphereEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vsphereEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"name":       plan.Name.ValueString(),
		"url":        plan.URL.ValueString(),
		"username":   plan.Username.ValueString(),
		"password":   plan.Password.ValueString(),
		"datacenter": plan.Datacenter.ValueString(),
	}
	if !plan.Insecure.IsNull() {
		body["insecure"] = plan.Insecure.ValueBool()
	}

	var result APIVsphereEndpoint
	if err := r.client.Post("/api/vsphere/endpoints", body, &result); err != nil {
		resp.Diagnostics.AddError("Create vSphere endpoint failed", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(result.ID))
	plan.Insecure = types.BoolValue(result.Insecure)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *vsphereEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vsphereEndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result APIVsphereEndpoint
	if err := r.client.Get(fmt.Sprintf("/api/vsphere/endpoints/%d", state.ID.ValueInt64()), &result); err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(result.Name)
	state.URL = types.StringValue(result.URL)
	state.Username = types.StringValue(result.Username)
	state.Datacenter = types.StringValue(result.Datacenter)
	state.Insecure = types.BoolValue(result.Insecure)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *vsphereEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vsphereEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{}
	setOptionalStr(body, "name", plan.Name)
	setOptionalStr(body, "url", plan.URL)
	setOptionalStr(body, "username", plan.Username)
	setOptionalStr(body, "password", plan.Password)
	setOptionalStr(body, "datacenter", plan.Datacenter)
	if !plan.Insecure.IsNull() {
		body["insecure"] = plan.Insecure.ValueBool()
	}

	var result APIVsphereEndpoint
	if err := r.client.Put(fmt.Sprintf("/api/vsphere/endpoints/%d", plan.ID.ValueInt64()), body, &result); err != nil {
		resp.Diagnostics.AddError("Update vSphere endpoint failed", err.Error())
		return
	}

	plan.Insecure = types.BoolValue(result.Insecure)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *vsphereEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vsphereEndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(fmt.Sprintf("/api/vsphere/endpoints/%d", state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Delete vSphere endpoint failed", err.Error())
	}
}
