---
page_title: "hitechcloud_domain_dns_record Resource - hitechcloud"
subcategory: "Domains"
description: |-
  Manages a DNS record of a registered HiTechCloud domain.
---

# hitechcloud_domain_dns_record (Resource)

Manages a DNS record of a registered HiTechCloud domain
(`/api/domain/{domain_id}/dns/record`).

## Example Usage

```terraform
resource "hitechcloud_domain_dns_record" "mail" {
  domain_id = "10"
  name      = "mail"
  type      = "A"
  content   = "203.0.113.25"
}
```

## Schema

### Required

- `content` (String) Record content / value.
- `domain_id` (String) ID of the registered domain. Changing this forces a new
  resource.
- `name` (String) Record name (relative to the domain).
- `type` (String) Record type (A, AAAA, CNAME, MX, TXT, ...).

### Optional

- `priority` (Number) Record priority (MX/SRV).
- `ttl` (Number) Time to live in seconds.

### Read-Only

- `id` (String) Composite identifier `domain_id/record_id`.
- `record_id` (String) Identifier of the record as returned by the API.

## Import

Import is supported using the composite ID `domain_id/record_id`:

```shell
terraform import hitechcloud_domain_dns_record.mail 10/rec-456
```
