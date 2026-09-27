---
page_title: "hitechcloud_dns_record Resource - hitechcloud"
subcategory: "DNS"
description: |-
  Manages a DNS record inside a HiTechCloud DNS zone.
---

# hitechcloud_dns_record (Resource)

Manages a DNS record inside a HiTechCloud DNS zone
(`/api/service/{service_id}/dns/{zone_id}/record`).

## Example Usage

```terraform
resource "hitechcloud_dns_record" "www" {
  service_id = "1"
  zone_id    = "zone-123"
  name       = "www"
  type       = "A"
  content    = "203.0.113.10"
  ttl        = 3600
}
```

## Schema

### Required

- `content` (String) Record content / value.
- `name` (String) Record name (relative to the zone).
- `service_id` (String) ID of the hosted service owning the zone. Changing this
  forces a new resource.
- `type` (String) Record type (A, AAAA, CNAME, MX, TXT, ...).
- `zone_id` (String) ID of the DNS zone. Changing this forces a new resource.

### Optional

- `priority` (Number) Record priority (MX/SRV). Defaults to the API default.
- `ttl` (Number) Time to live in seconds. Defaults to the API default.

### Read-Only

- `id` (String) Composite identifier `service_id/zone_id/record_id`.
- `record_id` (String) Identifier of the record as returned by the API.

## Import

Import is supported using the composite ID `service_id/zone_id/record_id`:

```shell
terraform import hitechcloud_dns_record.www 1/zone-123/rec-456
```
