---
page_title: "hitechcloud_domain_dns_types Data Source - hitechcloud"
subcategory: "Domains"
description: |-
  Lists the DNS record types supported by a registered HiTechCloud domain.
---

# hitechcloud_domain_dns_types (Data Source)

Lists the DNS record types supported by a registered HiTechCloud domain
(`GET /api/domain/{domain_id}/dns/types`).

## Example Usage

```terraform
data "hitechcloud_domain_dns_types" "example" {
  domain_id = "10"
}

output "supported_types" {
  value = data.hitechcloud_domain_dns_types.example.types
}
```

## Schema

### Required

- `domain_id` (String) ID of the registered domain.

### Read-Only

- `id` (String) Data source identifier (same as `domain_id`).
- `types` (Set of String) Supported DNS record types.
