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

The provider supports two authentication flows:

1. **Username / password login** — the provider calls `POST /api/login`, which
   returns two tokens (`token` + `refresh_token`). The access token is used as
   the bearer token and is renewed automatically with `POST /api/token` when it
   expires (the failed request is then replayed).
2. **Pre-issued token** — supply an existing access token via the `token`
   attribute. Optionally provide `refresh_token` so the provider can still
   renew the access token automatically.

Credentials may also come from the `HITECHCLOUD_USERNAME`,
`HITECHCLOUD_PASSWORD`, `HITECHCLOUD_TOKEN` and `HITECHCLOUD_REFRESH_TOKEN`
environment variables.

## Example Usage

```terraform
terraform {
  required_providers {
    hitechcloud = {
      source  = "hitechcloud-vietnam/hitechcloud"
      version = "~> 1.0"
    }
  }
}

provider "hitechcloud" {
  username = var.hitechcloud_username # or use HITECHCLOUD_USERNAME
  password = var.hitechcloud_password # or use HITECHCLOUD_PASSWORD
  endpoint = "https://api.hitechcloud.vn"
}
```

Or with a pre-issued token:

```terraform
provider "hitechcloud" {
  token         = var.hitechcloud_token # or use HITECHCLOUD_TOKEN
  refresh_token = var.hitechcloud_refresh_token # optional auto-renewal
}
```

## Schema

### Optional

- `endpoint` (String) Base URL of the HiTechCloud User API. May also be set
  with the `HITECHCLOUD_ENDPOINT` environment variable. Defaults to
  `https://api.hitechcloud.vn`.
- `password` (String, Sensitive) Account password used with `username` for
  `POST /api/login` (two-token flow). May also be set with the
  `HITECHCLOUD_PASSWORD` environment variable.
- `refresh_token` (String, Sensitive) Refresh token used to renew the access
  token automatically via `POST /api/token`. May also be set with the
  `HITECHCLOUD_REFRESH_TOKEN` environment variable. Filled automatically when
  logging in with `username`/`password`.
- `request_timeout` (String) Per-request timeout as a Go duration string, for
  example `60s` or `2m`. May also be set with the `HITECHCLOUD_REQUEST_TIMEOUT`
  environment variable. Defaults to `60s`.
- `token` (String, Sensitive) API bearer token used to authenticate against the
  HiTechCloud User API. May also be set with the `HITECHCLOUD_TOKEN`
  environment variable. When empty, `username`/`password` login is used.
- `username` (String) Account username for `POST /api/login` (two-token
  flow). May also be set with the `HITECHCLOUD_USERNAME` environment
  variable.
