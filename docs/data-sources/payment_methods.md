---
page_title: "hitechcloud_payment_methods Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Lists the payment methods available to the HiTechCloud account.
---

# hitechcloud_payment_methods (Data Source)

Lists the payment methods available to the HiTechCloud account
(`GET /api/payment`).

## Example Usage

```terraform
data "hitechcloud_payment_methods" "all" {}

output "method_names" {
  value = [for m in data.hitechcloud_payment_methods.all.methods : m.name]
}
```

## Schema

### Read-Only

- `id` (String) Data source identifier (`payment_methods`).
- `methods` (Attributes List) Available payment methods. (see
  [nested schema](#nestedatt--methods))

<a id="nestedatt--methods"></a>

### Nested Schema for `methods`

Read-Only:

- `id` (String) Payment module / gateway key.
- `name` (String) Display name.
