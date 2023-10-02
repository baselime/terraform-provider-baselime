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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AlertResource{}

var _ resource.ResourceWithImportState = &AlertResource{}

func NewAlertResource() resource.Resource {
	return &AlertResource{}
}

// AlertResource defines the resource implementation.
type AlertResource struct {
	client *client.Client
}

func (r *AlertResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert"
}

func (r *AlertResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Alert resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Alert name",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Alert description",
			},
			"service": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Alert service",
			},
			"enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Alert enabled",
			},
			"channels": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "Alert channels",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type": types.StringType,
						"targets": types.ListType{
							ElemType: types.StringType,
						},
					},
				},
			},
			"query": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Alert query",
			},
			"threshold": schema.ObjectAttribute{
				Required:            true,
				MarkdownDescription: "Alert threshold",
				AttributeTypes: map[string]attr.Type{
					"operator": types.StringType,
					"value":    types.NumberType,
				},
			},
			"frequency": schema.StringAttribute{
				Required: true,
			},
			"window": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *AlertResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.AlertResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.CreateAlert(ctx, data.ToApiModel())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create alert, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.AlertResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	alert, err := r.client.GetAlert(ctx, data.Service.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert, got error: %s", err))
		return
	}
	data.FromApiModel(alert)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.AlertResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAlert(ctx, data.ToApiModel())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update alert, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.AlertResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAlert(ctx, data.Service.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete alert, got error: %s", err))
		return
	}
}

func (r *AlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
