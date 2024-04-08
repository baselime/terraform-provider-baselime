terraform {
  required_providers {
    baselime = {
      version = "~> 0.1.6"
      source  = "baselime/baselime"
    }
  }
}

provider "baselime" {
  api_key = "your_api_key"
}

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
    operator = ">"
    value    = 0
  }
  frequency = "10m"
  window    = "5m"
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