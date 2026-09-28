---
page_title: "hitechcloud_affiliate_adv Data Source - hitechcloud"
subcategory: "Portal"
description: |-
  Reads the advanced affiliate data of one client: profile, statistics, referrals, vouchers and commissions.
---

# hitechcloud_affiliate_adv (Data Source)

Reads the advanced affiliate data of one client: profile, statistics, referrals,
vouchers and commissions
(`GET /api/affiliates_adv/{client_id}/info`, `/stats`, `/referrals`,
`/vouchers`, `/commissions`).

## Example Usage

```terraform
data "hitechcloud_affiliate_adv" "partner" {
  client_id = "100"
}

output "referral_count" {
  value = length(data.hitechcloud_affiliate_adv.partner.referrals)
}
```

## Schema

### Required

- `client_id` (String) Affiliate client ID.

### Read-Only

- `commissions` (Attributes List) Commissions earned by the affiliate. (see
  [nested schema](#nestedatt--commissions))
- `id` (String) The client ID (same as `client_id`).
- `info` (Map of String) Affiliate profile fields.
- `referrals` (Attributes List) Referred customers. (see
  [nested schema](#nestedatt--referrals))
- `stats` (Map of String) Aggregate affiliate statistics.
- `vouchers` (Attributes List) Vouchers created under the affiliate. (see
  [nested schema](#nestedatt--vouchers))

<a id="nestedatt--referrals"></a>

### Nested Schema for `referrals`

Read-Only:

- `amount` (String) Monetary amount, when present.
- `date` (String) Date of the row, when present.
- `extra` (Map of String) Remaining fields of the row.
- `id` (String) Row identifier.
- `name` (String) Display name of the row.
- `status` (String) Status of the row.

The same nested schema applies to `commissions` and `vouchers`.
