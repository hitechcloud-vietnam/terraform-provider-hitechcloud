---
page_title: "hitechcloud_services Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Lists all services of the HiTechCloud account.
---

# hitechcloud_services (Data Source)

Lists all services of the HiTechCloud account (`GET /api/service`).

## Example Usage

```terraform
data "hitechcloud_services" "all" {}

output "active_services" {
  value = [for s in data.hitechcloud_services.all.services : s.name if s.status == "active"]
}
```

## Schema

### Read-Only

- `id` (String) Data source identifier (`services`).
- `services` (Attributes List) Services of the account. (see
  [nested schema](#nestedatt--services))

<a id="nestedatt--services"></a>

### Nested Schema for `services`

Read-Only:

- `amount` (String) Recurring amount.
- `billing_cycle` (String) Billing cycle.
- `domain` (String) Primary domain of the service.
- `group` (String) Service group.
- `id` (String) Service identifier.
- `ip` (String) Primary IP address.
- `label` (String) Service label.
- `name` (String) Service name.
- `next_due` (String) Next due date.
- `register_date` (String) Registration date.
- `status` (String) Service status.
