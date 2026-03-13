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

var _ resource.Resource = &gitConnectionResource{}

type gitConnectionResource struct{ client *Client }

type gitConnectionResourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Provider types.String `tfsdk:"provider_name"`
	BaseURL  types.String `tfsdk:"base_url"`
	OrgName  types.String `tfsdk:"org_name"`
	Token    types.String `tfsdk:"token"`
}

func NewGitConnectionResource() resource.Resource { return &gitConnectionResource{} }

func (r *gitConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_git_connection"
}

func (r *gitConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Git connection for browsing repositories and branches in course configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for this Git connection.",
			},
			"provider_name": schema.StringAttribute{
				Required:    true,
				Description: "Git provider type: github or gitea.",
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "API base URL (for Gitea). Not needed for GitHub.",
			},
			"org_name": schema.StringAttribute{
				Required:    true,
				Description: "Organization or username to list repos from.",
			},
			"token": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "API access token (write-only).",
			},
		},
	}
}

func (r *gitConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *gitConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan gitConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"name":     plan.Name.ValueString(),
		"provider": plan.Provider.ValueString(),
		"org_name": plan.OrgName.ValueString(),
		"token":    plan.Token.ValueString(),
	}
	if !plan.BaseURL.IsNull() {
		body["base_url"] = plan.BaseURL.ValueString()
	}

	var result APIGitConnection
	if err := r.client.Post("/api/git-connections", body, &result); err != nil {
		resp.Diagnostics.AddError("Create Git connection failed", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(result.ID))
	plan.BaseURL = types.StringValue(result.BaseURL)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *gitConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state gitConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result APIGitConnection
	if err := r.client.Get(fmt.Sprintf("/api/git-connections/%d", state.ID.ValueInt64()), &result); err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(result.Name)
	state.Provider = types.StringValue(result.Provider)
	state.BaseURL = types.StringValue(result.BaseURL)
	state.OrgName = types.StringValue(result.OrgName)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *gitConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan gitConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{}
	if !plan.Name.IsNull() {
		body["name"] = plan.Name.ValueString()
	}
	if !plan.OrgName.IsNull() {
		body["org_name"] = plan.OrgName.ValueString()
	}
	if !plan.Token.IsNull() {
		body["token"] = plan.Token.ValueString()
	}
	if !plan.BaseURL.IsNull() {
		body["base_url"] = plan.BaseURL.ValueString()
	}

	var result APIGitConnection
	if err := r.client.Put(fmt.Sprintf("/api/git-connections/%d", plan.ID.ValueInt64()), body, &result); err != nil {
		resp.Diagnostics.AddError("Update Git connection failed", err.Error())
		return
	}

	plan.BaseURL = types.StringValue(result.BaseURL)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *gitConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state gitConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(fmt.Sprintf("/api/git-connections/%d", state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Delete Git connection failed", err.Error())
	}
}
