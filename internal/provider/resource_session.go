package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &sessionResource{}

type sessionResource struct{ client *Client }

type sessionResourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	CourseID   types.Int64  `tfsdk:"course_id"`
	TrainerIDs types.List   `tfsdk:"trainer_ids"`
	Status     types.String `tfsdk:"status"`
	Notes      types.String `tfsdk:"notes"`
	Days       types.List   `tfsdk:"days"`
}

type sessionDayModel struct {
	Date      types.String `tfsdk:"date"`
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
}

var sessionDayAttrTypes = map[string]attr.Type{
	"date":       types.StringType,
	"start_time": types.StringType,
	"end_time":   types.StringType,
}

func NewSessionResource() resource.Resource { return &sessionResource{} }

func (r *sessionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session"
}

func (r *sessionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a course session (class). A session represents a scheduled instance of a course with specific dates, trainers, and students.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"course_id": schema.Int64Attribute{
				Required:    true,
				Description: "ID of the course this session belongs to.",
			},
			"trainer_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "List of trainer user IDs assigned to this session.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("scheduled"),
				Description: "Session status: scheduled, active, completed, or cancelled.",
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Session notes.",
			},
			"days": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of session days with date and time ranges.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"date": schema.StringAttribute{
							Required:    true,
							Description: "Day date in YYYY-MM-DD format.",
						},
						"start_time": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("09:00"),
							Description: "Start time in HH:MM format (default: 09:00).",
						},
						"end_time": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("18:00"),
							Description: "End time in HH:MM format (default: 18:00).",
						},
					},
				},
			},
		},
	}
}

func (r *sessionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *sessionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sessionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"course_id": plan.CourseID.ValueInt64(),
	}

	// Trainer IDs
	if !plan.TrainerIDs.IsNull() {
		var tids []int64
		plan.TrainerIDs.ElementsAs(ctx, &tids, false)
		body["trainer_ids"] = tids
	}

	if !plan.Notes.IsNull() {
		body["notes"] = plan.Notes.ValueString()
	}

	// Days
	var days []sessionDayModel
	plan.Days.ElementsAs(ctx, &days, false)
	if len(days) > 0 {
		body["start_date"] = days[0].Date.ValueString()
		body["end_date"] = days[len(days)-1].Date.ValueString()
		apiDays := make([]map[string]string, len(days))
		for i, d := range days {
			apiDays[i] = map[string]string{
				"date":       d.Date.ValueString(),
				"start_time": d.StartTime.ValueString(),
				"end_time":   d.EndTime.ValueString(),
			}
		}
		body["days"] = apiDays
	}

	var result APISession
	if err := r.client.Post("/api/sessions", body, &result); err != nil {
		resp.Diagnostics.AddError("Create session failed", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(result.ID))
	plan.Status = types.StringValue(result.Status)
	plan.Notes = types.StringValue(result.Notes)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *sessionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sessionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result APISession
	if err := r.client.Get(fmt.Sprintf("/api/sessions/%d", state.ID.ValueInt64()), &result); err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.CourseID = types.Int64Value(int64(result.CourseID))
	state.Status = types.StringValue(result.Status)
	state.Notes = types.StringValue(result.Notes)

	// Trainer IDs
	if len(result.Trainers) > 0 {
		tids := make([]attr.Value, len(result.Trainers))
		for i, t := range result.Trainers {
			tids[i] = types.Int64Value(int64(t.ID))
		}
		state.TrainerIDs, _ = types.ListValue(types.Int64Type, tids)
	}

	// Days
	if len(result.Days) > 0 {
		dayValues := make([]attr.Value, len(result.Days))
		for i, d := range result.Days {
			startTime := d.StartTime
			if len(startTime) > 5 {
				startTime = startTime[:5]
			}
			endTime := d.EndTime
			if len(endTime) > 5 {
				endTime = endTime[:5]
			}
			// Handle ISO format
			if len(d.StartTime) > 11 {
				startTime = d.StartTime[11:16]
			}
			if len(d.EndTime) > 11 {
				endTime = d.EndTime[11:16]
			}
			dayDate := d.DayDate
			if len(dayDate) > 10 {
				dayDate = dayDate[:10]
			}
			obj, _ := types.ObjectValue(sessionDayAttrTypes, map[string]attr.Value{
				"date":       types.StringValue(dayDate),
				"start_time": types.StringValue(startTime),
				"end_time":   types.StringValue(endTime),
			})
			dayValues[i] = obj
		}
		state.Days, _ = types.ListValue(types.ObjectType{AttrTypes: sessionDayAttrTypes}, dayValues)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *sessionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sessionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{}

	if !plan.TrainerIDs.IsNull() {
		var tids []int64
		plan.TrainerIDs.ElementsAs(ctx, &tids, false)
		body["trainer_ids"] = tids
	}
	if !plan.Status.IsNull() {
		body["status"] = plan.Status.ValueString()
	}
	if !plan.Notes.IsNull() {
		body["notes"] = plan.Notes.ValueString()
	}

	var days []sessionDayModel
	plan.Days.ElementsAs(ctx, &days, false)
	if len(days) > 0 {
		apiDays := make([]map[string]string, len(days))
		for i, d := range days {
			apiDays[i] = map[string]string{
				"date":       d.Date.ValueString(),
				"start_time": d.StartTime.ValueString(),
				"end_time":   d.EndTime.ValueString(),
			}
		}
		body["days"] = apiDays
	}

	var result APISession
	if err := r.client.Put(fmt.Sprintf("/api/sessions/%d", plan.ID.ValueInt64()), body, &result); err != nil {
		resp.Diagnostics.AddError("Update session failed", err.Error())
		return
	}

	plan.Status = types.StringValue(result.Status)
	plan.Notes = types.StringValue(result.Notes)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *sessionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sessionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(fmt.Sprintf("/api/sessions/%d", state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Delete session failed", err.Error())
	}
}
