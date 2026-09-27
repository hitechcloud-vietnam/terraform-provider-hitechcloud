---
page_title: "hitechcloud_domain_settings Resource - hitechcloud"
subcategory: "Domains"
description: |-
  Manages settings of a registered HiTechCloud domain: nameservers, auto-renew, registrar lock and ID protection.
---

# hitechcloud_domain_settings (Resource)

Manages settings of a registered HiTechCloud domain: nameservers, auto-renewal,
registrar lock and ID protection (`/api/domain/{domain_id}/...`). Only the
settings you configure are managed; unconfigured settings are left untouched.
Deleting the resource only removes it from state (the settings themselves
persist).

## Example Usage

```terraform
resource "hitechcloud_domain_settings" "example" {
  domain_id = "10"

  nameservers = [
    "ns1.hitechcloud.vn",
    "ns2.hitechcloud.vn",
  ]

  autorenew     = true
  registrar_lock = true
}
```

## Schema

### Required

- `domain_id` (String) ID of the registered domain. Changing this forces a new
  resource.

### Optional

- `autorenew` (Boolean) Enable automatic renewal.
- `id_protection` (Boolean) Enable WHOIS ID protection.
- `nameservers` (Set of String) Nameservers of the domain. An empty set
  restores the default nameservers.
- `registrar_lock` (Boolean) Enable the registrar (transfer) lock.

### Read-Only

- `id` (String) Identifier of the domain settings (`domain_id`).

## Import

Import is supported using the domain ID:

```shell
terraform import hitechcloud_domain_settings.example 10
```
