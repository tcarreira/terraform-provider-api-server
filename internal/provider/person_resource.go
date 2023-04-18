package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tcarreira/api-server/pkg/client"
	apiTypes "github.com/tcarreira/api-server/pkg/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PersonResource{}
var _ resource.ResourceWithImportState = &PersonResource{}

func NewPersonResource() resource.Resource {
	return &PersonResource{}
}

// PersonResource defines the resource implementation.
type PersonResource struct {
	client *client.APIClient
}

// PersonResourceModel describes the resource data model.
type PersonResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Age         types.Int64  `tfsdk:"age"`
	Description types.String `tfsdk:"description"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *PersonResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_person"
}

func (r *PersonResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Person resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Person identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Person name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"age": schema.Int64Attribute{
				MarkdownDescription: "Person age",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Person description",
				Optional:            true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *PersonResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	cli, ok := req.ProviderData.(*client.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = cli
}

func (r *PersonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *PersonResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	person := &apiTypes.Person{
		Name:        data.Name.ValueString(),
		Age:         int(data.Age.ValueInt64()),
		Description: data.Description.ValueString(),
	}
	err := r.client.People().Create(person)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create person, got error: %s", err))
		return
	}

	data.Id = types.StringValue(strconv.Itoa(person.ID))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	tflog.Info(ctx, "created a resource", map[string]interface{}{
		"person": person,
		"data":   data,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PersonResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("(read) Error converting id to int", err.Error())
		return
	}
	person, err := r.client.People().Get(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting person", err.Error())
		return
	}

	data.Name = types.StringValue(person.Name)
	data.Age = types.Int64Value(int64(person.Age))
	data.Description = types.StringValue(person.Description)
	if data.LastUpdated.ValueString() == "" {
		data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	}

	tflog.Info(ctx, "read a resource", map[string]interface{}{
		"person": person,
		"data":   data,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *PersonResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("(update) Error converting id to int", err.Error())
		return
	}
	person := &apiTypes.Person{
		Name:        data.Name.ValueString(),
		Age:         int(data.Age.ValueInt64()),
		Description: data.Description.ValueString(),
	}
	err = r.client.People().Update(id, person)
	if err != nil {
		resp.Diagnostics.AddError("Error updating person", err.Error())
		return
	}

	data.Name = types.StringValue(person.Name)
	data.Age = types.Int64Value(int64(person.Age))
	data.Description = types.StringValue(person.Description)
	if data.LastUpdated.ValueString() == "" {
		data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	}

	tflog.Info(ctx, "updated a resource", map[string]interface{}{
		"person": person,
		"data":   data,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PersonResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error converting id to int", err.Error())
		return
	}
	err = r.client.People().Delete(id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting person", err.Error())
		return
	}
}

func (r *PersonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
