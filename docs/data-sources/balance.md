---
page_title: "hitechcloud_balance Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Retrieves the account balance of the HiTechCloud account.
---

# hitechcloud_balance (Data Source)

Retrieves the account balance of the HiTechCloud account (`GET /api/balance`).

## Example Usage

```terraform
data "hitechcloud_balance" "current" {}

output "balance" {
  value = data.hitechcloud_balance.current.balance
}
```

## Schema

### Read-Only

- `balance` (String) Account balance.
- `credit` (String) Available credit.
- `currency` (String) Currency.
- `id` (String) Data source identifier (`balance`).
