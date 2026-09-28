---
page_title: "hitechcloud_service_label Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Manages the display label of a HiTechCloud service.
---

# hitechcloud_service_label (Resource)

Manages the display label of a HiTechCloud service
(`POST /api/service/{id}/label`). Destroying the resource clears the label.

## Example Usage

```terraform
resource "hitechcloud_service_label" "web" {
  service_id = "1"
  label      = "production-web"
}
```

## Schema

### Required

- `label` (String) Display label of the service.
- `service_id` (String) ID of the service. Changing this forces a new resource.

### Read-Only

- `id` (String) The service ID (same as `service_id`).

## Import

Import is supported using the service ID:

```shell
terraform import hitechcloud_service_label.web 1
```
