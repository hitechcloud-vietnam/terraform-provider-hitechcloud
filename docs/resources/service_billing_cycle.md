---
page_title: "hitechcloud_service_billing_cycle Resource - hitechcloud"
subcategory: "Billing"
description: |-
  Manages the billing cycle of a HiTechCloud service.
---

# hitechcloud_service_billing_cycle (Resource)

Manages the billing cycle of a HiTechCloud service
(`POST /api/service/{id}/cycle`). Destroying the resource leaves the cycle
unchanged (the API has no reset).

## Example Usage

```terraform
resource "hitechcloud_service_billing_cycle" "vm" {
  service_id = "1"
  cycle      = "annually"
}
```

## Schema

### Required

- `cycle` (String) Billing cycle, for example `monthly`, `quarterly`,
  `annually`.
- `service_id` (String) ID of the service. Changing this forces a new resource.

### Read-Only

- `id` (String) The service ID (same as `service_id`).

## Import

Import is supported using the service ID:

```shell
terraform import hitechcloud_service_billing_cycle.vm 1
```
