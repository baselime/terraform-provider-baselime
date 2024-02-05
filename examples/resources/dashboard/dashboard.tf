
resource "baselime_query" "terraformed" {
  name        = "terraformed-query"
  description = "This query was created by Terraform"
  datasets    = ["lambda-logs"]
  filters = [
    {
      key       = "message"
      operation = "INCLUDES"
      value     = "error"
      type      = "string"
    }
  ]
  filter_combination = "AND"
  calculations = [
    {
      key      = ""
      operator = "COUNT"
      alias    = "count"
    }
  ]
  group_by = [
    {
      type  = "string"
      value = "message"
    }
  ]
  order_by = {
    value = "count"
    order = "DESC"
  }
  limit = 10
  needle = {
    value      = ".*"
    is_regex   = true
    match_case = false
  }
}

resource "baselime_dashboard" "terraformed" {
  name        = "terraformed-dashboard"
  description = "This alert was created by Terraform"
  widgets = [
    {
      query_id    = baselime_query.terraformed.id
      type        = "timeseries"
      name        = "Line Chart"
      description = "This is a line chart"
    }
  ]
}