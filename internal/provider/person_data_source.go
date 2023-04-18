package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tcarreira/api-server/pkg/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &PersonDataSource{}
	_ datasource.DataSourceWithConfigure = &PersonDataSource{}
)

func NewPersonDataSource() datasource.DataSource {
	return &PersonDataSource{}
}

// PersonDataSource defines the data source implementation.
type PersonDataSource struct {
	client *client.APIClient
}

// PersonDataSourceModel describes the data source data model.
type PersonDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Age         types.Int64  `tfsdk:"age"`
	Description types.String `tfsdk:"description"`
}

func (d *PersonDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_person"
}

func (d *PersonDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Person data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Person identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Person name",
				Computed:            true,
			},
			"age": schema.Int64Attribute{
				MarkdownDescription: "Person age",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Person description",
				Computed:            true,
			},
		},
	}
}

func (d *PersonDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *PersonDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PersonDataSourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error converting id to int", err.Error())
		return
	}
	person, err := d.client.People().Get(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting person", err.Error())
		return
	}

	data.Name = types.StringValue(person.Name)
	data.Age = types.Int64Value(int64(person.Age))
	data.Description = types.StringValue(person.Description)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
