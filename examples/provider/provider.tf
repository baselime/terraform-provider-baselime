terraform {
  required_providers {
    baselime = {
      version = "~> 1.0.0"
      source  = "baselime.io/baselime/baselime"
    }
  }
}

provider "baselime" {
  # example configuration here
  api_key = ""
  api_host = "go.baselime.io"
#  api_host = "localhost:32768"
#  api_scheme = "http"
}

resource "baselime_query" "terraformed" {
  name        = "terraformed-query"
  description = "This query was created by Terraform"
  service     = "default"
  datasets    = ["lambda-logs"]
  filters     = [
    {
      key       = "message"
      operation = "INCLUDES"
      value     = "error"
      type      = "string"
    }
  ]
  filter_combination = "AND"
  calculations       = [
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
  limit  = 10
  needle = {
    value      = ".*"
    is_regex   = true
    match_case = false
  }
}

resource "baselime_alert" "terraformed" {
  name        = "terraformed-alert"
  description = "This alert was created by Terraform"
  service     = "default"
  enabled     = true
  channels = [
    {
      type    = "email"
      targets = ["maksym@baselime.io"]
    }
  ]
  query     = baselime_query.terraformed.id
  threshold = {
    operator = "GREATER_THAN"
    value     = 0
  }
  frequency = "10m"
  window    = "5m"
}

resource "baselime_dashboard" "terraformed" {
  name        = "terraformed-dashboard"
  description = "This alert was created by Terraform"
  service     = "default"
  widgets     = [
    {
      query_id     = baselime_query.terraformed.id
      type        = "timeseries"
      name        = "Line Chart"
      description = "This is a line chart"
    }
  ]
}