
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

resource "baselime_alert" "terraformed" {
  name        = "terraformed-alert"
  description = "This alert was created by Terraform"
  enabled     = true
  channels = [
    {
      type    = "email"
      targets = ["foo@baselime.io"]
    }
  ]
  query = baselime_query.terraformed.id
  threshold = {
    operator = "GREATER_THAN"
    value    = 0
  }
  frequency = "5m"
  window    = "5m"
}