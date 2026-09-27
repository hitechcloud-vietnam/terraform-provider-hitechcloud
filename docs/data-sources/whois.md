---
page_title: "hitechcloud_whois Data Source - hitechcloud"
subcategory: "Domains"
description: |-
  Performs a WHOIS lookup for a domain.
---

# hitechcloud_whois (Data Source)

Performs a WHOIS lookup for a domain (`GET /api/whois/{domain}`). The raw WHOIS
fields are exposed as the `attributes` map.

## Example Usage

```terraform
data "hitechcloud_whois" "example" {
  domain = "example.com"
}

output "registrar" {
  value = try(data.hitechcloud_whois.example.attributes["Registrar"], "unknown")
}
```

## Schema

### Required

- `domain` (String) Domain name to look up.

### Read-Only

- `attributes` (Map of String) WHOIS attributes of the domain.
- `id` (String) Data source identifier (same as `domain`).
