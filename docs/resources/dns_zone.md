---
page_title: "hitechcloud_dns_zone Resource - hitechcloud"
subcategory: "DNS"
description: |-
  Manages a DNS zone of a HiTechCloud hosted service.
---

# hitechcloud_dns_zone (Resource)

Manages a DNS zone of a HiTechCloud hosted service (`POST/GET/DELETE
/api/service/{service_id}/dns`). Records of the zone are exposed as a computed
`records` list; to manage individual records declaratively use the
[`hitechcloud_dns_record`](dns_record.md) resource.

## Example Usage

```terraform
resource "hitechcloud_dns_zone" "example" {
  service_id = "1"
  name       = "example.com"
}
```

## Schema

### Required

- `name` (String) Zone name (domain). Changing this forces a new resource.
- `service_id` (String) ID of the hosted service owning the zone. Changing this
  forces a new resource.

### Optional

- `records` (Attributes Set) Initial records of the zone. Computed values
  reflect the records currently present. (see [nested schema](#nestedatt--records))

### Read-Only

- `id` (String) Composite identifier `service_id/zone_id`.
- `zone_id` (String) Identifier of the zone as returned by the API.

<a id="nestedatt--records"></a>

### Nested Schema for `records`

Read-Only:

- `content` (String) Record content / value.
- `id` (String) Record identifier.
- `name` (String) Record name.
- `priority` (Number) Record priority (MX/SRV).
- `ttl` (Number) Time to live in seconds.
- `type` (String) Record type (A, AAAA, CNAME, MX, TXT, ...).

## Import

Import is supported using the composite ID `service_id/zone_id`:

```shell
terraform import hitechcloud_dns_zone.example 1/zone-123
```
