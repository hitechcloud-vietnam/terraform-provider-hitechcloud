---
page_title: "hitechcloud_service Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Retrieves a single HiTechCloud service by ID.
---

# hitechcloud_service (Data Source)

Retrieves a single HiTechCloud service by ID (`GET /api/service/{service_id}`).

## Example Usage

```terraform
data "hitechcloud_service" "web" {
  service_id = "1"
}

output "service_status" {
  value = data.hitechcloud_service.web.status
}
```

## Schema

### Required

- `service_id` (String) ID of the service to look up.

### Read-Only

- `amount` (String) Recurring amount.
- `billing_cycle` (String) Billing cycle.
- `domain` (String) Primary domain of the service.
- `group` (String) Service group.
- `id` (String) Data source identifier (same as `service_id`).
- `ip` (String) Primary IP address.
- `label` (String) Service label.
- `name` (String) Service name.
- `next_due` (String) Next due date.
- `register_date` (String) Registration date.
- `status` (String) Service status.
