---
page_title: "hitechcloud_partner Data Source - hitechcloud"
subcategory: "Portal"
description: |-
  Reads the caller's HiTechCloud partner program data: profile, tiers, rate card, wallet and leads.
---

# hitechcloud_partner (Data Source)

Reads the caller's HiTechCloud partner program data: profile, tiers, rate card,
wallet and recent leads (`GET /api/partner`, `/partner/tiers`,
`/partner/ratecard`, `/partner/wallet`, `/partner/leads`).

## Example Usage

```terraform
data "hitechcloud_partner" "me" {}

output "partner_tiers" {
  value = [for t in data.hitechcloud_partner.me.tiers : t.name]
}
```

## Schema

### Read-Only

- `id` (String) Partner identifier (or email) reported by the API.
- `profile` (Map of String) Partner profile fields.
- `rate_card` (Map of String) Commission rate card values.
- `recent_leads` (Attributes List) Most recent partner leads registered by the
  caller. (see [nested schema](#nestedatt--recent_leads))
- `tiers` (Attributes List) Partner program tiers and their requirements. (see
  [nested schema](#nestedatt--tiers))
- `wallet` (Map of String) Partner wallet balances and payout settings.

<a id="nestedatt--recent_leads"></a>

### Nested Schema for `recent_leads`

Read-Only:

- `company` (String) Lead company.
- `created_at` (String) Registration timestamp.
- `email` (String) Lead email.
- `id` (String) Lead identifier.
- `status` (String) Lead status.

<a id="nestedatt--tiers"></a>

### Nested Schema for `tiers`

Read-Only:

- `commission` (String) Commission rate of the tier.
- `level` (Number) Tier level.
- `name` (String) Tier name.
- `requirement` (String) Requirement description.
