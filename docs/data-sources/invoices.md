---
page_title: "hitechcloud_invoices Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Lists the invoices of the HiTechCloud account.
---

# hitechcloud_invoices (Data Source)

Lists the invoices of the HiTechCloud account (`GET /api/invoice`).

## Example Usage

```terraform
data "hitechcloud_invoices" "all" {}

output "unpaid_invoices" {
  value = [for i in data.hitechcloud_invoices.all.invoices : i.number if i.status != "paid"]
}
```

## Schema

### Read-Only

- `id` (String) Data source identifier (`invoices`).
- `invoices` (Attributes List) Invoices of the account. (see
  [nested schema](#nestedatt--invoices))

<a id="nestedatt--invoices"></a>

### Nested Schema for `invoices`

Read-Only:

- `currency` (String) Currency.
- `due_date` (String) Due date.
- `id` (String) Invoice identifier.
- `number` (String) Invoice number.
- `status` (String) Invoice status.
- `total` (String) Total amount due.
