---
page_title: "hitechcloud_ai_instance_types Data Source - hitechcloud"
subcategory: "AI Factory"
description: |-
  Lists the AI instance types of a HiTechCloud service.
---

# hitechcloud_ai_instance_types (Data Source)

Lists the AI instance types of a HiTechCloud service
(`GET /api/service/{service_id}/ai/instancetype`).

## Example Usage

```terraform
data "hitechcloud_ai_instance_types" "all" {
  service_id = "1"
}

output "gpu_types" {
  value = [for t in data.hitechcloud_ai_instance_types.all.types : t.name]
}
```

## Schema

### Required

- `service_id` (String) ID of the hosted service.

### Read-Only

- `id` (String) Data source identifier (same as `service_id`).
- `types` (Attributes List) Available instance types. (see
  [nested schema](#nestedatt--types))

<a id="nestedatt--types"></a>

### Nested Schema for `types`

Read-Only:

- `description` (String) Type description.
- `name` (String) Type name.
- `vcpus` (Number) Number of vCPUs.
