---
page_title: "hitechcloud_payment_fees Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Lists the fees of the payment methods available to the HiTechCloud account.
---

# hitechcloud_payment_fees (Data Source)

Lists the fees of the payment methods available to the HiTechCloud account
(`GET /api/payment/fees`).

## Example Usage

```terraform
data "hitechcloud_payment_fees" "all" {}

output "fee_table" {
  value = { for f in data.hitechcloud_payment_fees.all.fees : f.name => f.fee }
}
```

## Schema

### Read-Only

- `fees` (Attributes List) Payment method fees. (see
  [nested schema](#nestedatt--fees))
- `id` (String) Data source identifier (`payment_fees`).

<a id="nestedatt--fees"></a>

### Nested Schema for `fees`

Read-Only:

- `fee` (String) Fee (percentage or amount).
- `module` (String) Payment module / gateway key.
- `name` (String) Display name.
