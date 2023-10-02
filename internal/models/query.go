package models

import (
	"github.com/baselime/terraform-provider-baselime/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type QueryFilter struct {
	Key       types.String `json:"key" tfsdk:"key"`
	Operation types.String `json:"operation" tfsdk:"operation"`
	Value     types.String `json:"value" tfsdk:"value"`
	Type      types.String `json:"type" tfsdk:"type"`
}

type SearchNeedle struct {
	Value     types.String `json:"value" tfsdk:"value"`
	IsRegex   types.Bool   `json:"isRegex" tfsdk:"is_regex"`
	MatchCase types.Bool   `json:"matchCase" tfsdk:"match_case"`
}

type FilterCombination string

var (
	FilterCombinationAnd FilterCombination = "AND"
	FilterCombinationOr  FilterCombination = "OR"
)

type QueryOrderBy struct {
	Value types.String `json:"value" tfsdk:"value"`
	Order types.String `json:"order" tfsdk:"order"`
}

type QueryGroupBy struct {
	Type  types.String `json:"type" tfsdk:"type"`
	Value types.String `json:"value" tfsdk:"value"`
}

type QueryCalculation struct {
	Key      types.String `json:"key" tfsdk:"key"`
	Operator types.String `json:"operator" tfsdk:"operator"`
	Alias    types.String `json:"alias" tfsdk:"alias"`
}

func (c *QueryCalculation) ToApiModel() client.QueryCalculation {
	return client.QueryCalculation{
		Key:      c.Key.ValueString(),
		Operator: c.Operator.ValueString(),
		Alias:    c.Alias.ValueString(),
	}
}

func (qf *QueryFilter) ToApiModel() *client.QueryFilter {
	return &client.QueryFilter{
		Key:       qf.Key.ValueString(),
		Operation: qf.Operation.ValueString(),
		Value:     qf.Value.ValueString(),
		Type:      qf.Type.ValueString(),
	}
}

func (qgb *QueryGroupBy) ToApiModel() *client.QueryGroupBy {
	return &client.QueryGroupBy{
		Type:  qgb.Type.ValueString(),
		Value: qgb.Value.ValueString(),
	}
}

func (qob *QueryOrderBy) ToApiModel() *client.QueryOrderBy {
	return &client.QueryOrderBy{
		Value: qob.Value.ValueString(),
		Order: qob.Order.ValueString(),
	}
}

func (sn *SearchNeedle) ToApiModel() *client.SearchNeedle {
	return &client.SearchNeedle{
		Value:     sn.Value.ValueString(),
		IsRegex:   sn.IsRegex.ValueBool(),
		MatchCase: sn.MatchCase.ValueBool(),
	}
}

// QueryResourceModel describes the resource data model.
type QueryResourceModel struct {
	Id                types.String       `tfsdk:"id"`
	Name              types.String       `tfsdk:"name"`
	Description       types.String       `tfsdk:"description"`
	Service           types.String       `tfsdk:"service"`
	Datasets          []string           `tfsdk:"datasets"`
	Filters           []QueryFilter      `tfsdk:"filters"`
	FilterCombination types.String       `tfsdk:"filter_combination"`
	Calculations      []QueryCalculation `tfsdk:"calculations"`
	GroupBy           []QueryGroupBy     `tfsdk:"group_by"`
	OrderBy           QueryOrderBy       `tfsdk:"order_by"`
	Limit             types.Int64        `tfsdk:"limit"`
	Needle            SearchNeedle       `tfsdk:"needle"`
}

func (data *QueryResourceModel) FromApiObject(obj *client.Query) {
	data.Name = types.StringValue(obj.Name)
	data.Description = types.StringValue(obj.Description)
	data.Service = types.StringValue(obj.Service)
	data.Datasets = obj.Parameters.Datasets
	if obj.Parameters.Filters != nil {
		data.Filters = func() []QueryFilter {
			filters := make([]QueryFilter, 0)
			for _, f := range obj.Parameters.Filters {
				filters = append(filters, QueryFilter{
					Key:       types.StringValue(f.Key),
					Operation: types.StringValue(f.Operation),
					Value:     types.StringValue(f.Value),
					Type:      types.StringValue(f.Type),
				})
			}
			return filters
		}()
	}
	data.FilterCombination = types.StringValue(string(obj.Parameters.FilterCombination))
	data.Calculations = func() []QueryCalculation {
		cals := make([]QueryCalculation, 0)
		for _, c := range obj.Parameters.Calculations {
			cals = append(cals, QueryCalculation{
				Key:      types.StringValue(c.Key),
				Operator: types.StringValue(c.Operator),
				Alias:    types.StringValue(c.Alias),
			})
		}
		return cals
	}()
	if obj.Parameters.GroupBy != nil {
		data.GroupBy = func() []QueryGroupBy {
			groups := make([]QueryGroupBy, 0)
			for _, g := range obj.Parameters.GroupBy {
				groups = append(groups, QueryGroupBy{
					Type:  types.StringValue(g.Type),
					Value: types.StringValue(g.Value),
				})
			}
			return groups
		}()
	}
	if obj.Parameters.OrderBy != nil {
		data.OrderBy = QueryOrderBy{
			Value: types.StringValue(obj.Parameters.OrderBy.Value),
			Order: types.StringValue(obj.Parameters.OrderBy.Order),
		}
	}
	data.Limit = types.Int64Value(obj.Parameters.Limit)
	if obj.Parameters.Needle != nil {
		data.Needle = SearchNeedle{
			Value:     types.StringValue(obj.Parameters.Needle.Value),
			IsRegex:   types.BoolValue(obj.Parameters.Needle.IsRegex),
			MatchCase: types.BoolValue(obj.Parameters.Needle.MatchCase),
		}
	}
}

func (data *QueryResourceModel) ToApiObject() *client.Query {
	return &client.Query{
		Id:          data.Name.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Service:     data.Service.ValueString(),
		Parameters: client.QueryParameters{
			Datasets: data.Datasets,
			Filters: func() []client.QueryFilter {
				filters := make([]client.QueryFilter, 0, len(data.Filters))
				for _, f := range data.Filters {
					filters = append(filters, *f.ToApiModel())
				}
				return filters
			}(),
			FilterCombination: data.FilterCombination.ValueString(),
			Calculations: func() []client.QueryCalculation {
				cals := make([]client.QueryCalculation, 0)
				for _, c := range data.Calculations {
					cals = append(cals, c.ToApiModel())
				}
				return cals
			}(),
			GroupBy: func() []client.QueryGroupBy {
				groups := make([]client.QueryGroupBy, 0)
				for _, g := range data.GroupBy {
					groups = append(groups, *g.ToApiModel())
				}
				return groups
			}(),
			OrderBy: data.OrderBy.ToApiModel(),
			Limit:   data.Limit.ValueInt64(),
			Needle:  data.Needle.ToApiModel(),
		},
	}
}
