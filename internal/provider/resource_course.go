package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &courseResource{}

type courseResource struct{ client *Client }

type courseResourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	GuideRepo       types.String `tfsdk:"guide_repo"`
	DurationDays    types.Int64  `tfsdk:"duration_days"`
	GuideBranch     types.String `tfsdk:"guide_branch"`
	GitConnectionID types.Int64  `tfsdk:"git_connection_id"`
}

func NewCourseResource() resource.Resource { return &courseResource{} }

func (r *courseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_course"
}

func (r *courseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a LabPlatform course.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Course name (e.g. 'CKA - Certified Kubernetes Admin').",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Course description.",
			},
			"guide_repo": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "GitHub guide repository (e.g. 'desotech-it/CKA').",
			},
			"duration_days": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Course duration in days (default: 5).",
			},
			"guide_branch": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("main"),
				Description: "Git branch for the guide (default: main).",
			},
			"git_connection_id": schema.Int64Attribute{
				Optional:    true,
				Description: "ID of the Git connection to use for repo/branch browsing.",
			},
		},
	}
}

func (r *courseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *courseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan courseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"name": plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body["description"] = plan.Description.ValueString()
	}
	if !plan.GuideRepo.IsNull() {
		body["guide_repo"] = plan.GuideRepo.ValueString()
	}
	if !plan.DurationDays.IsNull() {
		body["duration_days"] = plan.DurationDays.ValueInt64()
	}
	if !plan.GuideBranch.IsNull() {
		body["guide_branch"] = plan.GuideBranch.ValueString()
	}
	if !plan.GitConnectionID.IsNull() {
		body["git_connection_id"] = plan.GitConnectionID.ValueInt64()
	}

	var result APICourse
	if err := r.client.Post("/api/courses", body, &result); err != nil {
		resp.Diagnostics.AddError("Create course failed", err.Error())
		return
	}

	r.apiToState(&plan, &result)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *courseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state courseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result APICourse
	if err := r.client.Get(fmt.Sprintf("/api/courses/%d", state.ID.ValueInt64()), &result); err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	r.apiToState(&state, &result)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *courseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan courseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{}
	if !plan.Name.IsNull() {
		body["name"] = plan.Name.ValueString()
	}
	if !plan.Description.IsNull() {
		body["description"] = plan.Description.ValueString()
	}
	if !plan.GuideRepo.IsNull() {
		body["guide_repo"] = plan.GuideRepo.ValueString()
	}
	if !plan.DurationDays.IsNull() {
		body["duration_days"] = plan.DurationDays.ValueInt64()
	}
	if !plan.GuideBranch.IsNull() {
		body["guide_branch"] = plan.GuideBranch.ValueString()
	}
	if !plan.GitConnectionID.IsNull() {
		body["git_connection_id"] = plan.GitConnectionID.ValueInt64()
	}

	var result APICourse
	if err := r.client.Put(fmt.Sprintf("/api/courses/%d", plan.ID.ValueInt64()), body, &result); err != nil {
		resp.Diagnostics.AddError("Update course failed", err.Error())
		return
	}

	r.apiToState(&plan, &result)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *courseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state courseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(fmt.Sprintf("/api/courses/%d", state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Delete course failed", err.Error())
	}
}

func (r *courseResource) apiToState(state *courseResourceModel, api *APICourse) {
	state.ID = types.Int64Value(int64(api.ID))
	state.Name = types.StringValue(api.Name)
	state.Description = types.StringValue(api.Description)
	state.GuideRepo = types.StringValue(api.GuideRepo)
	state.DurationDays = types.Int64Value(int64(api.DurationDays))
	state.GuideBranch = types.StringValue(api.GuideBranch)
	if api.GitConnectionID != nil {
		state.GitConnectionID = types.Int64Value(int64(*api.GitConnectionID))
	}
}
