package provider

import (
	"context"
	"fmt"
	"github.com/baselime/terraform-provider-baselime/client"
	"github.com/baselime/terraform-provider-baselime/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &QueryResource{}

var _ resource.ResourceWithImportState = &QueryResource{}

func NewQueryResource() resource.Resource {
	return &QueryResource{}
}

// QueryResource defines the resource implementation.
type QueryResource struct {
	client *client.Client
}

func (r *QueryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_query"
}

func (r *QueryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Query resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Query name",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Query description",
			},
			"service": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Query service",
			},
			"datasets": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "Query datasets",
				ElementType:         types.StringType,
			},
			"filters": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "Query filters",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"key":       types.StringType,
						"operation": types.StringType,
						"value":     types.StringType,
						"type":      types.StringType,
					},
				},
			},
			"filter_combination": schema.StringAttribute{
				Optional:            true,
				Default:             stringdefault.StaticString("OR"),
				MarkdownDescription: "Query filter combination",
				Computed:            true,
			},
			"calculations": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Query calculations",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"key":      types.StringType,
						"operator": types.StringType,
						"alias":    types.StringType,
					},
				},
			},
			"group_by": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Query group by",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":  types.StringType,
						"value": types.StringType,
					},
				},
			},
			"order_by": schema.ObjectAttribute{
				Optional: true,
				AttributeTypes: map[string]attr.Type{
					"value": types.StringType,
					"order": types.StringType,
				},
			},
			"limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Query limit",
				Default:             int64default.StaticInt64(50),
				Computed:            true,
			},
			"needle": schema.ObjectAttribute{
				Optional: true,
				AttributeTypes: map[string]attr.Type{
					"value":      types.StringType,
					"is_regex":   types.BoolType,
					"match_case": types.BoolType,
				},
			},
		},
	}
}

func (r *QueryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*BaselimeResourceData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configure Type",
			fmt.Sprintf("Expected *BaselimeProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = provider.Client
}

func (r *QueryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.QueryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	data.Id = data.Name

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.CreateQuery(ctx, data.ToApiObject())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create query, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "query created", map[string]interface{}{
		"name": data.Name,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.QueryResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	query, err := r.client.GetQuery(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read query, got error: %s", err))
		return
	}
	if query != nil {
		data.FromApiObject(query)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.QueryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateQuery(ctx, data.ToApiObject())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update query, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.QueryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteQuery(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete query, got error: %s", err))
		return
	}
}

func (r *QueryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
