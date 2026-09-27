---
page_title: "hitechcloud_dns_zones Data Source - hitechcloud"
subcategory: "DNS"
description: |-
  Lists the DNS zones (with their records) of a HiTechCloud hosted service.
---

# hitechcloud_dns_zones (Data Source)

Lists the DNS zones of a HiTechCloud hosted service
(`GET /api/service/{service_id}/dns`), including each zone's records.

## Example Usage

```terraform
data "hitechcloud_dns_zones" "all" {
  service_id = "1"
}

output "zone_names" {
  value = [for z in data.hitechcloud_dns_zones.all.zones : z.name]
}
```

## Schema

### Required

- `service_id` (String) ID of the hosted service.

### Read-Only

- `id` (String) Data source identifier (same as `service_id`).
- `zones` (Attributes List) DNS zones of the service. (see
  [nested schema](#nestedatt--zones))

<a id="nestedatt--zones"></a>

### Nested Schema for `zones`

Read-Only:

- `id` (String) Zone identifier.
- `name` (String) Zone name.
- `records` (Attributes Set) Records of the zone. (see
  [nested schema](#nestedatt--zones--records))

<a id="nestedatt--zones--records"></a>

### Nested Schema for `zones.records`

Read-Only:

- `content` (String) Record content / value.
- `id` (String) Record identifier.
- `name` (String) Record name.
- `priority` (Number) Record priority (MX/SRV).
- `ttl` (Number) Time to live in seconds.
- `type` (String) Record type.
