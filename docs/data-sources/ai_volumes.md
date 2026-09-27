---
page_title: "hitechcloud_ai_volumes Data Source - hitechcloud"
subcategory: "AI Factory"
description: |-
  Lists the volumes of the HiTechCloud AI Factory.
---

# hitechcloud_ai_volumes (Data Source)

Lists the volumes of the HiTechCloud AI Factory
(`GET /api/service/{service_id}/volumes`).

## Example Usage

```terraform
data "hitechcloud_ai_volumes" "all" {
  service_id = "1"
}

output "volume_names" {
  value = [for v in data.hitechcloud_ai_volumes.all.volumes : v.name]
}
```

## Schema

### Required

- `service_id` (String) ID of the hosted service.

### Read-Only

- `id` (String) Data source identifier (same as `service_id`).
- `volumes` (Attributes List) AI volumes. (see
  [nested schema](#nestedatt--volumes))

<a id="nestedatt--volumes"></a>

### Nested Schema for `volumes`

Read-Only:

- `cloud` (String) Cloud of the volume.
- `id` (String) Volume identifier.
- `name` (String) Volume name.
- `region` (String) Region of the volume.
- `size_in_gb` (Number) Size in GB.
- `status` (String) Volume status.
