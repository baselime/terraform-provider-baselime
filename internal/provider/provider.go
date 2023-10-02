package provider

import (
	"context"
	"github.com/baselime/terraform-provider-baselime/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure BaselimeProvider satisfies various provider interfaces.
var _ provider.Provider = &BaselimeProvider{}

// BaselimeProvider defines the provider implementation.
type BaselimeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type BaselimeResourceData struct {
	Client *client.Client
}

type DataSourceData struct {
	Client *client.Client
}

// BaselimeProviderModel describes the provider data model.
type BaselimeProviderModel struct {
	ApiHost   types.String `tfsdk:"api_host"`
	ApiKey    types.String `tfsdk:"api_key" sensitive:"true"`
	ApiScheme types.String `tfsdk:"api_scheme"`
}

func (p *BaselimeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "baselime"
	resp.Version = p.version
}

func (p *BaselimeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_host": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Required: true,
			},
			"api_scheme": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *BaselimeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data BaselimeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	c := client.NewClient(&client.Config{
		APIKey:    data.ApiKey.ValueString(),
		APIHost:   data.ApiHost.ValueString(),
		ApiScheme: data.ApiScheme.ValueString(),
		Debug:     false,
	})
	resp.DataSourceData = &DataSourceData{
		Client: c,
	}
	resp.ResourceData = &BaselimeResourceData{
		Client: c,
	}
}

func (p *BaselimeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewQueryResource,
		NewAlertResource,
		NewDashboardResource,
	}
}

func (p *BaselimeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BaselimeProvider{
			version: version,
		}
	}
}
