package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &classStudentResource{}

type classStudentResource struct{ client *Client }

type classStudentResourceModel struct {
	ClassID     types.Int64 `tfsdk:"class_id"`
	UserID      types.Int64 `tfsdk:"user_id"`
	TemplateIDs types.List  `tfsdk:"template_ids"`
}

func NewClassStudentResource() resource.Resource { return &classStudentResource{} }

func (r *classStudentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_class_student"
}

func (r *classStudentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Assigns a student to a class with specific connection templates. Creates the lab and remote connections.",
		Attributes: map[string]schema.Attribute{
			"class_id": schema.Int64Attribute{
				Required:    true,
				Description: "ID of the class to assign the student to.",
			},
			"user_id": schema.Int64Attribute{
				Required:    true,
				Description: "ID of the student user.",
			},
			"template_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "List of connection template IDs to provision for this student.",
			},
		},
	}
}

func (r *classStudentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *classStudentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan classStudentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"user_id": plan.UserID.ValueInt64(),
	}
	if !plan.TemplateIDs.IsNull() {
		var tids []int64
		plan.TemplateIDs.ElementsAs(ctx, &tids, false)
		body["template_ids"] = tids
	}

	if err := r.client.Post(
		fmt.Sprintf("/api/sessions/%d/students", plan.ClassID.ValueInt64()),
		body, nil,
	); err != nil {
		resp.Diagnostics.AddError("Assign student failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *classStudentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state classStudentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Verify the student is still assigned by listing class students
	var students []APILab
	if err := r.client.Get(
		fmt.Sprintf("/api/sessions/%d/students", state.ClassID.ValueInt64()),
		&students,
	); err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	found := false
	for _, s := range students {
		if int64(s.UserID) == state.UserID.ValueInt64() {
			found = true
			break
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *classStudentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Re-create: remove and re-add
	var plan classStudentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove old
	_ = r.client.Delete(fmt.Sprintf("/api/sessions/%d/students/%d",
		plan.ClassID.ValueInt64(), plan.UserID.ValueInt64()))

	// Re-add
	body := map[string]interface{}{
		"user_id": plan.UserID.ValueInt64(),
	}
	if !plan.TemplateIDs.IsNull() {
		var tids []int64
		plan.TemplateIDs.ElementsAs(ctx, &tids, false)
		body["template_ids"] = tids
	}

	if err := r.client.Post(
		fmt.Sprintf("/api/sessions/%d/students", plan.ClassID.ValueInt64()),
		body, nil,
	); err != nil {
		resp.Diagnostics.AddError("Update student assignment failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *classStudentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state classStudentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(
		fmt.Sprintf("/api/sessions/%d/students/%d",
			state.ClassID.ValueInt64(), state.UserID.ValueInt64()),
	); err != nil {
		resp.Diagnostics.AddError("Remove student failed", err.Error())
	}
}
