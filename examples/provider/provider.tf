terraform {
  required_providers {
    baselime = {
      version = "~> 1.0.0"
#      source  = "baselime.io/baselimeio/baselimeio"
      source  = "baselime.io/baselimeio/baselimeio"
    }
  }
}

provider "baselime" {
  # example configuration here
  api_key = "foo"
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
  frequency = "5m"
  window    = "5m"
}