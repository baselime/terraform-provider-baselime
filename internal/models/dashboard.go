package models

import (
	"github.com/baselime/terraform-provider-baselime/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DashboardResourceModel struct {
	Name        types.String      `tfsdk:"name"`
	Description types.String      `tfsdk:"description"`
	Widgets     []DashboardWidget `tfsdk:"widgets"`
}

type DashboardWidget struct {
	QueryId     types.String `tfsdk:"query_id"`
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (d *DashboardResourceModel) ToApiModel() *client.Dashboard {
	return &client.Dashboard{
		Parameters: client.DashboardParameters{
			Widgets: func() []client.DashboardWidget {
				widgets := make([]client.DashboardWidget, 0, len(d.Widgets))
				for _, widget := range d.Widgets {
					widgets = append(widgets, client.DashboardWidget{
						QueryId:     widget.QueryId.ValueString(),
						Type:        client.WidgetType(widget.Type.ValueString()),
						Name:        widget.Name.ValueString(),
						Description: widget.Description.ValueString(),
					})
				}
				return widgets
			}(),
		},
		Id:          d.Name.ValueString(),
		Name:        d.Name.ValueString(),
		Description: d.Description.ValueString(),
	}
}

func (d *DashboardResourceModel) FromApiModel(dashboard *client.Dashboard) {
	if dashboard == nil {
		return
	}
	d.Name = types.StringValue(dashboard.Name)
	d.Description = types.StringValue(dashboard.Description)
	d.Widgets = func() []DashboardWidget {
		widgets := make([]DashboardWidget, 0, len(dashboard.Parameters.Widgets))
		for _, widget := range dashboard.Parameters.Widgets {
			widgets = append(widgets, DashboardWidget{
				QueryId:     types.StringValue(widget.QueryId),
				Type:        types.StringValue(string(widget.Type)),
				Name:        types.StringValue(widget.Name),
				Description: types.StringValue(widget.Description),
			})
		}
		return widgets
	}()
}
