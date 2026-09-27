---
page_title: "hitechcloud_domain Data Source - hitechcloud"
subcategory: "Domains"
description: |-
  Retrieves a single registered HiTechCloud domain by ID or name.
---

# hitechcloud_domain (Data Source)

Retrieves a single registered HiTechCloud domain by ID or name
(`GET /api/domain/{domain_id}` / lookup by name). Exactly one of `domain_id`
or `name` must be set.

## Example Usage

```terraform
data "hitechcloud_domain" "example" {
  name = "example.com"
}

output "domain_expiry" {
  value = data.hitechcloud_domain.example.expires
}
```

## Schema

### Optional

- `domain_id` (String) ID of the domain. Exactly one of `domain_id` or `name`.
- `name` (String) Domain name. Exactly one of `domain_id` or `name`.

### Read-Only

- `autorenew` (Boolean) Auto-renewal enabled.
- `expires` (String) Expiry date.
- `id` (String) Domain identifier.
- `id_protection` (Boolean) WHOIS ID protection enabled.
- `nameservers` (Set of String) Nameservers of the domain.
- `reg_date` (String) Registration date.
- `registrar` (String) Registrar.
- `registrar_lock` (Boolean) Registrar lock enabled.
- `status` (String) Domain status.
