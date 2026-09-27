---
page_title: "hitechcloud_domain_availability Data Source - hitechcloud"
subcategory: "Domains"
description: |-
  Checks whether a domain name is available for registration.
---

# hitechcloud_domain_availability (Data Source)

Checks whether a domain name is available for registration
(`POST /api/domain/lookup`).

## Example Usage

```terraform
data "hitechcloud_domain_availability" "example" {
  name = "example.com"
}

output "is_available" {
  value = data.hitechcloud_domain_availability.example.available
}
```

## Schema

### Required

- `name` (String) Domain name to check.

### Read-Only

- `available` (Boolean) Whether the domain is available for registration.
- `id` (String) Data source identifier (same as `name`).
- `price` (String) Registration price, when returned by the API.
