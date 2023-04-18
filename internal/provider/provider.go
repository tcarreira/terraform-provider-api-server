package provider

import (
	"context"

	"github.com/tcarreira/api-server/pkg/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure APIServerProvider satisfies various provider interfaces.
var _ provider.Provider = &APIServerProvider{}

// APIServerProvider defines the provider implementation.
type APIServerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version   string
	APIClient *client.APIClient
}

// APIServerProviderModel describes the provider data model.
type APIServerProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *APIServerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "apiserver"
	resp.Version = p.version
}

func (p *APIServerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "API Server endpoint",
				Required:            true,
			},
		},
	}
}

func (p *APIServerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data APIServerProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if p.APIClient == nil && !data.Endpoint.IsNull() {
		cli, err := client.NewAPIClient(client.Config{
			Endpoint: data.Endpoint.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("failed to create api client", err.Error())
			return
		}
		p.APIClient = cli
	}

	resp.DataSourceData = p.APIClient
	resp.ResourceData = p.APIClient
}

func (p *APIServerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *APIServerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPersonDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &APIServerProvider{
			version: version,
		}
	}
}
