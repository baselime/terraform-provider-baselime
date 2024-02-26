# Baselime Observability as Code Terraform Provider

Observability as Code with the [Baselime](https://baselime.io) Terraform provider.

## Resources

[Documentation](https://registry.terraform.io/providers/baselime/baselime/latest/docs)
[Examples](https://github.com/baselime/terraform-provider-baselime/tree/main/examples/resources)

## Community
If you have any questions or want to discuss Baselime, please join our [Slack community](https://join.slack.com/t/baselimecommunity/shared_invite/zt-24fbumkc5-9O6qIj92xW_CbQSHeKT7CQ).

## Using the provider
```terraform
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
```

#### Api Key
To find your key:
1. Navigate to https://console.baselime.io
2. Select the workspace you need
3. Select the environment you want to get the key for
4. Click on the "API Keys" button on the left-hand side menu (key icon)

The API key can be supplied via the `BASELIME_API_KEY` environment variable as well.


#### Resource types
- [Query](https://registry.terraform.io/providers/baselime/baselime/latest/docs/resources/query)
- [Dashboard](https://registry.terraform.io/providers/baselime/baselime/latest/docs/resources/dashboard)
- [Alert](https://registry.terraform.io/providers/baselime/baselime/latest/docs/resources/alert)