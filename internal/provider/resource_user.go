package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &userResource{}

type userResource struct{ client *Client }

type userResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Username  types.String `tfsdk:"username"`
	Email     types.String `tfsdk:"email"`
	Password  types.String `tfsdk:"password"`
	Role      types.String `tfsdk:"role"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Company   types.String `tfsdk:"company"`
	Phone     types.String `tfsdk:"phone"`
	Language  types.String `tfsdk:"language"`
}

func NewUserResource() resource.Resource { return &userResource{} }

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a LabPlatform user (student, trainer, or admin).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "Unique username for authentication.",
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Email address.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Password for the user (write-only, not returned by API).",
			},
			"role": schema.StringAttribute{
				Required:    true,
				Description: "User role: student, trainer, or admin.",
			},
			"first_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "First name.",
			},
			"last_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Last name.",
			},
			"company": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Company name.",
			},
			"phone": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Phone number.",
			},
			"language": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("it"),
				Description: "Preferred language (default: it).",
			},
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"username": plan.Username.ValueString(),
		"password": plan.Password.ValueString(),
		"role":     plan.Role.ValueString(),
	}
	if !plan.Email.IsNull() {
		body["email"] = plan.Email.ValueString()
	}
	if !plan.FirstName.IsNull() {
		body["first_name"] = plan.FirstName.ValueString()
	}
	if !plan.LastName.IsNull() {
		body["last_name"] = plan.LastName.ValueString()
	}
	if !plan.Company.IsNull() {
		body["company"] = plan.Company.ValueString()
	}
	if !plan.Phone.IsNull() {
		body["phone"] = plan.Phone.ValueString()
	}
	if !plan.Language.IsNull() {
		body["language"] = plan.Language.ValueString()
	}

	var result APIUser
	if err := r.client.Post("/api/users", body, &result); err != nil {
		resp.Diagnostics.AddError("Create user failed", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(result.ID))
	plan.Email = types.StringValue(result.Email)
	plan.FirstName = types.StringValue(result.FirstName)
	plan.LastName = types.StringValue(result.LastName)
	plan.Company = types.StringValue(result.Company)
	plan.Phone = types.StringValue(result.Phone)
	plan.Language = types.StringValue(result.Language)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result APIUser
	if err := r.client.Get(fmt.Sprintf("/api/users/%d", state.ID.ValueInt64()), &result); err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Username = types.StringValue(result.Username)
	state.Email = types.StringValue(result.Email)
	state.Role = types.StringValue(result.Role)
	state.FirstName = types.StringValue(result.FirstName)
	state.LastName = types.StringValue(result.LastName)
	state.Company = types.StringValue(result.Company)
	state.Phone = types.StringValue(result.Phone)
	state.Language = types.StringValue(result.Language)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{}
	if !plan.Email.IsNull() {
		body["email"] = plan.Email.ValueString()
	}
	if !plan.FirstName.IsNull() {
		body["first_name"] = plan.FirstName.ValueString()
	}
	if !plan.LastName.IsNull() {
		body["last_name"] = plan.LastName.ValueString()
	}
	if !plan.Company.IsNull() {
		body["company"] = plan.Company.ValueString()
	}
	if !plan.Phone.IsNull() {
		body["phone"] = plan.Phone.ValueString()
	}
	if !plan.Language.IsNull() {
		body["language"] = plan.Language.ValueString()
	}
	if !plan.Role.IsNull() {
		body["role"] = plan.Role.ValueString()
	}

	var result APIUser
	if err := r.client.Put(fmt.Sprintf("/api/users/%d", plan.ID.ValueInt64()), body, &result); err != nil {
		resp.Diagnostics.AddError("Update user failed", err.Error())
		return
	}

	plan.Email = types.StringValue(result.Email)
	plan.FirstName = types.StringValue(result.FirstName)
	plan.LastName = types.StringValue(result.LastName)
	plan.Company = types.StringValue(result.Company)
	plan.Phone = types.StringValue(result.Phone)
	plan.Language = types.StringValue(result.Language)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(fmt.Sprintf("/api/users/%d", state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Delete user failed", err.Error())
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "Expected numeric ID")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(id))...)
}
