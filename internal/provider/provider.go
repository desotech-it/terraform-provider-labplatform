package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &labplatformProvider{}

type labplatformProvider struct{}

type labplatformProviderModel struct {
	URL      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func New() provider.Provider {
	return &labplatformProvider{}
}

func (p *labplatformProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "labplatform"
}

func (p *labplatformProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for DesoLabs LabPlatform. Manages users, courses, classes, connection templates, and lab environments.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "LabPlatform base URL (e.g. https://labplatform.desolabs.it). Can also be set via LABPLATFORM_URL env var.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Admin username for LabPlatform API. Can also be set via LABPLATFORM_USERNAME env var.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Admin password for LabPlatform API. Can also be set via LABPLATFORM_PASSWORD env var.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *labplatformProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config labplatformProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("LABPLATFORM_URL")
	username := os.Getenv("LABPLATFORM_USERNAME")
	password := os.Getenv("LABPLATFORM_PASSWORD")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddError("Missing URL", "Set 'url' in provider config or LABPLATFORM_URL env var")
		return
	}
	if username == "" {
		resp.Diagnostics.AddError("Missing Username", "Set 'username' in provider config or LABPLATFORM_USERNAME env var")
		return
	}
	if password == "" {
		resp.Diagnostics.AddError("Missing Password", "Set 'password' in provider config or LABPLATFORM_PASSWORD env var")
		return
	}

	client, err := NewClient(url, username, password)
	if err != nil {
		resp.Diagnostics.AddError("Authentication Failed", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *labplatformProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewCourseResource,
		NewConnectionTemplateResource,
		NewClassResource,
		NewClassStudentResource,
		NewGitConnectionResource,
		NewVsphereEndpointResource,
	}
}

func (p *labplatformProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUsersDataSource,
		NewCoursesDataSource,
		NewTemplatesDataSource,
	}
}
