package models

import (
	"github.com/baselime/terraform-provider-baselime/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AlertResourceModel struct {
	Name        types.String    `tfsdk:"name"`
	Description types.String    `tfsdk:"description"`
	Enabled     types.Bool      `tfsdk:"enabled"`
	Channels    []AlertChannel  `tfsdk:"channels"`
	Query       types.String    `tfsdk:"query"`
	Threshold   *AlertThreshold `tfsdk:"threshold"`
	Frequency   types.String    `tfsdk:"frequency"`
	Window      types.String    `tfsdk:"window"`
}

type AlertChannel struct {
	Type    types.String   `tfsdk:"type"`
	Targets []types.String `tfsdk:"targets"`
}

type AlertThreshold struct {
	Operator types.String `tfsdk:"operator"`
	Value    types.Number `tfsdk:"value"`
}

func (a *AlertResourceModel) ToApiModel() *client.Alert {
	return &client.Alert{
		Parameters: client.AlertParameters{
			QueryId: a.Query.ValueString(),
			Threshold: client.AlertThreshold{
				Operation: a.Threshold.Operator.ValueString(),
				Value:     a.Threshold.Value.ValueBigFloat(),
			},
			Frequency: a.Frequency.ValueString(),
			Window:    a.Window.ValueString(),
		},
		Id:          a.Name.ValueString(),
		Description: a.Description.ValueString(),
		Enabled:     a.Enabled.ValueBool(),
		Channels: func() []client.AlertChannel {
			channels := make([]client.AlertChannel, len(a.Channels))
			for i, channel := range a.Channels {
				channels[i] = client.AlertChannel{
					Type: channel.Type.ValueString(),
					Targets: func() []string {
						targets := make([]string, len(channel.Targets))
						for i, target := range channel.Targets {
							targets[i] = target.ValueString()
						}
						return targets
					}(),
				}
			}
			return channels
		}(),
	}
}

func (a *AlertResourceModel) FromApiModel(alert *client.Alert) {
	a.Name = types.StringValue(alert.Id)
	a.Description = types.StringValue(alert.Description)
	a.Enabled = types.BoolValue(alert.Enabled)
	a.Channels = func() []AlertChannel {
		channels := make([]AlertChannel, len(alert.Channels))
		for i, channel := range alert.Channels {
			channels[i] = AlertChannel{
				Type: types.StringValue(channel.Type),
				Targets: func() []types.String {
					targets := make([]types.String, len(channel.Targets))
					for i, target := range channel.Targets {
						targets[i] = types.StringValue(target)
					}
					return targets
				}(),
			}
		}
		return channels
	}()
	a.Threshold.Operator = types.StringValue(alert.Parameters.Threshold.Operation)
	a.Threshold.Value = types.NumberValue(alert.Parameters.Threshold.Value)
	a.Frequency = types.StringValue(alert.Parameters.Frequency)
	a.Window = types.StringValue(alert.Parameters.Window)
	a.Query = types.StringValue(alert.Parameters.QueryId)
}
