---
page_title: "hitechcloud_affiliate Data Source - hitechcloud"
subcategory: "Portal"
description: |-
  Reads the caller's affiliate program data: summary, campaigns, commissions, payouts, vouchers and plans.
---

# hitechcloud_affiliate (Data Source)

Reads the caller's affiliate program data: summary, campaigns, commissions,
payouts, vouchers and commission plans
(`GET /api/affiliates/summary`, `/affiliates/campaigns`, `/affiliates/commissions`,
`/affiliates/payouts`, `/affiliates/vouchers`, `/affiliates/commissionplans`).

## Example Usage

```terraform
data "hitechcloud_affiliate" "me" {}

output "campaign_names" {
  value = [for c in data.hitechcloud_affiliate.me.campaigns : c.name]
}
```

## Schema

### Read-Only

- `campaigns` (Attributes List) Affiliate campaigns. (see
  [nested schema](#nestedatt--campaigns))
- `commission_plans` (Attributes List) Available commission plans. (see
  [nested schema](#nestedatt--commission_plans))
- `commissions` (Attributes List) Earned commissions. (see
  [nested schema](#nestedatt--commissions))
- `id` (String) Static identifier for this query (`affiliate`).
- `payouts` (Attributes List) Payout history. (see
  [nested schema](#nestedatt--payouts))
- `summary` (Map of String) Affiliate account summary counters and balances.
- `vouchers` (Attributes List) Affiliate vouchers. (see
  [nested schema](#nestedatt--vouchers))

<a id="nestedatt--campaigns"></a>

### Nested Schema for `campaigns`

Read-Only:

- `amount` (String) Monetary amount, when present.
- `date` (String) Date of the row, when present.
- `extra` (Map of String) Remaining fields of the row.
- `id` (String) Row identifier.
- `name` (String) Display name of the row.
- `status` (String) Status of the row.

The same nested schema applies to `commission_plans`, `commissions`, `payouts`
and `vouchers`.
