---
page_title: "HiTechCloud Provider"
description: |-
  Terraform provider for the HiTechCloud platform: compute, AI Factory, DNS, domains, S3 storage, billing and more.
---

# HiTechCloud Provider

The HiTechCloud provider manages resources on the [HiTechCloud](https://hitechcloud.vn)
platform through its User API. It supports compute (VMs, interfaces, firewall
rules, IPs), AI Factory (GPU instances, SSH keys, volumes, templates, clusters),
DNS zones and records, domain registration settings, S3 storage (buckets and
sub-users), billing data and portal utilities.

Use the navigation to the left to read about the available resources and data
sources.

## Authentication

The provider authenticates with a HiTechCloud API bearer token (issued by
`POST /api/login`). Provide it via the `token` attribute or the
`HITECHCLOUD_TOKEN` environment variable.

## Example Usage

```terraform
terraform {
  required_providers {
    hitechcloud = {
      source  = "hitechcloud-vietnam/hitechcloud"
      version = "~> 0.1"
    }
  }
}

provider "hitechcloud" {
  token    = var.hitechcloud_token # or use HITECHCLOUD_TOKEN
  endpoint = "https://api.hitechcloud.vn"
}
```

## Schema

### Optional

- `endpoint` (String) Base URL of the HiTechCloud User API. May also be set
  with the `HITECHCLOUD_ENDPOINT` environment variable. Defaults to
  `https://api.hitechcloud.vn`.
- `request_timeout` (String) Per-request timeout as a Go duration string, for
  example `60s` or `2m`. May also be set with the `HITECHCLOUD_REQUEST_TIMEOUT`
  environment variable. Defaults to `60s`.
- `token` (String, Sensitive) API bearer token used to authenticate against the
  HiTechCloud User API. May also be set with the `HITECHCLOUD_TOKEN`
  environment variable.
