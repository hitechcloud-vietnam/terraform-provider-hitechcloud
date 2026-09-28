---
page_title: "hitechcloud_service_resources Data Source - hitechcloud"
subcategory: "Compute"
description: |-
  Reads the resource summary and upgrade options of one service.
---

# hitechcloud_service_resources (Data Source)

Reads the resource summary and upgrade options of one service
(`GET /api/service/{id}/resources`).

## Example Usage

```terraform
data "hitechcloud_service_resources" "web" {
  service_id = "1"
}

output "upgrade_options" {
  value = data.hitechcloud_service_resources.web.upgrade_options
}
```

## Schema

### Required

- `service_id` (String) ID of the service to read.

### Read-Only

- `id` (String) The service ID (same as `service_id`).
- `resources` (Map of String) Resource counters of the service.
- `upgrade_options` (Map of String) Available upgrade options reported by the
  API.
