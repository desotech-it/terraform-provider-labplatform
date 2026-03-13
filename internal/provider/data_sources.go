package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// --- Users Data Source ---

var _ datasource.DataSource = &usersDataSource{}

type usersDataSource struct{ client *Client }
type usersDataSourceModel struct {
	Role  types.String `tfsdk:"role"`
	Users types.List   `tfsdk:"users"`
}

func NewUsersDataSource() datasource.DataSource { return &usersDataSource{} }

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists LabPlatform users, optionally filtered by role.",
		Attributes: map[string]schema.Attribute{
			"role": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by role: student, trainer, or admin.",
			},
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.Int64Attribute{Computed: true},
						"username":   schema.StringAttribute{Computed: true},
						"email":      schema.StringAttribute{Computed: true},
						"role":       schema.StringAttribute{Computed: true},
						"first_name": schema.StringAttribute{Computed: true},
						"last_name":  schema.StringAttribute{Computed: true},
						"company":    schema.StringAttribute{Computed: true},
						"language":   schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*Client)
}

var userObjAttrTypes = map[string]attr.Type{
	"id":         types.Int64Type,
	"username":   types.StringType,
	"email":      types.StringType,
	"role":       types.StringType,
	"first_name": types.StringType,
	"last_name":  types.StringType,
	"company":    types.StringType,
	"language":   types.StringType,
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config usersDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := "/api/users?per_page=500"
	if !config.Role.IsNull() {
		path += "&role=" + config.Role.ValueString()
	}

	var users []APIUser
	if err := d.client.Get(path, &users); err != nil {
		resp.Diagnostics.AddError("List users failed", err.Error())
		return
	}

	userValues := make([]attr.Value, len(users))
	for i, u := range users {
		obj, _ := types.ObjectValue(userObjAttrTypes, map[string]attr.Value{
			"id":         types.Int64Value(int64(u.ID)),
			"username":   types.StringValue(u.Username),
			"email":      types.StringValue(u.Email),
			"role":       types.StringValue(u.Role),
			"first_name": types.StringValue(u.FirstName),
			"last_name":  types.StringValue(u.LastName),
			"company":    types.StringValue(u.Company),
			"language":   types.StringValue(u.Language),
		})
		userValues[i] = obj
	}
	config.Users, _ = types.ListValue(types.ObjectType{AttrTypes: userObjAttrTypes}, userValues)

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

// --- Courses Data Source ---

var _ datasource.DataSource = &coursesDataSource{}

type coursesDataSource struct{ client *Client }
type coursesDataSourceModel struct {
	Courses types.List `tfsdk:"courses"`
}

func NewCoursesDataSource() datasource.DataSource { return &coursesDataSource{} }

func (d *coursesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_courses"
}

func (d *coursesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all LabPlatform courses.",
		Attributes: map[string]schema.Attribute{
			"courses": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.Int64Attribute{Computed: true},
						"name":          schema.StringAttribute{Computed: true},
						"description":   schema.StringAttribute{Computed: true},
						"guide_repo":    schema.StringAttribute{Computed: true},
						"duration_days": schema.Int64Attribute{Computed: true},
						"guide_branch":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *coursesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*Client)
}

var courseObjAttrTypes = map[string]attr.Type{
	"id":            types.Int64Type,
	"name":          types.StringType,
	"description":   types.StringType,
	"guide_repo":    types.StringType,
	"duration_days": types.Int64Type,
	"guide_branch":  types.StringType,
}

func (d *coursesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var courses []APICourse
	if err := d.client.Get("/api/courses?per_page=500", &courses); err != nil {
		resp.Diagnostics.AddError("List courses failed", err.Error())
		return
	}

	values := make([]attr.Value, len(courses))
	for i, c := range courses {
		obj, _ := types.ObjectValue(courseObjAttrTypes, map[string]attr.Value{
			"id":            types.Int64Value(int64(c.ID)),
			"name":          types.StringValue(c.Name),
			"description":   types.StringValue(c.Description),
			"guide_repo":    types.StringValue(c.GuideRepo),
			"duration_days": types.Int64Value(int64(c.DurationDays)),
			"guide_branch":  types.StringValue(c.GuideBranch),
		})
		values[i] = obj
	}

	var state coursesDataSourceModel
	state.Courses, _ = types.ListValue(types.ObjectType{AttrTypes: courseObjAttrTypes}, values)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// --- Templates Data Source ---

var _ datasource.DataSource = &templatesDataSource{}

type templatesDataSource struct{ client *Client }
type templatesDataSourceModel struct {
	Templates types.List `tfsdk:"templates"`
}

func NewTemplatesDataSource() datasource.DataSource { return &templatesDataSource{} }

func (d *templatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_templates"
}

func (d *templatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all connection templates.",
		Attributes: map[string]schema.Attribute{
			"templates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.Int64Attribute{Computed: true},
						"name":     schema.StringAttribute{Computed: true},
						"protocol": schema.StringAttribute{Computed: true},
						"hostname": schema.StringAttribute{Computed: true},
						"port":     schema.Int64Attribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *templatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*Client)
}

var templateObjAttrTypes = map[string]attr.Type{
	"id":       types.Int64Type,
	"name":     types.StringType,
	"protocol": types.StringType,
	"hostname": types.StringType,
	"port":     types.Int64Type,
}

func (d *templatesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var templates []APIConnectionTemplate
	if err := d.client.Get("/api/templates", &templates); err != nil {
		resp.Diagnostics.AddError("List templates failed", err.Error())
		return
	}

	values := make([]attr.Value, len(templates))
	for i, t := range templates {
		obj, _ := types.ObjectValue(templateObjAttrTypes, map[string]attr.Value{
			"id":       types.Int64Value(int64(t.ID)),
			"name":     types.StringValue(t.Name),
			"protocol": types.StringValue(t.Protocol),
			"hostname": types.StringValue(t.Hostname),
			"port":     types.Int64Value(int64(t.Port)),
		})
		values[i] = obj
	}

	var state templatesDataSourceModel
	state.Templates, _ = types.ListValue(types.ObjectType{AttrTypes: templateObjAttrTypes}, values)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
