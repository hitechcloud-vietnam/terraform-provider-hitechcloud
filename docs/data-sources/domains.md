---
page_title: "hitechcloud_domains Data Source - hitechcloud"
subcategory: "Domains"
description: |-
  Lists the registered domains of the HiTechCloud account.
---

# hitechcloud_domains (Data Source)

Lists the registered domains of the HiTechCloud account (`GET /api/domain`),
including nameservers and lock status.

## Example Usage

```terraform
data "hitechcloud_domains" "all" {}

output "domain_names" {
  value = [for d in data.hitechcloud_domains.all.domains : d.name]
}
```

## Schema

### Read-Only

- `domains` (Attributes List) Registered domains. (see
  [nested schema](#nestedatt--domains))
- `id` (String) Data source identifier (`domains`).

<a id="nestedatt--domains"></a>

### Nested Schema for `domains`

Read-Only:

- `autorenew` (Boolean) Auto-renewal enabled.
- `expires` (String) Expiry date.
- `id` (String) Domain identifier.
- `id_protection` (Boolean) WHOIS ID protection enabled.
- `name` (String) Domain name.
- `nameservers` (Set of String) Nameservers of the domain.
- `reg_date` (String) Registration date.
- `registrar` (String) Registrar.
- `registrar_lock` (Boolean) Registrar lock enabled.
- `status` (String) Domain status.
