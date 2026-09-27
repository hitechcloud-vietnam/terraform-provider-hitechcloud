---
page_title: "hitechcloud_domain_tlds Data Source - hitechcloud"
subcategory: "Domains"
description: |-
  Lists the TLDs available for domain registration at HiTechCloud.
---

# hitechcloud_domain_tlds (Data Source)

Lists the TLDs available for domain registration at HiTechCloud
(`GET /api/domain/tld`).

## Example Usage

```terraform
data "hitechcloud_domain_tlds" "all" {}

output "tld_names" {
  value = [for t in data.hitechcloud_domain_tlds.all.tlds : t.name]
}
```

## Schema

### Read-Only

- `id` (String) Data source identifier (`domain_tlds`).
- `tlds` (Attributes List) Available TLDs. (see
  [nested schema](#nestedatt--tlds))

<a id="nestedatt--tlds"></a>

### Nested Schema for `tlds`

Read-Only:

- `max_years` (Number) Maximum registration period in years.
- `min_years` (Number) Minimum registration period in years.
- `name` (String) TLD name (for example `com`).
- `price` (String) Registration price.
