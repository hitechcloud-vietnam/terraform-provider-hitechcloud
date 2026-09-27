---
page_title: "hitechcloud_certificates Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Lists the SSL certificates of the HiTechCloud account.
---

# hitechcloud_certificates (Data Source)

Lists the SSL certificates of the HiTechCloud account (`GET /api/certificate`).

## Example Usage

```terraform
data "hitechcloud_certificates" "all" {}

output "certificate_domains" {
  value = [for c in data.hitechcloud_certificates.all.certificates : c.domain]
}
```

## Schema

### Read-Only

- `certificates` (Attributes List) SSL certificates. (see
  [nested schema](#nestedatt--certificates))
- `id` (String) Data source identifier (`certificates`).

<a id="nestedatt--certificates"></a>

### Nested Schema for `certificates`

Read-Only:

- `domain` (String) Certificate domain.
- `expires` (String) Expiry date.
- `id` (String) Certificate identifier.
- `status` (String) Certificate status.
- `type` (String) Certificate product / type.
